package repos

type Logger interface {
	Infof(string, ...any)
	Warnf(string, ...any)
	Errorf(string, ...any)
	Debugf(string, ...any)
}
