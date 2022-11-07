package infra

import (
	"context"
	"github.com/linkedin/goavro"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"sync"
	"time"
)

type Writer struct {
	conn *kafka.Conn
}

type Reader struct {
	conn *kafka.Conn
}

var busCodec, _ = goavro.NewCodec(`{
		"type": "record",
		"name": "temp",
		"fields": [
			{
				"name": "name",
				"type": "string"
			},
			{
				"name": "count",
				"type": "long"
			},
			{
				"name": "location",
				"type": "string"
			},
			{
				"name": "user_id",
				"type" "long"
			},
		]
}`)

var lock = &sync.Mutex{}
var err error

func (w *Writer) GetKafkaWriter(port string, topic string, partition int) *Writer {
	if w.conn != nil {
		return w
	}
	//lock.Lock()
	//defer lock.Unlock()
	HOST := os.Getenv("KAFKA_HOST")
	w.conn, err = kafka.DialLeader(context.Background(), "tcp", HOST+":"+port, topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	return w
}

func (w *Writer) WriteMessage(message []byte) {
	err = w.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
	log.Println(message)
	_, b, _ := busCodec.NativeFromBinary(message)
	log.Println(b)
	_, err = w.conn.Write(b)
	log.Println("end")
}

func (w *Writer) CloseKafkaWriter() {
	err := w.conn.Close()
	if err != nil {
		log.Fatal("failed to close Writer:", err)
	}
}

func (r *Reader) GetKafkaReader() {

}
