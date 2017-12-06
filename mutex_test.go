package etcdsync

import (
	"log"
	"testing"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

func newKeysAPI(machines []string) client.KeysAPI {
	cfg := client.Config{
		Endpoints:               machines,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(cfg)
	if err != nil {
		return nil
	}

	return client.NewKeysAPI(c)
}

func checkKeyExists(key string, kapi client.KeysAPI) bool {
	// Get the already node's value.
	_, err := kapi.Get(context.TODO(), key, nil)
	if err != nil {
		return false
	}
	return true
}

func TestMutex(t *testing.T) {
	log.SetFlags(log.Ltime | log.Ldate | log.Lshortfile)
	lockKey := "/etcdsync"
	machines := []string{"http://127.0.0.1:2379"}
	kapi := newKeysAPI(machines)
	m, err := New(lockKey, 60, machines)
	if err != nil {
		t.Error(err)
		return
	}
	if m == nil {
		t.Errorf("New Mutex ERROR")
	}
	err = m.Lock()
	if err != nil {
		t.Errorf("failed")
	}

	if checkKeyExists(lockKey, kapi) == false {
		t.Errorf("The mutex have been locked but the key node does not exists.")
		t.Fail()
	}
	//do something here

	err = m.Unlock()
	if err != nil {
		t.Errorf("failed")
	}

	_, err = m.kapi.Get(context.Background(), lockKey, nil)
	if e, ok := err.(client.Error); !ok {
		t.Errorf("Get key %v failed from etcd", lockKey)
	} else if e.Code != client.ErrorCodeKeyNotFound {
		t.Errorf("ERROR %v", err)
	}
}

func TestLockConcurrently(t *testing.T) {
	slice := make([]int, 0, 3)
	lockKey := "/etcd_sync"
	machines := []string{"http://127.0.0.1:2379"}
	kapi := newKeysAPI(machines)
	m1, err := New(lockKey, 60, machines)
	m2, err := New(lockKey, 60, machines)
	m3, err := New(lockKey, 60, machines)
	if err != nil {
		t.Error(err)
		return
	}
	if m1 == nil || m2 == nil || m3 == nil {
		t.Errorf("New Mutex ERROR")
	}
	m1.Lock()
	if checkKeyExists(lockKey, kapi) == false {
		t.Errorf("The mutex have been locked but the key node does not exists.")
		t.Fail()
	}
	ch1 := make(chan bool)
	go func() {
		ch2 := make(chan bool)
		m2.Lock()
		if checkKeyExists(lockKey, kapi) == false {
			t.Errorf("The mutex have been locked but the key node does not exists.")
			t.Fail()
		}
		go func() {
			m3.Lock()
			if checkKeyExists(lockKey, kapi) == false {
				t.Errorf("The mutex have been locked but the key node does not exists.")
				t.Fail()
			}
			slice = append(slice, 2)
			m3.Unlock()
			ch2 <- true
		}()
		slice = append(slice, 1)
		time.Sleep(1 * time.Second)
		m2.Unlock()
		<-ch2
		ch1 <- true
	}()
	slice = append(slice, 0)
	time.Sleep(1 * time.Second)
	m1.Unlock()
	<-ch1
	if len(slice) != 3 {
		t.Fail()
	}
	for n, i := range slice {
		if n != i {
			t.Fail()
		}
	}
}

func TestLockTimeout(t *testing.T) {
	slice := make([]int, 0, 2)
	m1, err := New("key", 2, []string{"http://127.0.0.1:2379"})
	m2, err := New("key", 2, []string{"http://127.0.0.1:2379"})
	if err != nil {
		t.Error(err)
		return
	}
	m1.Lock()
	ch := make(chan bool)
	go func() {
		m2.Lock()
		slice = append(slice, 1)
		m2.Unlock()
		ch <- true
	}()
	slice = append(slice, 0)
	<-ch
	for n, i := range slice {
		if n != i {
			t.Fail()
		}
	}
}

func TestRefreshLockTTL(t *testing.T) {
	lockKey := "/etcd_sync"
	machines := []string{"http://127.0.0.1:2379"}
	kapi := newKeysAPI(machines)
	m, err := New(lockKey, 10, machines)
	if err != nil {
		t.Error(err)
		return
	}
	m.Lock()
	if checkKeyExists(lockKey, kapi) == false {
		t.Errorf("The mutex have been refreshed but the key node does not exists.")
		t.Fail()
	}
	time.Sleep(5 * time.Second)
	m.RefreshLockTTL(10 * time.Second)
	time.Sleep(5 * time.Second)
	if checkKeyExists(lockKey, kapi) == false {
		t.Errorf("The mutex's TTL has been refreshed but the key node still expired.")
		t.Fail()
	}
	m.Unlock()
}
