package mqb

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

type amqpProducerInterface interface {
	start() error
	stop() error
	getState() int
	setState(int)
	getPublishChan() chan *AmqpRequest
	setConnection(amqpConnectionInterface)
	setNotify(chan error)
	getNotify() chan error
	publish(exchange string, key string, mandatory bool, immediate bool, msg AmqpPublishing) error
	sentinel()
	produce()
}

type amqpProducer struct {
	sync.RWMutex
	name          string
	conn          amqpConnectionInterface
	channel       extAmqpChannelInterface
	publishChan   chan *AmqpRequest
	options       *AmqpProducerOptions
	notifyChannel chan *amqp.Error
	notify        chan error
	state         int
}

func (producer *amqpProducer) setConnection(conn amqpConnectionInterface) {
	producer.conn = conn
}

func (producer *amqpProducer) getPublishChan() chan *AmqpRequest {
	return producer.publishChan
}

func (producer *amqpProducer) setNotify(notify chan error) {
	producer.Lock()
	producer.notify = notify
	producer.Unlock()
}

func (producer *amqpProducer) getNotify() chan error {
	producer.RLock()
	defer producer.RUnlock()
	return producer.notify
}

func (producer *amqpProducer) start() (err error) {
	if producer.conn == nil {
		return fmt.Errorf("can't Start produce with <nil> connection")
	}
	channel, err := producer.conn.channel()
	if err != nil {
		return
	}

	if producer.options.ExchangeType == "" {
		producer.options.ExchangeType = AmqpExchangeDirect
	}
	if err = channel.ExchangeDeclare(
		producer.options.Exchange,
		producer.options.ExchangeType,
		producer.options.Durable,
		producer.options.AutoDelete,
		producer.options.Internal,
		producer.options.NoWait,
		producer.options.Args,
	); err != nil {
		return fmt.Errorf("error while exchange declare [%s]", err)
	}
	producer.channel = channel
	producer.notifyChannel = make(chan *amqp.Error)
	producer.setNotify(make(chan error))
	producer.produce()
	producer.sentinel()
	producer.setState(opened)
	Logger.Debug(fmt.Sprintf("producer started: name=[%s]", producer.name))
	return
}

func (producer *amqpProducer) stop() (err error) {
	if producer.getState() == opened {
		err = producer.channel.Close()
	}
	return
}

func (producer *amqpProducer) sentinel() {
	go func() {
		<-producer.channel.NotifyClose(producer.notifyChannel)
		producer.setState(closed)
		Logger.Debug(fmt.Sprintf("producer stopped: name=[%s]", producer.name))
		producer.getNotify() <- nil
		<-producer.getNotify()
	}()
}

func (producer *amqpProducer) publish(exchange string, key string, mandatory bool, immediate bool, msg AmqpPublishing) error {
	return producer.channel.Publish(exchange, key, mandatory, immediate, msg)
}

func (producer *amqpProducer) produce() {
	go func() {
	Loop:
		for {
			select {
			case req := <-producer.publishChan:
				if e := producer.publish(req.Exchange, req.Key, req.Mandatory, req.Immediate, req.convertToAmqpPublishing()); e != nil {
					log.Println(e)
				}
			case <-producer.getNotify():
				break Loop
			}
		}
		producer.getNotify() <- nil
	}()
}

func (producer *amqpProducer) getState() int {
	producer.RLock()
	defer producer.RUnlock()
	return producer.state
}

func (producer *amqpProducer) setState(state int) {
	producer.Lock()
	producer.state = state
	producer.Unlock()
}

func producerName(exchange string) string {
	if exchange == "" {
		exchange = "default"
	}
	return exchange
}
