package conf

//server general config
type Config struct {
	Id               string `xml:"id"`
	Name             string `xml:"name"`
	Addr             string `xml:"address"`
	Namespace        string `xml:"namespace"`
	RegisterTTL      int    `xml:"registerTTL"`
	RegisterInterval int    `xml:"registerInterval"`
}

func (c *Config) ServerName() string {
	return c.Namespace + "." + c.Name
}
