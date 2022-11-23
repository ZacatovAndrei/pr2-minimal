package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

var QueueAccess sync.Mutex

type Consumer struct {
	id int
}

func (c *Consumer) Start(id int, storage *list.List) {
	c.Init(id)
	log.Printf("Consumer #%v initialised", id)
	for {
		if p, ok := c.popResource(storage); ok {
			log.Println("Payload acquired, processing")
			p = c.process(p)
			c.Send(p)
		} else {
			log.Println("Buffer is empty, waiting")
		}
		time.Sleep(TimeUnit)
	}
}

func (c *Consumer) Init(id int) {
	c.id = id
}

func (c *Consumer) popResource(l *list.List) (Payload, bool) {
	QueueAccess.Lock()
	defer QueueAccess.Unlock()
	if l.Len() == 0 {
		return Payload{}, false
	}
	return (l.Remove(l.Front())).(Payload), true
}

func (c *Consumer) process(pl Payload) Payload {
	log.Println("processing item #", pl.Id)
	pl.ConsumerId = c.id
	pl.Payload = "PONG!"
	return pl
}

func (c *Consumer) Send(pl Payload) {
	log.Printf("updated payload:%v\n", pl)
	var serialised []byte
	serialised, ok := json.Marshal(pl)
	if ok != nil {
		log.Panicln(ok)
	}
	for {
		resp, ok := http.Post(aggregatorPort+"/ctp", "application/json", bytes.NewBuffer(serialised))
		if ok != nil {
			panic(ok)
		}
		if resp.StatusCode == 200 {
			resp.Body.Close()
			log.Println("Sent successfully")
			return
		}
		if resp.StatusCode == 507 {
			log.Println("Error 507,Aggregator buffer full,waiting")
			// closing response body because go doc said so
			if ok := resp.Body.Close(); ok != nil {
				panic(ok)
			}
			// wait a bit in hopes that the buffer will free up
			time.Sleep(TimeUnit)

		}
	}
}
