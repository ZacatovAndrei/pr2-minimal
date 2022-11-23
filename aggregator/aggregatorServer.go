package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	consumerAddress   = "http://localhost:8088/receive"
	producerAddress   = "http://localhost:8086/receive"
	aggregatorPort    = ":8087"
	TimeUnit          = 2 * time.Second
	StorageSize       = 10
	AggregatorThreads = 5
)

var (
	ProducerStack = list.New()
	ConsumerStack = list.New()
)

func init() {
	time.Sleep(2 * time.Second)
	fmt.Println("starting producer server on port 8086")
}

func main() {
	for i := 0; i < AggregatorThreads; i++ {
		go aggregator(ConsumerStack, ProducerStack)
	}
	http.HandleFunc("/ctp", storeFinished)
	http.HandleFunc("/ptc", storeNew)
	http.ListenAndServe(aggregatorPort, nil)
}

func storeNew(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not allowed")
		return
	}
	if ConsumerStack.Len() > StorageSize {
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
	log.Printf("pushing %v to the consumer stack", o)
	ConsumerStack.PushBack(o)
}

func storeFinished(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		fmt.Fprintf(w, "Method not allowed")
		return
	}
	if ProducerStack.Len() > StorageSize {
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
	if ok != nil {
		panic(ok)
	}
	log.Printf("Received %v a thing from consumer", o)
	ProducerStack.PushBack(o)
}
