package configcenter

import (
	"context"
	"crypto/tls"
	"github.com/coreos/etcd/clientv3"
	"github.com/micro/go-micro/util/log"
	"github.com/segdumping/shared/util"
	"strings"
	"time"
)

const (
	prefix = "/micro/config/center/"
)

type config struct {
	options   Options
	client    *clientv3.Client
	eventChan chan *event
}

//watch channel transfer
type event struct {
	action  string
	key     string
	value   []byte
	version int64
}

func (e *config) Init(opts ...Option) error {
	return configure(e, opts...)
}

func (e *config) String() string {
	return "configCenter"
}

// get kv from etcd
func (e *config) Get(key string) (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	newKey := prefix + key
	r, err := e.client.KV.Get(ctx, newKey, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	//return map, maybe more than one
	m := make(map[string][]byte, len(r.Kvs))
	for _, v := range r.Kvs {
		k := util.Byte2String(v.Key)
		k = strings.TrimPrefix(k, prefix)
		m[k] = v.Value
	}

	return m, nil
}

// store kv to etcd
func (e *config) Put(key, val string) (map[string][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	newKey := prefix + key
	r, err := e.client.KV.Put(ctx, newKey, val, clientv3.WithPrevKV())
	if err != nil {
		return nil, err
	}

	if r.PrevKv == nil {
		return nil, nil
	}

	k := util.Byte2String(r.PrevKv.Key)
	k = strings.TrimPrefix(k, prefix)

	//return previous kv
	return map[string][]byte{k: r.PrevKv.Value}, nil
}

//etcd watch
func (e *config) watch() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := e.client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithPrevKV())
	for r := range w {
		if r.Err() != nil {
			log.Logf("watch error: %s", r.Err().Error())
			return
		}

		for _, ev := range r.Events {
			w := &event{}
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					w.action = "create"
				}

				if ev.IsModify() {
					w.action = "modify"
				}

				w.key = util.Byte2String(ev.Kv.Key)
				w.key = strings.TrimPrefix(w.key, prefix)
				w.value = ev.Kv.Value
				w.version = ev.Kv.Version
			case clientv3.EventTypeDelete:
				w.action = "delete"
				w.key = util.Byte2String(ev.PrevKv.Key)
				w.key = strings.TrimPrefix(w.key, prefix)
				w.value = ev.PrevKv.Value
				w.version = ev.PrevKv.Version
			}

			e.eventChan <- w
		}
	}
}

//配置选项
func configure(e *config, opts ...Option) error {
	config := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},	//etcd默认地址
	}

	for _, o := range opts {
		o(&e.options)
	}

	//operate timeout
	if e.options.Timeout == 0 {
		e.options.Timeout = 5 * time.Second
	}

	//tls connect
	if e.options.Secure || e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		config.TLS = tlsConfig
	}

	var cAddrs []string
	for _, addr := range e.options.Addrs {
		if len(addr) == 0 {
			continue
		}
		cAddrs = append(cAddrs, addr)
	}

	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	cli, err := clientv3.New(config)
	if err != nil {
		return err
	} else {
		e.client = cli
	}

	return nil
}

func NewConfigCenter(opts ...Option) *config {
	e := &config{
		options:   Options{},
		eventChan: make(chan *event, 10),
	}

	configure(e, opts...)
	return e
}
