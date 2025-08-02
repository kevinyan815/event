package event

import (
	"context"
	"fmt"
	"strings"
)

// LogLevel 定义日志级别
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

// LogParam
type LogParam struct {
	Key   string
	Value interface{}
}

// ILogger
type ILogger interface {
	Debug(ctx context.Context, msg string, params ...LogParam)
	Info(ctx context.Context, msg string, params ...LogParam)
	Warn(ctx context.Context, msg string, params ...LogParam)
	Error(ctx context.Context, msg string, params ...LogParam)
}

// Logger instance
var _Logger ILogger = &defaultLogger{}

// SetLogger set custom logger
func SetLogger(logger ILogger) {
	if logger != nil {
		_Logger = logger
	}
}

// GetLogger
func GetLogger() ILogger {
	return _Logger
}

// helper function to create log param
func P(key string, value interface{}) LogParam {
	return LogParam{Key: key, Value: value}
}

// default logger implementation
type defaultLogger struct{}

func formatParams(params []LogParam) string {
	if len(params) == 0 {
		return ""
	}

	var parts []string
	for _, p := range params {
		parts = append(parts, fmt.Sprintf("%q, %v", p.Key, p.Value))
	}
	return ", " + strings.Join(parts, ", ")
}

func (l *defaultLogger) Debug(ctx context.Context, msg string, params ...LogParam) {
}

func (l *defaultLogger) Info(ctx context.Context, msg string, params ...LogParam) {
	fmt.Printf("[INFO] %s%s\n", msg, formatParams(params))
}

func (l *defaultLogger) Warn(ctx context.Context, msg string, params ...LogParam) {
	fmt.Printf("[WARN] %s%s\n", msg, formatParams(params))
}

func (l *defaultLogger) Error(ctx context.Context, msg string, params ...LogParam) {
	fmt.Printf("[ERROR] %s%s\n", msg, formatParams(params))
}
