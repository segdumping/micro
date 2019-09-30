package logic

import (
	"context"
	mon "github.com/micro/go-micro/monitor"
	"github.com/segdumping/micro/monitor/conf"
	registry2 "github.com/segdumping/micro/registry"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/micro/go-micro/client"
	pb "github.com/micro/go-micro/debug/proto"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/cache"
)

type monitor struct {
	options mon.Options

	exit     chan bool
	registry cache.Cache
	client   client.Client

	sync.RWMutex
	running  bool
	services map[string]*mon.Status
}

//debug.Client check
func (m *monitor) check(service string) (*mon.Status, error) {
	//service list
	services, err := m.registry.GetService(service)
	if err != nil {
		return nil, err
	}

	debug := pb.NewDebugService(service, m.client)

	var hasValid bool
	var checkErr error

	for _, service := range services {
		for _, node := range service.Nodes {
			if node.Metadata["server"] != m.client.String() {
				continue
			}

			//check health
			rsp, err := debug.Health(
				context.Background(),
				&pb.HealthRequest{},
				client.WithAddress(node.Address),
				client.WithRetries(3),
			)

			if err != nil {
				_ = m.registry.Deregister(&registry.Service{
					Name:    service.Name,
					Version: service.Version,
					Nodes:   []*registry.Node{node},
				})

				checkErr = err
				logrus.Errorf("node:[%s] health check error: %s", node.Id, err.Error())
				continue
			}

			// TODO
			if rsp.Status != "ok" {
				logrus.Infof("node:[%s] health check status: %s", node.Id, rsp.Status)
				continue
			}

			hasValid = true
		}
	}

	//ping ok
	if hasValid {
		return &mon.Status{Code: mon.StatusRunning, Info: "running"}, nil
	}

	//ping failed
	if checkErr != nil {
		return &mon.Status{
			Code:  mon.StatusFailed,
			Info:  "not running",
			Error: checkErr.Error(),
		}, nil
	}

	// otherwise unknown status
	return &mon.Status{
		Code: mon.StatusUnknown,
		Info: "unknown status",
	}, nil
}

//clean not exist service node
func (m *monitor) reap() {
	services, err := m.registry.ListServices()
	if err != nil {
		return
	}

	serviceMap := make(map[string]bool)
	for _, service := range services {
		serviceMap[service.Name] = true
	}

	m.Lock()
	defer m.Unlock()

	//delete??
	for service, _ := range m.services {
		if !serviceMap[service] {
			delete(m.services, service)
		}
	}
}

func (m *monitor) run() {
	conf := conf.Config.Micro
	checkTimer := time.NewTimer(time.Second)
	defer checkTimer.Stop()

	duration := time.Duration(conf.ReapDuration)
	reapTimer := time.NewTicker(duration * time.Second)
	defer reapTimer.Stop()

	check := make(chan string, 20)

	for {
		select {
		// just exit when told
		case <-m.exit:
			return
		case service := <-check:
			status, err := m.check(service)
			if err != nil {
				status = &mon.Status{
					Code: mon.StatusUnknown,
					Info: "unknown status",
				}
			}

			logrus.Debugf("check service %s, status: %v", service, status)

			m.Lock()
			m.services[service] = status
			m.Unlock()
		case <-checkTimer.C:
			// create a list of services
			serviceMap := make(map[string]bool)

			m.RLock()
			for service, _ := range m.services {
				serviceMap[service] = true
			}
			m.RUnlock()

			go func() {
				for service, _ := range serviceMap {
					select {
					//stop check, just exit
					case <-m.exit:
						return
					case check <- service:
					default:
						// barf if we block
					}
				}

				// check and watch new service
				services, _ := m.registry.ListServices()
				for _, service := range services {
					if ok := serviceMap[service.Name]; !ok {
						m.Watch(service.Name)
					}
				}
			}()

			duration := time.Duration(conf.CheckDuration)
			checkTimer.Reset(duration * time.Second)
		case <-reapTimer.C:
			m.reap()
		}
	}
}

//TODO
//notice when service status changed
func (m *monitor) notify() {

}

//deregister a service from registry
func (m *monitor) Reap(service string) error {
	services, err := m.registry.GetService(service)
	if err != nil {
		return nil
	}

	m.Lock()
	defer m.Unlock()
	delete(m.services, service)
	for _, service := range services {
		m.registry.Deregister(service)
	}

	return nil
}

//service status
func (m *monitor) Status(service string) (mon.Status, error) {
	m.RLock()
	defer m.RUnlock()
	if status, ok := m.services[service]; ok {
		return *status, nil
	}

	return mon.Status{}, mon.ErrNotWatching
}

// check or add
func (m *monitor) Watch(service string) error {
	m.Lock()
	defer m.Unlock()

	//check watched
	if _, ok := m.services[service]; ok {
		return nil
	}

	//check status
	status, err := m.check(service)
	if err != nil {
		return err
	}

	//add watch
	m.services[service] = status
	return nil
}

func (m *monitor) Run() error {
	m.Lock()
	defer m.Unlock()

	if m.running {
		return nil
	}

	m.exit = make(chan bool)
	m.registry = cache.New(m.options.Registry)

	go m.run()
	m.running = true

	return nil
}


func (m *monitor) Stop() error {
	m.Lock()
	defer m.Unlock()

	if !m.running {
		return nil
	}

	select {
	case <-m.exit:
		return nil
	default:
		close(m.exit)
		for s, _ := range m.services {
			delete(m.services, s)
		}
		m.registry.Stop()
		m.running = false
		return nil
	}

	return nil
}

func NewMonitor(opts ...mon.Option) mon.Monitor {
	options := mon.Options{
		Client:   client.DefaultClient,
		Registry: registry2.NewRegistry(),
	}

	for _, o := range opts {
		o(&options)
	}

	return &monitor{
		options:  options,
		exit:     make(chan bool),
		client:   options.Client,
		registry: cache.New(options.Registry),
		services: make(map[string]*mon.Status),
	}
}
