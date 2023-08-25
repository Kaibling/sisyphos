package cluster

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

type logger interface {
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Debugf(string, ...any)
}

type Backend interface {
	Lock(key string, lockDuration time.Duration) (bool, error)
	UnLock(key string) error
}

type ClusterMaster struct {
	cfg      ClusterConfig
	backend  Backend
	IsMaster bool // TODO mutex
	id       string
}

type ClusterConfig struct {
	StartHook    func()
	StopHook     func()
	HeatBeatRate time.Duration // milliseconds
	Log          logger
}

func New(cfg ClusterConfig, be Backend) (*ClusterMaster, error) {
	ulid, _ := NewUlid()
	if cfg.Log == nil {
		return nil, fmt.Errorf("logger missing")
	}
	return &ClusterMaster{cfg, be, false, ulid}, nil
}

func (cm *ClusterMaster) ID() string {
	return cm.id
}

func (cm *ClusterMaster) Run() error {
	cm.cfg.Log.Infof("start cluster node. ID: %s", cm.id)
	waittime := rand.Intn(900)
	time.Sleep(time.Duration(waittime) * time.Millisecond)
	cm.cfg.Log.Debugf("waiting %d", waittime)
	gracePeriod := cm.cfg.HeatBeatRate
	// todo random start time
	for {
		// try being master
		ok, err := cm.backend.Lock(cm.id, cm.cfg.HeatBeatRate)
		if err != nil {
			cm.cfg.Log.Errorf(err.Error())
		}
		// if succeed and was not master, start hook
		if ok && !cm.IsMaster {
			cm.IsMaster = true
			cm.cfg.Log.Debugf("got cluster master")
			go cm.cfg.StartHook()
			cm.cfg.Log.Debugf("cluster start hook executed")
			gracePeriod = cm.cfg.HeatBeatRate - time.Duration(300)*time.Millisecond
		}
		// if no lock is abtained and node was master, stop hook
		if !ok && cm.IsMaster {
			cm.IsMaster = false
			cm.cfg.Log.Debugf("lost cluster master")
			go cm.cfg.StopHook()
			cm.cfg.Log.Debugf("stop hook executed")
			gracePeriod = cm.cfg.HeatBeatRate
		}
		time.Sleep(gracePeriod)
	}
}

func NewUlid() (string, error) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	if nl, err := ulid.New(ms, entropy); err != nil {
		return "", err
	} else {
		return nl.String(), nil
	}

}
