package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

// custom Log override Log and Logf of logrus in order to replace
// micro default log with logrus
var std *Log

func init() {
	std = New()
}

type Log struct {
	*logrus.Logger
	*CustomFormatter
	level   logrus.Level
}

func New() *Log {
	l := &Log{
		Logger: logrus.New(),
		level:  logrus.InfoLevel,
		CustomFormatter: &CustomFormatter{},
	}

	l.SetFormatter(l.CustomFormatter)

	return l
}

func (l *Log) SetLogger(logger *logrus.Logger) {
	l.Logger = logger
}

func (l *Log) Log(v ...interface{}) {
	l.Logger.Log(l.level, v...)
}

func (l *Log) Logf(format string, v ...interface{}) {
	l.Logger.Logf(l.level, format, v...)
}

func (l *Log) SetLevel(level string) {
	lv, ok := levelMapping[level]
	if !ok {
		return
	}

	l.Logger.SetLevel(lv)
}

func (l *Log) Redirect(file string) error {
	return redirect(l, file)
}

func SetLevel(level string) {
	std.SetLevel(level)
}

func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

func Redirect(file string) error {
	return redirect(std, file)
}

func redirect(l *Log, file string) error {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	l.SetOutput(f)
	return nil
}

//TODO
func rotate() {

}

var levelMapping = map[string]logrus.Level{
	"debug": logrus.DebugLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"error": logrus.ErrorLevel,
	"trace": logrus.TraceLevel,
	"fatal": logrus.FatalLevel,
	"panic": logrus.PanicLevel,
}
