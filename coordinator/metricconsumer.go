package coordinator

import (
	"bytes"
	"encoding/gob"
	"github.com/distributed-go/dto"
	"github.com/distributed-go/monitoring"
	"github.com/distributed-go/qutils"
	"github.com/streadway/amqp"
	"time"
)

type MetricConsumer struct {
	rc      monitoring.ReadingCounterInterface
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	queue   *amqp.Queue
	sources []string
}

func NewMetricConsumer(er EventRaiser) *MetricConsumer {
	mc := MetricConsumer{
		er: er,
		rc: monitoring.NewReadingCounter(),
	}
	mc.conn, mc.ch = qutils.GetChannel(url)
	mc.queue = qutils.GetQueue(qutils.LiveReadingsQueue, mc.ch, false)

	mc.er.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		mc.SubscribeToDataEvent(eventData.(string))
	})
	return &mc
}

func (mc *MetricConsumer) SubscribeToDataEvent(eventName string) {
	for _, v := range mc.sources {
		if v == eventName {
			return
		}
	}
	callback := mc.callbackGenerator()

	mc.er.AddListener("MessageReceived_"+eventName, callback)
}

func (mc *MetricConsumer) callbackGenerator() func(interface{}) {
	prevTime := time.Unix(0, 0)
	buf := new(bytes.Buffer)

	return func(eventData interface{}) {
		ed := eventData.(EventData)

		if time.Since(prevTime) > 700*time.Millisecond {
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

			err := mc.ch.Publish("", qutils.LiveReadingsQueue, false, false, msg)

			if err == nil {
				mc.rc.Increment("coordinator", ed.Name)
			}
		}
	}
}
