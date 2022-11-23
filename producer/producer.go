package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	PayloadNum  atomic.Int32
	QueueAccess sync.Mutex
)

type Producer struct {
	id int
}

func (p *Producer) Start(id int, storage *list.List) {
	var pl Payload
	p.Init(id)
	log.Printf("Producer #%v initialised", id)
	for {
		if resp, ok := p.receive(storage); ok {
			log.Printf("received a %v from consumer #%v", resp.Payload, resp.ConsumerId)
			time.Sleep(TimeUnit / 2)
			continue
		}
		pl = p.generate()
		p.send(pl)
		time.Sleep(TimeUnit)
	}
}

func (p *Producer) Init(id int) {
	p.id = id
}

func (p *Producer) receive(l *list.List) (Payload, bool) {
	QueueAccess.Lock()
	defer QueueAccess.Unlock()
	if l.Len() == 0 {
		return Payload{}, false
	}
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(Payload).ProducerId == p.id {
			res := l.Remove(e)
			log.Printf("Received payload %v back!", res.(Payload).Id)
			return res.(Payload), true
		}
	}
	return Payload{}, false
}

func (p *Producer) generate() (res Payload) {
	pid := PayloadNum.Add(1)
	return Payload{
		Id:         int(pid),
		ProducerId: p.id,
		Payload:    "PING?",
	}
}

func (p *Producer) send(pl Payload) {
	serialised, ok := json.Marshal(pl)
	if ok != nil {
		panic(ok)
	}
	for {
		resp, ok := http.Post(aggregatorAddress+"/ptc", "application/json", bytes.NewBuffer(serialised))
		if ok != nil {
			panic(ok)
		}
		if resp.StatusCode == 200 {
			// if sent successfully then close the body because doc says so
			if ok := resp.Body.Close(); ok != nil {
				panic(ok)
			}
			log.Printf("succesfully sent payload %v", pl.Id)
			// nothing more to do. return
			return
		}
		// otherwise there is some issue
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
