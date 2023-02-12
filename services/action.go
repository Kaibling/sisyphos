package services

import (
	"errors"
	"fmt"
	"strings"

	"sisyphos/lib/apperrors"
	"sisyphos/lib/config"
	"sisyphos/lib/metadata"
	"sisyphos/lib/ssh"
	"sisyphos/lib/utils"
	"sisyphos/models"
)

type actionRepo interface {
	Update(name string, d *models.Action) (*models.Action, error)
	Create(actions []models.Action) ([]models.Action, error)
	ReadByName(name interface{}) (*models.Action, error)
	ReadIDs(ids []interface{}) ([]models.Action, error)
	ReadRuns(actionname interface{}) ([]models.Run, error)
}

type ActionService struct {
	repo          actionRepo
	permService   *PermissionService
	runLogService *RunService
	hostService   *HostService
}

func NewActionService(repo actionRepo) *ActionService {
	return &ActionService{repo: repo}
}

func (s *ActionService) AddPermissionService(p *PermissionService) {
	s.permService = p
}

func (s *ActionService) AddRunService(r *RunService) {
	s.runLogService = r
}

func (s *ActionService) AddHostService(r *HostService) {
	s.hostService = r
}

func (s *ActionService) Create(models []models.Action) ([]models.Action, error) {
	for i := 0; i < len(models); i++ {
		// add to default group, if no group is provided
		if len(models[i].Groups) == 0 {
			models[i].Groups = []string{defaultGroupName}
		}
		if err := models[i].Validate(); err != nil {
			return nil, err
		}
		models[i].Default()
	}
	return s.repo.Create(models)
}

func (s *ActionService) ReadByName(name interface{}) (*models.Action, error) {
	return s.repo.ReadByName(name)
}

func (s *ActionService) ReadRuns(actionname interface{}) ([]models.Run, error) {
	return s.repo.ReadRuns(actionname)
}

func (s *ActionService) ReadAllExtendedPermission(username string) ([]models.Action, error) {
	if s.permService == nil {
		return nil, errors.New("no permission service instantiated")
	}
	ids, err := s.permService.GetActionIDs(username)
	if err != nil {
		return nil, err
	}
	return s.repo.ReadIDs(ids)
}

func (s *ActionService) ReadIDs(ids []interface{}) ([]models.Action, error) {
	return s.repo.ReadIDs(ids)
}

func (s *ActionService) ReadAllFiltered(md metadata.MetaData, f filter) ([]models.Action, error) {
	id, err := f.Filter(md.Filter)
	if err != nil {
		return nil, err
	}
	return s.repo.ReadIDs(id)
}

func (s *ActionService) Update(name string, data *models.Action) (*models.Action, error) {
	return s.repo.Update(name, data)
}

type JSON = map[string]interface{}

func (s *ActionService) InitRun(r *models.Action) ([]models.Run, error) {
	if s.runLogService == nil {
		return nil, errors.New("no Run service instantiated")
	}
	if s.hostService == nil {
		return nil, errors.New("no host service instantiated")
	}
	r.Variables["failonerrors"] = utils.PtrRead(r.FailOnErrors)
	r.Variables = CombineVars(r.Variables, config.Config.GlobalVars)
	s.run(r, "")
	return s.runLogService.ReadByReqID()
}

func (s *ActionService) run(r *models.Action, parentID string) error {
	fmt.Printf("Start run %s pid:%s\n", utils.PtrRead(r.Name), parentID)
	execLog := models.NewRun(utils.PtrRead(r.Name),
		s.runLogService.repo.GetUsername(),
		s.runLogService.repo.GetRequestID(),
		parentID,
	)
	execLog.Status = apperrors.ScriptRunSuccess

	for _, tr := range r.Actions {
		t, err := s.ReadByName(tr.Name)
		if err != nil {
			execLog.Error = err.Error()
			execLog.SetEndTime()
			s.runLogService.Create(execLog)
			// return err
		}
		t.Variables["failonerrors"] = utils.PtrRead(t.FailOnErrors)
		t.Variables = CombineVars(r.Variables, t.Variables)
		if len(t.Hosts) == 0 {
			t.Hosts = r.Hosts
		}
		rerr := s.run(t, execLog.RunID)
		if t.Variables["failonerrors"].(bool) && rerr != nil {
			e := fmt.Errorf("script '%s' failed: %w", utils.PtrRead(r.Name), err)
			execLog.Error = e.Error()
			execLog.SetEndTime()
			s.runLogService.Create(execLog)
			return rerr
		}
		// TODO cancel if error ???
	}
	if utils.PtrRead(r.Script) != "" {
		fmt.Printf("run %s has script %s\n", *r.Name, utils.PtrRead(r.Script))
		if len(r.Hosts) == 0 {
			e := fmt.Errorf("no hosts for '%s'", utils.PtrRead(r.Name))
			execLog.Error = e.Error()
			execLog.SetEndTime()
			s.runLogService.Create(execLog)
			//return e
		}

		for _, connection := range r.Hosts {
			fmt.Printf("try ssh run %s on %s\n", utils.PtrRead(r.Name), utils.PtrRead(&connection.Name))
			hostExecLog := models.NewRun(*r.Name,
				s.runLogService.repo.GetUsername(),
				s.runLogService.repo.GetRequestID(),
				execLog.RequestID)
			hostExecLog.Host = &connection.Name

			sshc := ssh.NewSSHConnector()
			sshService := NewSSHService(sshc)
			cfg, err := s.hostService.GetSSHConfig(connection.Name)
			if err != nil {
				execLog.Error = err.Error()
				execLog.SetEndTime()
				s.runLogService.Create(execLog)
				return err
			}
			//cfg := connection.ToSSHConfig()

			cmd := replaceVariables(utils.PtrRead(r.Script), r.Variables)
			output, serr := sshService.RunCommand(*cfg, cmd)
			if serr != nil {
				hostExecLog.Error = serr.Error()
				hostExecLog.Status = apperrors.ScriptRunFailed
				execLog.Status = apperrors.ScriptRunFailed
			} else {
				hostExecLog.Status = apperrors.ScriptRunSuccess
			}
			hostExecLog.Output = output

			hostExecLog.SetEndTime()
			if _, err := s.runLogService.Create(hostExecLog); err != nil {
				fmt.Println(err.Error())
			}

			if r.Variables["failonerrors"].(bool) && serr != nil {
				return serr
			}
		}
	}
	execLog.SetEndTime()
	s.runLogService.Create(execLog)
	return nil
}

func CombineVars(base, override JSON) JSON {
	for k, v := range override {
		base[k] = v
	}
	return base
}

func replaceVariables(cmd string, variables map[string]interface{}) string {
	// suche {}, suche im map danache und ersetzt
	startActione := "{"
	endActione := "}"

	str := ""
	for _, s := range strings.Split(cmd, startActione) {
		// fmt.Println(s)
		ss := strings.Split(s, endActione)
		if len(ss) > 1 {
			found := ss[0]
			if val, ok := variables[found].(string); ok {
				str += val
			} else {
				str += startActione + found + endActione
			}
			str += strings.Join(ss[1:], "")
		} else {
			str += s
		}
	}
	return str
}
