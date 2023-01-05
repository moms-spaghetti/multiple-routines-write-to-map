package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

func main() {
	ch := make(chan map[int]interface{})
	done := make(chan struct{})
	wait := make(chan os.Signal, 1)
	signal.Notify(wait, syscall.SIGINT)

	store := createStore()

	go func() {
		for {
			d := time.Duration(rand.Intn(3))
			time.Sleep(d * time.Second)
			ch <- map[int]interface{}{rand.Intn(40): "data from routine 1"}
		}
	}()

	go func() {
		for {
			d := time.Duration(rand.Intn(3))
			time.Sleep(d * time.Second)
			ch <- map[int]interface{}{rand.Intn(40): "data from routine 2"}
		}
	}()

	go func() {
		for {
			d := time.Duration(rand.Intn(3))
			time.Sleep(d * time.Second)
			ch <- map[int]interface{}{rand.Intn(40): "data from routine 3"}
		}
	}()

	go func() {
		for {
			d := time.Duration(rand.Intn(3))
			time.Sleep(d * time.Second)
			ch <- map[int]interface{}{rand.Intn(40): "data from routine 4"}
		}
	}()

	go func() {
		for {
			select {
			case m := <-ch:
				for k, v := range m {
					store.add(k, v)
				}
			case <-done:
				return
			}
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Second)
			store.get()
		}
	}()

	<-wait
	close(wait)

}

type storage struct {
	store map[int]interface{}
}

func createStore() *storage {
	store := map[int]interface{}{}

	return &storage{
		store: store,
	}
}

func (s *storage) add(key int, value interface{}) {
	mutex := sync.Mutex{}

	mutex.Lock()
	s.store[key] = value
	mutex.Unlock()
}

func (s *storage) get() {
	var keys []int

	for k := range s.store {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	fmt.Printf("%s\n", time.Now().Format(time.RFC1123))
	for _, k := range keys {
		fmt.Printf("%d : %v\n", k, s.store[k])
	}
	fmt.Println("")
}
