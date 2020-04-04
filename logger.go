package lunar

import "log"

// Logger is logger interface
type Logger interface {
	Printf(string, ...interface{})
}

// LoggerFunc is a bridge between Logger and any third party logger
type LoggerFunc func(string, ...interface{})

// Printf implements Logger interface
func (f LoggerFunc) Printf(msg string, args ...interface{}) { f(msg, args...) }

// defaultLogger writes nothing
var defaultLogger = LoggerFunc(func(string, ...interface{}) {})

// Printf is a logger which wraps log.Printf
var Printf = LoggerFunc(log.Printf)
