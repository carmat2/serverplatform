package server

import (
	"os"
	"github.com/mattermost/logr/v2"
	"github.com/mattermost/logr/v2/formatters"
	"github.com/mattermost/logr/v2/targets"
)

var lgr *logr.Logr
var logger logr.Sugar

// CreateLogger creates the logger targets.
// For unit tests log output is sent to os.Stdout from the Trace level.
// For application run log output is sent to ./logs/serverplatform.log from the Info level.
// Both targets will log the stack trace from the Error level.
func CreateLogger(isTesting bool) {
	lgr,_ = logr.New()
	logger = lgr.NewLogger().Sugar()

	if isTesting {
		formatterStdOut := &formatters.Plain{DisableTimestamp: true, Delim: " | "}
		filterStdOut := &logr.StdFilter{Lvl: logr.Trace, Stacktrace: logr.Error}
		tStdOut := targets.NewWriterTarget(os.Stdout)
		lgr.AddTarget(tStdOut, "StdOut", filterStdOut, formatterStdOut, 1000)
	} else {
		formatterLogFile := &formatters.JSON{}
		filtersLogFile := &logr.StdFilter{Lvl: logr.Info, Stacktrace: logr.Error}
		// max file size 10MB, keep log files for 30 days, keep up to 5 old backup log files, no file compression
		opts := targets.FileOptions{
			Filename:   "./logs/serverplatform.log",
			MaxSize:    10,
			MaxAge:     30,
			MaxBackups: 5,
			Compress:   false,
		}
		tLogFile := targets.NewFileTarget(opts)
		lgr.AddTarget(tLogFile, "LogFile", filtersLogFile, formatterLogFile, 1000)
	}
}

// ShutdownLogger ensures targets are drained before application exit.
func ShutdownLogger() {
	lgr.Shutdown()
}

