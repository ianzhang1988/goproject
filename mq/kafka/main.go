package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	. "github.com/segmentio/kafka-go"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func pub(brokers []string) {
	log.Println("pub ...")
	// Make a writer that publishes messages
	// The topic will be created if it is missing.
	w := &Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  "kafka-go-test",
		AllowAutoTopicCreation: true,
		BatchTimeout:           1 * time.Second,
		// BatchSize: 2,
		// Logger: log.Default(),
	}

	log.Printf("BatchTimeout %v\n", w.BatchTimeout)

	messages := []kafka.Message{
		{
			//Key:   []byte("Key-A"),
			Value: []byte("Hello World!"),
		},
		{
			//Key:   []byte("Key-B"),
			Value: []byte("One!"),
		},
		{
			//Key:   []byte("Key-C"),
			Value: []byte("Two!"),
		},
	}

	log.Println("send ...")
	var err error
	const retries = 3
	for i := 0; i < retries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// attempt to create topic prior to publishing the message
		err = w.WriteMessages(ctx, messages...)
		if errors.Is(err, LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Fatalf("unexpected error %v", err)
		} else {
			break
		}
	}

	for i := 0; i < 15; i++ {
		message := kafka.Message{
			//Key:   []byte("Key-A"),
			Value: []byte(fmt.Sprintf("num: %d", i)),
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		log.Printf("send num: %d\n", i)
		err = w.WriteMessages(ctx, message)
		if errors.Is(err, LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 250)
			continue
		}

		if err != nil {
			log.Fatalf("unexpected error %v", err)
		}

	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

	log.Println("send fin ...")
}

func sub(brokers []string) {
	log.Println("sub ...")
	// make a new reader that consumes messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "kafka-go-test-grp",
		Topic:    "kafka-go-test",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		// Logger:   log.Default(),
		// 下面两个变量影响read的时间，目前都1s就可以一个一个读数据。
		ReadBatchTimeout: 1 * time.Second,
		MaxWait:          1 * time.Second,
	})

	log.Println("read ...")
	for {
		ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
		// ctx := context.Background()
		m, err := r.ReadMessage(ctx)
		if err != nil {
			break
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	log.Println("read fin ...")
}

func batchPub(w *Writer) {
	log.Println("pub ...")

	// log.Printf("BatchTimeout %v\n", w.BatchTimeout)

	fill := func(in *[]byte, b byte) {
		for i, _ := range *in {
			(*in)[i] = b
		}
	}

	bytesA := make([]byte, 1024/2)
	fill(&bytesA, byte('A'))
	bytesB := make([]byte, 1024/2)
	fill(&bytesB, byte('B'))
	bytesC := make([]byte, 1024/2)
	fill(&bytesC, byte('C'))

	messages_template := []kafka.Message{
		{
			//Key:   []byte("Key-A"),
			Value: bytesA,
		},
		{
			//Key:   []byte("Key-B"),
			Value: bytesB,
		},
		{
			//Key:   []byte("Key-C"),
			Value: bytesC,
		},
	}
	template_len := len(messages_template)

	messages := []kafka.Message{}
	for i := 0; i < 1000; i++ {
		messages = append(messages, messages_template[i%template_len])
	}

	log.Println("send ...")
	var err error
	const retries = 3

	for c := 0; c < 1500; c++ {
		for i := 0; i < retries; i++ {
			ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
			// defer cancel()

			// attempt to create topic prior to publishing the message
			err = w.WriteMessages(ctx, messages...)
			if errors.Is(err, LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
				log.Printf("write err: %s, retry: %d\n", err, i)
				time.Sleep(time.Millisecond * 250)
				continue
			}

			if err != nil {
				log.Printf("unexpected error %v, retry %d\n", err, i)
				time.Sleep(5 * time.Second)
				continue
			} else {
				break
			}
		}
	}

	// if err := w.Close(); err != nil {
	// 	log.Fatal("failed to close writer:", err)
	// }

	log.Println("send fin ...")
}

func batchSub(brokers []string) {
	begin := time.Now()
	var end time.Time
	counter := 0
	log.Println("sub ...")

	// make a new reader that consumes messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  "kafka-go-test-grp",
		Topic:    "kafka-go-test",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		// Logger:   log.Default(),
		// 下面两个变量影响read的时间，目前都1s就可以一个一个读数据。
		// ReadBatchTimeout: 1 * time.Second,
		MaxWait:        1 * time.Second,
		QueueCapacity:  100,
		CommitInterval: 5 * time.Second,
	})

	log.Println("read ...")
	for {
		ctx, _ := context.WithTimeout(context.Background(), 60*time.Second)
		// ctx := context.Background()
		m, err := r.ReadMessage(ctx)
		if errors.Is(err, context.DeadlineExceeded) {
			log.Printf("time out\n")
			break
		}
		if err != nil {
			log.Printf("read err: %s\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		counter += len(m.Value)
		end = time.Now()
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}

	log.Printf("speed: %f", float64(counter)/1024/1024/float64(end.Sub(begin).Seconds()))

	log.Println("read fin ...")
}

func InitWriter(brokers []string) *Writer {
	// Make a writer that publishes messages
	// The topic will be created if it is missing.
	w := &Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  "kafka-go-test",
		AllowAutoTopicCreation: true,
		// BatchTimeout:           1 * time.Second,
		BatchSize:  100,
		BatchBytes: 1000000,
		// Logger: log.Default(),
		// Async: true,
	}
	return w
}

func CloseWrite(w *Writer) {
	log.Println("close writer ...")
	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile | log.Ltime | log.Ldate)
	addr := getEnv("KAFKA_ADDR", "localhost:9092")
	brokers := strings.Split(addr, ",")
	var w *Writer
	// w := InitWriter(brokers)
	// defer CloseWrite(w)
	// https://github.com/segmentio/kafka-go/issues/417
	for i := 0; i < 50; i++ {
		if i%50 == 0 {
			w = InitWriter(brokers)
			defer CloseWrite(w)
		}
		go batchPub(w)
	}

	// for i := 0; i < 19; i++ {
	// 	go batchSub(brokers)
	// }
	// batchSub(brokers)
	time.Sleep(300 * time.Second)
}
