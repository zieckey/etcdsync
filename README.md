# go-etcd-lock

A distributed lock library in Go using etcd. It easy to use like sync.Mutex.


In fact, there are many similar implementation, but which are all obsolete 
depending on library `github.com/coreos/go-etcd/etcd` which is official marked `deprecated`,
and the usage is a little bit complicated. 
Otherwise this library is very very simple. The usage is simple, the code is simple.

## Import
    
    go get github.com/zieckey/go-etcd-lock

## Simplest usage

Steps:

1. m := etcdsync.New()
2. m.Lock()
3. Do your business here
4. m.Unlock()

```go
package main

import (
	"github.com/zieckey/go-etcd-lock"
	"log"
)

func main() {
	log.SetFlags(log.Ldate|log.Ltime|log.Lshortfile)
	m := etcdsync.New("/etcdsync", "123", []string{"http://127.0.0.1:2379"})
	if m == nil {
		log.Printf("etcdsync.NewMutex failed")
	}
	err := m.Lock()
	if err != nil {
		log.Printf("etcdsync.Lock failed")
	} else {
		log.Printf("etcdsync.Lock OK")
	}

	log.Printf("Get the lock. Do something here.")

	err = m.Unlock()
	if err != nil {
		log.Printf("etcdsync.Unlock failed")
	} else {
		log.Printf("etcdsync.Unlock OK")
	}
}

```

## Test

You need a etcd instance running on localhost:2379, then:

    go test
    