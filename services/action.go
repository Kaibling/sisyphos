package services

import (
	"errors"
	"fmt"
	"strings"

	"sisyphos/lib/apperrors"
	"sisyphos/lib/config"
	"sisyphos/lib/cron"
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
		if models[i].ScheduleExpr != nil {
			if !cron.Validate(*models[i].ScheduleExpr) {
				return nil, fmt.Errorf("schedule expression invalid")
			}
		}
		models[i].Default()
	}
	return s.repo.Create(models)
}

func (s *ActionService) ReadByName(name interface{}) (*models.Action, error) {
	return s.repo.ReadByName(name)
}

func (s *ActionService) ReadRuns(actionname interface{}) ([]models.Run, error) {
	runs, err := s.repo.ReadRuns(actionname)
	if err != nil {
		return nil, err
	}
	n := sortRunsChild(runs)
	return n, nil
}

func sortRunsChild(runs []models.Run) []models.Run {
	t := map[string]*models.Run{}
	for _, r := range runs {
		r := r
		t[r.ID] = &r
	}
	topLevel := []models.Run{}
	tmp := map[string]struct{}{}

	for _, a := range t {
		if a.ParentID != "" {
			tmp[a.ID] = struct{}{}
		}
	}

	loopbreak := 0
	for len(tmp) > 0 {
		for id := range tmp {
			child, ok := t[id]
			if !ok {
				fmt.Printf("child %s not found\n", id)
			}
			parent, ok := t[child.ParentID]
			if !ok {
				fmt.Printf("parent %s not found\n", child.ParentID)
				delete(tmp, id)
				continue
			}
			utils.PrettyJSON(parent)
			if parent.Childs == nil {
				parent.Childs = []*models.Run{}
			}
			parent.Childs = append(parent.Childs, child)
			delete(tmp, id)
		}
		loopbreak++
		if loopbreak > 10 {
			fmt.Println("broke")
			break
		}
	}

	// cleanup
	for _, u := range t {
		if u.ParentID == "" {
			topLevel = append(topLevel, *u)
		}
	}

	return topLevel
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
	if data.ScheduleExpr != nil {
		if !cron.Validate(*data.ScheduleExpr) {
			return nil, fmt.Errorf("schedule expression invalid")
		}
	}
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
	if rerr := s.run(r, ""); rerr != nil {
		return nil, rerr
	}

	return s.runLogService.ReadByReqID()
}

func (s *ActionService) run(r *models.Action, parentID string) error {
	execLog := models.NewRun(utils.PtrRead(r.Name),
		s.runLogService.repo.GetUsername(),
		s.runLogService.repo.GetRequestID(),
		parentID,
	)
	execLog.Status = apperrors.ScriptRunSuccess
	fmt.Printf("Start run %s pid:%s:execlog: %s\n", utils.PtrRead(r.Name), parentID, execLog.ID)
	for _, tr := range r.Actions {
		t, err := s.ReadByName(tr.Name)
		if err != nil {
			execLog.Error = err.Error()
			execLog.SetEndTime()
			if _, cerr := s.runLogService.Create(execLog); cerr != nil {
				return cerr
			}
		}
		t.Variables["failonerrors"] = utils.PtrRead(t.FailOnErrors)
		t.Variables = CombineVars(r.Variables, t.Variables)
		if len(t.Hosts) == 0 {
			t.Hosts = r.Hosts
		}
		rerr := s.run(t, execLog.ID)
		if t.Variables["failonerrors"].(bool) && rerr != nil {
			e := fmt.Errorf("script '%s' failed: %w", utils.PtrRead(r.Name), err)
			execLog.Error = e.Error()
			execLog.SetEndTime()
			if _, cerr := s.runLogService.Create(execLog); cerr != nil {
				return cerr
			}
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
			if _, cerr := s.runLogService.Create(execLog); cerr != nil {
				return cerr
			}
		}

		for _, connection := range r.Hosts {
			fmt.Printf("try ssh run %s on %s\n", utils.PtrRead(r.Name), utils.PtrRead(&connection.Name))
			hostExecLog := models.NewRun(*r.Name,
				s.runLogService.repo.GetUsername(),
				s.runLogService.repo.GetRequestID(),
				execLog.ID)
			hostExecLog.Host = &connection.Name

			sshc := ssh.NewSSHConnector()
			sshService := NewSSHService(sshc)
			cfg, err := s.hostService.GetSSHConfig(connection.Name)
			if err != nil {
				execLog.Error = err.Error()
				execLog.SetEndTime()
				if _, cerr := s.runLogService.Create(execLog); cerr != nil {
					return cerr
				}
				return err
			}
			// cfg := connection.ToSSHConfig()

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
	if _, cerr := s.runLogService.Create(execLog); cerr != nil {
		return cerr
	}
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
