package conf

import "errors"

type Config struct {
	RegList []Registry `xml:"regList>registry" json:"regList"`
}

func (c *Config)GetValid() (*Registry, error) {
	for _, v := range c.RegList {
		if v.Valid {
			return &v, nil
		}
	}

	return nil, errors.New("no valid registry center")
}

type Registry struct {
	Type     string    `xml:"type" json:"type"`
	Valid    bool      `xml:"valid" json:"valid"`
	Endpoint Endpoint `xml:"endpoint" json:"endpoint"`
}

type Endpoint struct {
	Addr string `xml:"address" json:"address"`
	Port string `xml:"port" json:"port"`
}
