package etcdsync

import (
	"log"
	"sync"
	"time"

	"github.com/coreos/etcd/client"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"fmt"
)

const (
	defaultTTL = 60
	defaultTry = 3
)

type Mutex struct {
	key    string
	id     string
	client client.Client
	kapi   client.KeysAPI
	state  int32
	mutex  *sync.Mutex
}

func New(key string, id string, machines []string) *Mutex {
	cfg := client.Config{
		Endpoints:               machines,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil
	}


	return &Mutex{
		key:    key,
		id:     id,
		client: c,
		kapi:   client.NewKeysAPI(c),
		mutex:  new(sync.Mutex),
	}
}

func (m *Mutex) Lock() (err error) {
	m.mutex.Lock()
	for try := 1; try <= defaultTry; try++ {
		debugf("id=%v Try to creating a node : key=%v", m.id, m.key)
		_, err = m.kapi.Create(context.Background(), m.key, m.id)
		if err == nil {
			debugf("id=%v Create node %v OK", m.id, m.key)
			return nil
		}
		debugf("id=%v create ERROR: %v", m.id, err)
		if e, ok := err.(client.Error); ok {
			if e.Code != client.ErrorCodeNodeExist {
				continue
			}

			wait:
			// Get the already node's value.
			resp, err := m.kapi.Get(context.Background(), m.key, nil)
			if err != nil {
				// Always try.
				try--
				continue
			}
			debugf("id=%v get ok", m.id)
			ops := &client.WatcherOptions {
				AfterIndex : resp.Index,
				Recursive:false,
			}
			watcher := m.kapi.Watcher(m.key, ops)
			for {
				debugf("id=%v watching ...", m.id)
				resp, err := watcher.Next(context.Background())
				if err != nil {
					goto wait
				}
				debugf("id=%v Received an event : %q", m.id, resp)
				break
			}
			debugf("id=%v watch done", m.id)
			continue // try to creating the lock node again
		}
	}
	return err
}

// Unlock
func (m *Mutex) Unlock() (err error) {
	defer m.mutex.Unlock()
	for i := 1; i <= defaultTry; i++ {
		var resp *client.Response
		resp, err = m.kapi.Delete(context.Background(), m.key, nil)
		debugf("id=%v Delete %q", m.id, resp)
		if err != nil {
			if _, ok := err.(client.Error); !ok {
				// retry.
				continue
			}
		}
		break
	}
	return
}

var debug = false
func debugf(format string, v ...interface{}) {
	if debug {
		log.Output(2, fmt.Sprintf(format, v...))
	}
}
func SetDebug(flag bool) {
	debug = flag
}

