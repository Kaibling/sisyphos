package tasker

type Task struct {
	ID           string
	ScheduleExpr string
	Action       Action
}

type Action interface {
	Run() error
}
