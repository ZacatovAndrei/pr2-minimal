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
	consumerPort    = ":8088"
	producerPort    = "http://localhost:8086"
	aggregatorPort  = "http://localhost:8087"
	TimeUnit        = 2 * time.Second
	ConsumerThreads = 5
	StorageSize     = 10
)

var (
	Consumers = make([]Consumer, ConsumerThreads)
	Storage   = list.New()
)

func init() {
	time.Sleep(2 * time.Second)
	fmt.Println("starting consumer server on port 8088")
}

func main() {
	for i := 0; i < ConsumerThreads; i++ {
		go Consumers[i].Start(i, Storage)
	}
	http.HandleFunc("/receive", receive)
	if ok := http.ListenAndServe(consumerPort, nil); ok != nil {
		panic(ok)
	}

}

// function to receive stuff from aggregator
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
