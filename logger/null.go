package logger

// NullLogger is a logger.Logger and logger.Factory implementation that does nothing.
type NullLogger struct {
}

// Out is a no-op function.
func (n *NullLogger) Out(_ []byte) {
}

// Err is a no-op function.
func (n *NullLogger) Err(_ []byte) {
}

// Create implements logger.LogFormatterFactory and returns a NullLogger.
func (n *NullLogger) Create(_ string) LogFormatter {
	return &NullLogger{}
}
