package coordinator

import (
	"bytes"
	"encoding/gob"
	"github.com/distributed-go/dto"
	"github.com/distributed-go/qutils"
	"github.com/streadway/amqp"
	"time"
)

const maxRate = 5 * time.Second

type DatabaseConsumer struct {
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	queue   *amqp.Queue
	sources []string
}

func NewDatabaseConsumer(er EventRaiser) *DatabaseConsumer {
	dc := DatabaseConsumer{
		er: er,
	}
	dc.conn, dc.ch = qutils.GetChannel(url)
	dc.queue = qutils.GetQueue(qutils.PersistentReadingsQueue, dc.ch, false)

	dc.er.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		dc.SubscribeToDataEvent(eventData.(string))
	})
	return &dc
}

func (dc *DatabaseConsumer) SubscribeToDataEvent(eventName string) {
	for _, v := range dc.sources {
		if v == eventName {
			return
		}
	}
	callback := dc.callbackGenerator()

	dc.er.AddListener("MessageReceived_"+eventName, callback)
}

func (dc *DatabaseConsumer) callbackGenerator() func(interface{}) {
	prevTime := time.Unix(0, 0)
	buf := new(bytes.Buffer)

	return func(eventData interface{}) {
		ed := eventData.(EventData)
		if time.Since(prevTime) > maxRate {
			prevTime = time.Now()

			sm := dto.SensorMessage{
				Name:      ed.Name,
				Value:     ed.Value,
				Timestamp: ed.Timestamp,
			}
			buf.Reset()
			enc := gob.NewEncoder(buf)
			_ = enc.Encode(sm)

			msg := amqp.Publishing{
				Body: buf.Bytes(),
			}

			_ = dc.ch.Publish("", qutils.PersistentReadingsQueue, false, false, msg)
		}
	}
}
