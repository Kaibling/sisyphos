package tasker

// TODO addd logging
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/adhocore/gronx"
)

type Tasker struct {
	tasks map[string]Task
	ctx   context.Context
	mux   sync.RWMutex
}

func New(ctx context.Context) *Tasker {
	return &Tasker{tasks: map[string]Task{}, ctx: ctx, mux: sync.RWMutex{}}
}

func (t *Tasker) Start() {
	go t.schedule()
}

func (t *Tasker) Remove(taskId string) {
	t.mux.Lock()
	defer t.mux.Unlock()
	delete(t.tasks, taskId)
}

func (t *Tasker) Add(task Task) {
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
				if due, err := gron.IsDue(v.ScheduleExpr, refTime); err != nil {
					fmt.Println(err.Error())
					continue
				} else if due {
					fmt.Printf("execution of %s\n", k)
					if err := v.Action.Run(); err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
	}
}
