package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/distributed-go/dto"
	"github.com/distributed-go/monitoring"
	"github.com/distributed-go/qutils"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var url = qutils.BrokerUrl
var name = flag.String("name", "sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "frequency in cycle/sec")
var max = flag.Float64("max", 5., "maximum value for generated reading")
var min = flag.Float64("min", 1., "minimum value for generated reading")
var stepSize = flag.Float64("step", 0.1, "maximum allowable cange per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var value float64
var nom float64

func main() {
	go monitoring.MetricExporter()
	flag.Parse()
	value = r.Float64()*(*max-*min) + *min
	nom = (*max-*min)/2 + *min

	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch, false)

	//////////////////////////-----------------\\\\\\\\\\\\\\\\\\\\\\\\\\\\\
	// Create central queue to hold names of each queue created by a sensor
	//sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)
	publishQueueName(ch)
	discoveryQueue := qutils.GetQueue("", ch, true)

	_ = ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil)

	_ = ch.QueueBind(
		discoveryQueue.Name,
		"",
		qutils.SensorDiscoveryExchange,
		false,
		nil)

	go listenForDiscoverRequests(discoveryQueue.Name, ch)
	//////////////////////////-----------------\\\\\\\\\\\\\\\\\\\\\\\\\\\\\

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	signal := time.Tick(dur)

	buf := new(bytes.Buffer)

	for range signal {
		calcValue()

		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}

		enc := gob.NewEncoder(buf)

		_ = enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		delay := rand.Intn(400)
		time.Sleep(time.Duration(delay) * time.Millisecond)

		_ = ch.Publish("", dataQueue.Name, false, false, msg)

		log.Printf("%s Reading Sent. Value: %v\n", *name, value)

		buf.Reset()

	}
}

func calcValue() {
	var maxStep, minStep float64
	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}

func publishQueueName(ch *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(*name)}
	_ = ch.Publish("amq.fanout", "", false, false, msg)
}

func listenForDiscoverRequests(name string, ch *amqp.Channel) {
	msgs, _ := ch.Consume(name, "", true, false, false, false, nil)

	for range msgs {
		fmt.Print("Received discovery request")
		publishQueueName(ch)
	}
}
