package main

import (
	"github.com/zieckey/go-etcd-lock"
	"log"
	"os"
)

func main() {
	//logger := log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	m := etcdsync.New("/mylock", 10, []string{"http://127.0.0.1:2379"})
	m.SetDebugLogger(os.Stdout)
	if m == nil {
		log.Printf("etcdsync.New failed")
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
