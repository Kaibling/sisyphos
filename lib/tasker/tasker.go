package tasker

// TODO addd logging
import (
	"context"
	"fmt"
	"sisyphos/models"
	"sync"
	"time"

	"github.com/adhocore/gronx"
)

type actionService interface {
	InitRun(r *models.Action) ([]models.Run, error)
}
type Tasker struct {
	tasks map[string]models.Action
	ctx   context.Context
	mux   sync.RWMutex
	as    actionService
}

func New(ctx context.Context, as actionService) *Tasker {
	return &Tasker{tasks: map[string]models.Action{}, ctx: ctx, mux: sync.RWMutex{}, as: as}
}

func (t *Tasker) Start() {
	go t.schedule()
}

func (t *Tasker) Remove(taskId string) {
	t.mux.Lock()
	defer t.mux.Unlock()
	delete(t.tasks, taskId)
}

func (t *Tasker) Add(task models.Action) {
	t.mux.Lock()
	defer t.mux.Unlock()
	t.tasks[task.ID] = task
}

func (tasker *Tasker) schedule() {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-tasker.ctx.Done():
			fmt.Println("canceled")
			return
		case <-ticker.C:
			gron := gronx.New()
			currentTime := time.Now().UTC()
			refTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), 0, 0, time.UTC)
			for k, v := range tasker.tasks {
				v := &v
				if due, err := gron.IsDue(*v.ScheduleExpr, refTime); err != nil {
					fmt.Println(err.Error())
					continue
				} else if due {
					fmt.Printf("execution of %s\n", k)
					if _, err := tasker.as.InitRun(v); err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
	}
}
