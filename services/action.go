package services

import (
	"errors"
	"fmt"
	"strings"

	"sisyphos/lib/config"
	"sisyphos/lib/metadata"
	"sisyphos/lib/ssh"
	"sisyphos/models"
)

type actionRepo interface {
	Update(name string, d *models.Action) (*models.Action, error)
	Create(actions []models.Action) ([]models.Action, error)
	ReadByName(name interface{}) (*models.Action, error)
	ReadIDs(ids []interface{}) ([]models.Action, error)
	ReadExtIDs(ids []interface{}) ([]models.ActionExt, error)
	ReadRuns(actionname interface{}) ([]models.Run, error)
	ReadExt(name interface{}) (*models.ActionExt, error)
}

type ActionService struct {
	repo          actionRepo
	permService   *PermissionService
	runLogService *RunService
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

func (s *ActionService) Create(models []models.Action) ([]models.Action, error) {
	for i := 0; i < len(models); i++ {
		// add to default group, if no group is provided
		if len(models[i].Groups) == 0 {
			models[i].Groups = []string{defaultGroupName}
		}
	}
	return s.repo.Create(models)
}

func (s *ActionService) ReadByName(name interface{}) (*models.Action, error) {
	return s.repo.ReadByName(name)
}

// func (s *ActionService) ReadAll() ([]models.Action, error) {
// 	return s.repo.ReadAll()
// }

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

func (s *ActionService) ReadExt(name interface{}) (*models.ActionExt, error) {
	return s.repo.ReadExt(name)
}

type JSON = map[string]interface{}

func (s *ActionService) InitRun(r *models.ActionExt) ([]models.Run, error) {
	if s.runLogService == nil {
		return nil, errors.New("no Run service instantiated")
	}
	r.Variables = CombineVars(r.Variables, config.Config.GlobalVars)
	s.run(r)
	return s.runLogService.ReadByReqID()
}

func (s *ActionService) run(r *models.ActionExt) error {
	fmt.Printf("Start run %s\n", r.Name)
	execLog := models.NewRun(*r.Name,
		s.runLogService.repo.GetUsername(),
		s.runLogService.repo.GetRequestID())

	for _, tr := range r.Triggers {
		t, err := s.ReadExt(tr.Name)
		if err != nil {
			execLog.Error = err.Error()
			execLog.SetEndTime()
			s.runLogService.Create(*execLog)
			// return err
		}
		t.Variables = CombineVars(t.Variables, r.Variables)
		if len(t.Hosts) == 0 {
			t.Hosts = r.Hosts
		}
		s.run(t)
		// TODO cancel if error ???
	}
	if *r.Script != "" {
		fmt.Printf("run %s has script %s\n", r.Name, r.Script)
		if len(r.Hosts) == 0 {
			e := fmt.Errorf("no hosts for '%s'", r.Name)
			execLog.Error = e.Error()
			execLog.SetEndTime()
			s.runLogService.Create(*execLog)
			return e
		}
		for _, connection := range r.Hosts {
			fmt.Printf("try ssh run %s on %s\n", r.Name, connection.Name)
			sshc := ssh.NewSSHConnector()
			sshService := NewSSHService(sshc)

			cfg := SSHConfig{
				Address:  connection.Address,
				Port:     connection.Port,
				Username: r.Variables["ssh_user"].(string),
				Password: r.Variables["ssh_password"].(string),
			}
			cmd := replaceVariables(*r.Script, r.Variables)
			output, err := sshService.RunCommand(cfg, cmd)
			// TODO cancel if error ???
			execLog.Output = output
			if err != nil {
				execLog.Error = err.Error()
			}

			// if err != nil {
			// 		//return err
			// }
		}
	}
	execLog.SetEndTime()
	s.runLogService.Create(*execLog)
	return nil
}

func CombineVars(v1, v2 JSON) JSON {
	for k, v := range v2 {
		v1[k] = v
	}
	return v1
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
