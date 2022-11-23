package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	CSAccess sync.Mutex
	PSAccess sync.Mutex
)

func aggregator(Cl *list.List, Pl *list.List) {
	rand.Seed(time.Now().Unix())
	var o Payload
	var itemFound bool
	for {
		log.Printf("ConStack:%v\tProdStack:%v", ConsumerStack.Len(), ProducerStack.Len())
		itemFound = false
		// checking one of the stacks randomly
		state := rand.Intn(2)
		// Taking random Element from top
		switch state {
		case 0:
			o, itemFound = getResource(Cl, &CSAccess)
		case 1:
			o, itemFound = getResource(Pl, &PSAccess)
		}
		if !itemFound {
			log.Println("Buffer is empty, waiting")
			time.Sleep(TimeUnit / 2)
			continue
		}
		// Send to the appropriate server
		switch state {
		case 0:
			sendTo(consumerAddress, o)
		case 1:
			sendTo(producerAddress, o)
		}
		time.Sleep(TimeUnit)
	}

}

func getResource(l *list.List, m *sync.Mutex) (Payload, bool) {
	m.Lock()
	defer m.Unlock()
	if l.Len() == 0 {
		log.Println("The buffer is empty")
		return Payload{}, false
	}
	res := l.Remove(l.Front())
	return res.(Payload), true
}

func sendTo(address string, pl Payload) {
	serialised, ok := json.Marshal(pl)
	if ok != nil {
		panic(ok)
	}
	for {
		resp, ok := http.Post(address, "application/json", bytes.NewBuffer(serialised))
		if ok != nil {
			panic(ok)
		}
		if resp.StatusCode == 200 {
			// if sent successfully then close the body because doc says so
			if ok := resp.Body.Close(); ok != nil {
				panic(ok)
			}
			// nothing more to do. return
			return
		}
		// otherwise there is some issue
		if resp.StatusCode == 507 {
			log.Println("Error 507,Server buffer is full,waiting")
			// closing response body because go doc said so
			if ok := resp.Body.Close(); ok != nil {
				panic(ok)
			}
			// wait a bit in hopes that the buffer will free up
			time.Sleep(TimeUnit)
		}
	}
}
