package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"debug":             "broker,topic,msg",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			default:
				fmt.Println("default")
			}
		}
	}()

	topic := "message-log"
	/*
		for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          []byte(word),
			}, nil)
		}
	*/

	for _, word := range []string{"Welcome", "to", "the", "Confluent", "Kafka", "Golang", "client"} {
		p.ProduceChannel() <- &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}
	}
  
	go func() {
		for {
			fmt.Println(p.Len())
			if p.Len() == 0 {
				sigchan <- os.Interrupt
				fmt.Println("send sigterm")
			}
			// p.Flush(2000)
			time.Sleep(time.Millisecond * 500)
		}
	}()

	<-sigchan
}
