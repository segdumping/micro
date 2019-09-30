package log

import (
	"errors"
	"fmt"
	log2 "github.com/micro/go-micro/util/log"
	"github.com/sirupsen/logrus"
	"github.com/segdumping/shared/log"
	"os"
	"path"
	"runtime"
)

type Config struct {
	Level   string `xml:"level"`
	File    string `xml:"file"`
	Path    string `xml:"path"`
}

func (c *Config) LogPath() string {
	return c.Path + "/" + c.File
}

//replace go-micro log with logrus
//after call this method,  just direct call logrus
func Init(conf *Config) error {
	var levelMapping = map[string]log2.Level{
		"trace": log2.LevelTrace,
		"debug": log2.LevelDebug,
		"info":  log2.LevelInfo,
		"error": log2.LevelError,
		"fatal": log2.LevelFatal,
	}

	l := log.New()
	//std logger
	l.SetLogger(logrus.StandardLogger())
	//redirect
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("runtime caller error")
	}

	path := path.Dir(f) + "/" + conf.LogPath()
	err := l.Redirect(path)
	if err != nil {
		return err
	}

	l.SetOutput(os.Stdout)
	l.SetLevel(conf.Level)
	level, ok := levelMapping[conf.Level]
	if !ok {
		return errors.New(fmt.Sprintf("config level: %s error", conf.Level))
	}

	//go-micro log
	log2.SetLevel(level)
	log2.SetLogger(l)

	return nil
}
