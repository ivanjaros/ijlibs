// This is refactored copy from github.com/sirupsen/logrus/hooks/syslog
package syslogrus

import (
	"fmt"
	"github.com/hashicorp/go-syslog"
	"github.com/sirupsen/logrus"
	"os"
)

type SyslogHook struct {
	Writer gsyslog.Syslogger
}

func NewSyslogHook(logger gsyslog.Syslogger) *SyslogHook {
	return &SyslogHook{Writer: logger}
}

func (hook *SyslogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return hook.Writer.WriteLevel(gsyslog.LOG_CRIT, line)
	case logrus.FatalLevel:
		return hook.Writer.WriteLevel(gsyslog.LOG_CRIT, line)
	case logrus.ErrorLevel:
		return hook.Writer.WriteLevel(gsyslog.LOG_ERR, line)
	case logrus.WarnLevel:
		return hook.Writer.WriteLevel(gsyslog.LOG_WARNING, line)
	case logrus.InfoLevel:
		return hook.Writer.WriteLevel(gsyslog.LOG_INFO, line)
	case logrus.DebugLevel, logrus.TraceLevel:
		return hook.Writer.WriteLevel(gsyslog.LOG_DEBUG, line)
	default:
		return nil
	}
}

func (hook *SyslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
