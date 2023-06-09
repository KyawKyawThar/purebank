package worker

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)




type Logger struct {
}

func NewLoger() *Logger {
	//return empty logger struct cos it doesn't need to store any data
	return &Logger{}
}

func (l *Logger) PrintLog(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.PrintLog(zerolog.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.PrintLog(zerolog.InfoLevel, args...)
}
func (l *Logger) Warn(args ...interface{}) {
	l.PrintLog(zerolog.WarnLevel, args...)
}
func (l *Logger) Error(args ...interface{}) {
	l.PrintLog(zerolog.ErrorLevel, args...)
}
func (l *Logger) Fatal(args ...interface{}) {
	l.PrintLog(zerolog.FatalLevel, args...)
}
