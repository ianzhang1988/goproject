package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"time"
	"flag"
)

var (
	BrokerAddress = flag.StringVar("kafka")
)

func main() {
	// Kafka broker configuration
	brokerAddress := "your_kafka_broker_address"
	topic := "your_topic"

	// Create a new Kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		MinBytes: 10e3, // Minimum number of bytes to fetch from Kafka
		MaxBytes: 10e6, // Maximum number of bytes to fetch from Kafka
	})

	// Get the partitions for the topic
	partitions, err := reader.ReadInfo(context.Background())
	if err != nil {
		panic(err)
	}

	// Convert desired time to Unix timestamp
	desiredTime := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC).Unix()

	// Set the offset for each partition to the earliest offset whose timestamp is greater than or equal to the desired time
	for _, partition := range partitions {
		_, offset, err := reader.Seek(kafka.SeekByTime(desiredTime), kafka.Assigned(partition.ID))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Partition %d offset set to %d\n", partition.ID, offset)
	}

	// Start consuming messages
	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			panic(err)
		}

		fmt.Printf("Received message from partition %d: %s\n", message.Partition, string(message.Value))
	}
}
