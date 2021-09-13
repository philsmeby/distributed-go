package main

import (
	"bytes"
	"encoding/gob"
	"github.com/distributed-go/datamanager"
	"github.com/distributed-go/dto"
	"github.com/distributed-go/qutils"
	"log"
)

const url = "amqp://guest@localhost:5672"

func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()

	msgs, err := ch.Consume(
		qutils.PersistentReadingsQueue,
		"",
		false,
		true,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalln("Failed to get access to messages")
	}

	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)

		sd := &dto.SensorMessage{}
		_ = dec.Decode(sd)

		err := datamanager.SaveReading(sd)

		if err != nil {
			log.Printf("Failed to save reading from sensor %v. Error: %s", sd.Name, err.Error())
		} else {
			_ = msg.Ack(false)
		}
	}
}
