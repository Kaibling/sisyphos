package log

import (
	"context"
	"fmt"

	"sisyphos/lib/reqctx"
	"sisyphos/lib/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field struct {
	Key   string
	Value any
}

func (f *Field) ToZap() zapcore.Field {
	return zap.Any(f.Key, f.Value)
}

type Logger struct {
	l             *zap.Logger
	defaultFields map[string]Field
}

func New() (*Logger, error) {
	cfg := zap.Config{
		Encoding:         "console", // json
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}
	return &Logger{l: logger, defaultFields: map[string]Field{}}, nil
}

func (l *Logger) AddDefaultField(key, value string) {
	f := Field{Key: key, Value: value}
	l.defaultFields[key] = f
}

func (l *Logger) Copy() *Logger {
	l2, _ := New()
	for k, v := range l.defaultFields {
		l2.defaultFields[k] = v
	}
	return l2
}

func (l *Logger) DE() {
	utils.PrettyJSON(l.defaultFields)
}

func (l *Logger) Raw() *zap.Logger {
	return l.l
}

func (l *Logger) prepareFields(ctx context.Context, fields ...Field) []zap.Field {
	f := []zapcore.Field{}
	if pk := ctx.Value("user"); pk != nil {
		f = append(f, zap.Any("user", pk.(string)))
	}

	if reqID := ctx.Value("requestid"); reqID != nil {
		f = append(f, zap.Any("request_id", reqID.(string)))
	}
	for _, fd := range fields {
		f = append(f, fd.ToZap())
	}
	return f
}

func (l *Logger) Printf(msg string, data ...interface{}) {
	f := Field{
		Key:   "gorm",
		Value: data,
	}
	l.info(context.TODO(), "gorm", f)
}

func get(ctx context.Context) *Logger {
	return ctx.Value(reqctx.String("logger")).(*Logger)
}

func (l *Logger) debug(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Debug(msg, f...)
}

func (l *Logger) info(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Info(msg, f...)
}

func (l *Logger) warn(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Warn(msg, f...)
}

func (l *Logger) error(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Error(msg, f...)
}

func Debug(ctx context.Context, msg string, fields ...Field) {
	l := get(ctx)
	zf := []Field{}
	for k, v := range l.defaultFields {
		found := false
		for _, f := range fields {
			if f.Key == k {
				zf = append(zf, f)
				found = true
			}
		}
		if !found {
			zf = append(zf, v)
		}
	}
	l.debug(ctx, msg, zf...)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	l := get(ctx)
	zf := []Field{}
	for k, v := range l.defaultFields {
		found := false
		for _, f := range fields {
			if f.Key == k {
				zf = append(zf, f)
				found = true
			}
		}
		if !found {
			zf = append(zf, v)
		}
	}
	l.info(ctx, msg, zf...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	l := get(ctx)
	zf := []Field{}
	for k, v := range l.defaultFields {
		found := false
		for _, f := range fields {
			if f.Key == k {
				zf = append(zf, f)
				found = true
			}
		}
		if !found {
			zf = append(zf, v)
		}
	}
	l.warn(ctx, msg, zf...)
}

func Error(ctx context.Context, err error, fields ...Field) {
	l := get(ctx)
	zf := []Field{}
	for k, v := range l.defaultFields {
		found := false
		for _, f := range fields {
			if f.Key == k {
				zf = append(zf, f)
				found = true
			}
		}
		if !found {
			zf = append(zf, v)
		}
	}
	l.error(ctx, err.Error(), zf...)
}

func (l *Logger) Infof(s string, a ...any) {
	zf := []zapcore.Field{}
	for _, f := range l.defaultFields {
		zf = append(zf, f.ToZap())
	}
	l.l.Info(fmt.Sprintf(s, a...), zf...)
}

func (l *Logger) Warnf(s string, a ...any) {
	zf := []zapcore.Field{}
	for _, f := range l.defaultFields {
		zf = append(zf, f.ToZap())
	}
	l.l.Warn(fmt.Sprintf(s, a...), zf...)
}

func (l *Logger) Errorf(s string, a ...any) {
	zf := []zapcore.Field{}
	for _, f := range l.defaultFields {
		zf = append(zf, f.ToZap())
	}
	l.l.Error(fmt.Sprintf(s, a...), zf...)
}

func (l *Logger) Debugf(s string, a ...any) {
	zf := []zapcore.Field{}
	for _, f := range l.defaultFields {
		zf = append(zf, f.ToZap())
	}
	l.l.Debug(fmt.Sprintf(s, a...), zf...)
}
