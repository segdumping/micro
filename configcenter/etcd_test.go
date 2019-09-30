package configcenter

import "testing"

func TestEtcd(t *testing.T) {
	c := NewConfigCenter()
	p, err := c.Put("invalid-services", "go.micro.api.echo-1")
	if err != nil {
		t.Log("put error: ", err)
		return
	}

	t.Log("put response: ", p)

	v, err := c.Get("invalid-service")
	if err != nil {
		t.Log("get error: ", err)
		return
	}

	t.Log("get response", v)
}
