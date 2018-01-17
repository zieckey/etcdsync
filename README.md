# etcdsync

etcdsync is a distributed lock library in Go using etcd. It easy to use like sync.Mutex.


In fact, there are many similar implementation which are all obsolete 
depending on library `github.com/coreos/go-etcd/etcd` which is official marked `deprecated`,
and the usage is a little bit complicated. 
Otherwise this library is very very simple. The usage is simple, the code is simple.

## Import
    
    go get -u github.com/zieckey/etcdsync

## Simplest usage

Steps:

1. m, err := etcdsync.New()
2. m.Lock()
3. Do your business here
4. m.Unlock()

```go
package main

import (
	"log"
	"os"

	"github.com/zieckey/etcdsync"
)

func main() {
	m, err := etcdsync.New("/mylock", 10, []string{"http://127.0.0.1:2379"})
	if m == nil || err != nil {
		log.Printf("etcdsync.New failed")
		return
	}
	err = m.Lock()
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

`etcd2` or `etcd3` test OK.

You need a `etcd` instance running on http://localhost:2379, then:

    go test
    
