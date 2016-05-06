package logger

// LogFormatterFactory defines methods a factory should implement, to create a Logger
// based on the specified name.
type LogFormatterFactory interface {
	Create(name string) LogFormatter
}

// LogFormatter defines methods to implement for being a compose log logger.
type LogFormatter interface {
	Out(bytes []byte)
	Err(bytes []byte)
}

// Logger defines methods to implement for being a logger.
type Logger interface {
	Errorf(format string, values ...interface{})
	Error(values ...interface{})
	Infof(format string, values ...interface{})
	Info(values ...interface{})
	Debugf(format string, values ...interface{})
	Debug(values ...interface{})
	Warnf(format string, values ...interface{})
	Warn(values ...interface{})
}

// Wrapper is a wrapper around Logger that implements the Writer interface,
// mainly use by docker/pkg/stdcopy functions.
type Wrapper struct {
	Err    bool
	Logger LogFormatter
}

func (l *Wrapper) Write(bytes []byte) (int, error) {
	if l.Err {
		l.Logger.Err(bytes)
	} else {
		l.Logger.Out(bytes)
	}
	return len(bytes), nil
}
