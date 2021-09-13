package main

import (
	"bytes"
	"encoding/gob"
	"github.com/distributed-go/dto"
	"github.com/distributed-go/monitoring"
	"github.com/distributed-go/qutils"
	"log"
	"math/rand"
	"time"
)

const url = qutils.BrokerUrl

func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()

	msgs, err := ch.Consume(
		qutils.LiveReadingsQueue,
		"",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		log.Fatalln("Failed to get access to messages")
	}

	go monitoring.MetricExporter()

	rc := monitoring.NewReadingCounter()

	rg := monitoring.NewReadingGauge()

	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)

		sd := &dto.SensorMessage{}
		_ = dec.Decode(sd)

		rg.Set(sd.Value, sd.Name)
		rc.Increment("monitoring", sd.Name)

		if err != nil {
			log.Printf("Failed to save reading from sensor %v. Error: %s", sd.Name, err.Error())
		} else {
			_ = msg.Ack(false)
		}

		time.Sleep(time.Duration(rand.Intn(400)) * time.Millisecond) // Delay to simulate network latency
	}
}
