package external

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type (
	Setup struct {
		Env     string
		Logname string
		Display bool
		Level   string
	}
	MyFormatter struct{}
)

var levelList = []string{
	"PANIC",
	"FATAL",
	"ERROR",
	"WARN",
	"INFO",
	"DEBUG",
	"TRACE",
}

func (mf *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	level := levelList[int(entry.Level)]
	strList := strings.Split(entry.Caller.File, "/")
	fileName := strList[len(strList)-1]
	b.WriteString(fmt.Sprintf("%s - %s - [line:%d] - %s - %s\n",
		entry.Time.Format("2006-01-02 15:04:05,678"), fileName,
		entry.Caller.Line, level, entry.Message))
	return b.Bytes(), nil
}

func MakeLogger(s Setup) *logrus.Logger {
	f, err := os.OpenFile(s.Logname, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		panic(err.Error())
	}
	logger := logrus.New()
	if s.Display {
		logger.SetOutput(io.MultiWriter(os.Stdout, f))
	} else {
		logger.SetOutput(io.MultiWriter(f))
	}
	logger.SetReportCaller(true)

	if s.Env == "Production" {
		logger.SetFormatter(&logrus.JSONFormatter{})
	} else if s.Env == "Staging" {
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})
	}

	s.Level = strings.ToUpper(s.Level)
	if s.Level == "INFO" {
		logger.Level = logrus.InfoLevel
	} else if s.Level == "DEBUG" {
		logger.Level = logrus.DebugLevel
	} else if s.Level == "WARN" {
		logger.Level = logrus.WarnLevel
	} else if s.Level == "ERROR" {
		logger.Level = logrus.ErrorLevel
	} else if s.Level == "FATAL" {
		logger.Level = logrus.FatalLevel
	} else if s.Level == "PANIC" {
		logger.Level = logrus.PanicLevel
	} else if s.Level == "TRACE" {
		logger.Level = logrus.TraceLevel
	}

	/* rotatelog.NewHook("./"+filename+".%Y%m%d",
	rotatelog.WithMaxAge(24*time.Hour),
	rotatelog.WithRotationTime(time.Hour),
	rotatelog.WithClock(rotatelog.Local)) */

	/* if err != nil {
		log.Hooks.Add(hook)
	} */

	return logger
}
