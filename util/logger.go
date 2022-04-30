package util

import (
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {

	custLog := logrus.New()
	custLog.SetReportCaller(true)
	custLog.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename + ":" + strconv.Itoa(f.Line)

		},
	}

	return custLog
}
func getLoggerInstance() *logrus.Logger {
	return nil
}
