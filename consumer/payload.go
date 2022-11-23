package main

type Payload struct {
	Id         int    `json:"id"`
	ProducerId int    `json:"producer_id"`
	ConsumerId int    `json:"consumer_id"`
	Payload    string `json:"payload"`
}
