package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	consumerAddress   = "http://localhost:8088"
	producerPort      = "localhost:8086"
	aggregatorAddress = "http://localhost:8087"
	TimeUnit          = 2 * time.Second
	ProducerThreads   = 5
	StorageSize       = 10
)

var (
	Storage   = list.New()
	Producers = make([]Producer, ProducerThreads)
)

func init() {
	time.Sleep(2 * time.Second)
	fmt.Printf("starting producer server on port 8086\nThere are %v Producer threads\n", ProducerThreads)
}

func main() {
	for i := 0; i < ProducerThreads; i++ {
		go Producers[i].Start(i, Storage)
	}
	http.HandleFunc("/receive", receive)
	if ok := http.ListenAndServe(producerPort, nil); ok != nil {
		return
	}
}

func receive(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not allowed")
		return
	}
	if Storage.Len() >= StorageSize {
		w.WriteHeader(507)
		fmt.Fprintf(w, "The buffer is full")
		return
	}
	var o Payload
	body, ok := io.ReadAll(r.Body)
	if ok != nil {
		panic(ok)
	}
	ok = json.Unmarshal(body, &o)
	Storage.PushBack(o)
}
