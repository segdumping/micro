package log

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"runtime"
)

type CustomFormatter struct {
	host    string
}

func (c *CustomFormatter) Format(e *logrus.Entry) ([]byte, error) {
	if len(c.host) == 0 {
		c.host, _ = os.Hostname()
	}


	m := make(map[string]interface{}, 3)
	m["host"] = c.host
	m["level"] = e.Level.String()
	m["msg"] = e.Message
	m["datatime"] = e.Time.Format("2006-01-02 15:04:05")

	_, file, line, ok := runtime.Caller(6)
	if ok {
		m["file"] = file
		m["line"] = line
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	} else {
		b = append(b, '\n')
		return b, nil
	}
}
