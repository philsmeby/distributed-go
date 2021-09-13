package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/distributed-go/dto"
	"github.com/distributed-go/qutils"
	"github.com/streadway/amqp"
	"math/rand"
	"time"
)

const url = qutils.BrokerUrl

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener(ea *EventAggregator) *QueueListener {
	ql := QueueListener{}

	ql.conn, ql.ch = qutils.GetChannel(url)
	ql.sources = map[string]<-chan amqp.Delivery{}
	ql.ea = ea
	return &ql
}

func (ql *QueueListener) DiscoverSensor() {
	_ = ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil)

	_ = ql.ch.Publish(qutils.SensorDiscoveryExchange, "", false, false, amqp.Publishing{})
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch, true)

	_ = ql.ch.QueueBind(
		q.Name,
		"",
		"amq.fanout",
		false,
		nil)

	msgs, _ := ql.ch.Consume(q.Name, "", true, false, false, false, nil)

	ql.DiscoverSensor()

	for msg := range msgs {
		fmt.Print("New sensor discovered")
		ql.ea.PublishEvent("DataSourceDiscovered", string(msg.Body))
		sourceChan, _ := ql.ch.Consume(
			string(msg.Body),
			"",
			true,
			false,
			false,
			false,
			nil)
		if ql.sources[string(msg.Body)] == nil {
			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)

		sd := new(dto.SensorMessage)

		_ = d.Decode(sd)

		ed := EventData{
			Name:      sd.Name,
			Value:     sd.Value,
			Timestamp: sd.Timestamp,
		}

		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)

		fmt.Printf("Received message: %v\n", sd)
		delay := rand.Intn(400)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}
