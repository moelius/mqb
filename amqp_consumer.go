package mqb

import (
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"sync"
)

//extAmqpChannelInterface interface for "github.com/streadway/amqp" Channel
type extAmqpChannelInterface interface {
	io.Closer
	amqp.Acknowledger
	NotifyClose(chan *amqp.Error) chan *amqp.Error
	ExchangeDeclare(name string, kind string, durable bool, autoDelete bool, internal bool, noWait bool, args amqp.Table) error
	QueueDeclare(name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error)
	QueueBind(name string, key string, exchange string, noWait bool, args amqp.Table) error
	Consume(queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Qos(prefetchCount int, prefetchSize int, global bool) error
	Publish(exchange string, key string, mandatory bool, immediate bool, msg amqp.Publishing) error
}

//AmqpConsumerInterface interface for task
type AmqpConsumerInterface interface {
	Destroy() error
	Start() error
	Stop() error
}

type amqpConsumerInterface interface {
	ConsumerInterface
	getChannel() extAmqpChannelInterface
	setConnection(amqpConnectionInterface)
	getConnection() amqpConnectionInterface
	setChannel(extAmqpChannelInterface)
	setDelivery(<-chan amqp.Delivery)
	getDelivery() <-chan amqp.Delivery
	setNotify(chan error)
	getNotify() chan error
	sentinel()
	consume()
}

type amqpConsumer struct {
	sync.RWMutex
	name          string
	queueName     string
	broker        amqpBrokerInterface
	conn          amqpConnectionInterface
	channel       extAmqpChannelInterface
	deliveryChan  <-chan amqp.Delivery
	options       *AmqpConsumerOptions
	pool          poolInterface
	notifyChannel chan *amqp.Error
	notify        chan error
	state         int
	ctx           context.Context
}

func (consumer *amqpConsumer) setName(name string) {
	consumer.name = name
}

func (consumer *amqpConsumer) getName() string {
	return consumer.name
}

func (consumer *amqpConsumer) setChannel(channel extAmqpChannelInterface) {
	consumer.Lock()
	consumer.channel = channel
	consumer.Unlock()
}

func (consumer *amqpConsumer) getChannel() extAmqpChannelInterface {
	consumer.RLock()
	defer consumer.RUnlock()
	return consumer.channel
}

func (consumer *amqpConsumer) setDelivery(delivery <-chan amqp.Delivery) {
	consumer.Lock()
	consumer.deliveryChan = delivery
	consumer.Unlock()
}

func (consumer *amqpConsumer) getDelivery() <-chan amqp.Delivery {
	consumer.RLock()
	defer consumer.RUnlock()
	return consumer.deliveryChan
}

func (consumer *amqpConsumer) setNotify(notify chan error) {
	consumer.Lock()
	consumer.notify = notify
	consumer.Unlock()
}

func (consumer *amqpConsumer) getNotify() chan error {
	consumer.RLock()
	defer consumer.RUnlock()
	return consumer.notify
}

func (consumer *amqpConsumer) setConnection(conn amqpConnectionInterface) {
	consumer.Lock()
	consumer.conn = conn
	consumer.Unlock()
}

func (consumer *amqpConsumer) getConnection() amqpConnectionInterface {
	consumer.RLock()
	defer consumer.RUnlock()
	return consumer.conn
}

func (consumer *amqpConsumer) getPool() poolInterface {
	return consumer.pool
}

func (consumer *amqpConsumer) poolStart() error {
	return consumer.pool.start()
}

func (consumer *amqpConsumer) poolStop() {
	consumer.pool.stop()
}

func (consumer *amqpConsumer) Start() (err error) {
	if consumer.getConnection() == nil {
		return fmt.Errorf("can't start consume with <nil> connection")
	}

	channel, err := consumer.conn.channel()
	if err != nil {
		return
	}

	if consumer.options.ChannelOptions.Exchange != "" {
		if err = channel.ExchangeDeclare(
			consumer.options.ChannelOptions.Exchange,
			consumer.options.ChannelOptions.ExchangeType,
			consumer.options.ChannelOptions.Durable,
			consumer.options.ChannelOptions.AutoDelete,
			consumer.options.ChannelOptions.Internal,
			consumer.options.ChannelOptions.NoWait,
			consumer.options.ChannelOptions.Args,
		); err != nil {
			return fmt.Errorf("error while exchange declare [%s]", err)
		}
	}

	if consumer.options.BindingOptions != nil {
		if consumer.options.BindingOptions.Key == AmqpDirectReplyKey {
			consumer.queueName = AmqpDirectReplyKey
		} else {
			queue, e := channel.QueueDeclare(
				consumer.options.BindingOptions.Queue.Name,
				consumer.options.BindingOptions.Queue.Durable,
				consumer.options.BindingOptions.Queue.AutoDelete,
				consumer.options.BindingOptions.Queue.Exclusive,
				consumer.options.BindingOptions.Queue.NoWait,
				consumer.options.BindingOptions.Queue.Args,
			)
			if e != nil {
				return fmt.Errorf("error while queue declare [%s]", err)
			}
			consumer.queueName = queue.Name
		}
	}

	if consumer.options.QosOptions != nil {
		err = channel.Qos(consumer.options.QosOptions.PrefetchCount, consumer.options.QosOptions.PrefetchSize, consumer.options.QosOptions.Global)
		if err != nil {
			return fmt.Errorf("error while setting `qos` [%s]", err)
		}
	}

	if consumer.options.BindingOptions != nil && consumer.options.BindingOptions.Key != AmqpDirectReplyKey {
		if err = channel.QueueBind(
			consumer.queueName,
			consumer.options.BindingOptions.Key,
			consumer.options.ChannelOptions.Exchange,
			consumer.options.BindingOptions.NoWait,
			consumer.options.BindingOptions.Args,
		); err != nil {
			return fmt.Errorf("error while queue bind [%s]", err)
		}
	}
	if consumer.options.ConsumerTag == "" {
		consumer.options.ConsumerTag = consumer.name
	}

	deliveryChan, err := channel.Consume(
		consumer.queueName,
		consumer.options.ConsumerTag,
		consumer.options.AutoAck,
		consumer.options.Exclusive,
		consumer.options.NoLocal,
		consumer.options.NoWait,
		consumer.options.Args,
	)

	if err != nil {
		return fmt.Errorf("can't consume from queue [%s], error=[%s]", consumer.queueName, err)
	}

	consumer.setDelivery(deliveryChan)
	consumer.setChannel(channel)
	consumer.notifyChannel = make(chan *amqp.Error)
	consumer.setNotify(make(chan error))
	consumer.consume()
	consumer.sentinel()
	consumer.setState(opened)
	Logger.Debug(fmt.Sprintf("consumer started: name=[%s]", consumer.getName()))
	return
}

func (consumer *amqpConsumer) Stop() (err error) {
	if consumer.getState() == opened {
		err = consumer.getChannel().Close()
	}
	return
}

func (consumer *amqpConsumer) Destroy() (err error) {
	if err = consumer.Stop(); err != nil {
		return
	}
	if consumer.getPool() != nil {
		consumer.poolStop()
	}
	return consumer.broker.deleteConsumer(consumer.name)
}

func (consumer *amqpConsumer) sentinel() {
	go func() {
		e := <-consumer.getChannel().NotifyClose(consumer.notifyChannel)
		consumer.setState(closed)
		if e != nil {
			Logger.Warning(fmt.Sprintf("consumer name=[%s] stopped with error=[%s]", consumer.name, e))
		}
		Logger.Debug(fmt.Sprintf("consumer stopped: name=[%s]", consumer.getName()))
		consumer.getNotify() <- nil
		<-consumer.getNotify()
	}()
}

func (consumer *amqpConsumer) consume() {
	go func() {
	Loop:
		for {
			select {
			case delivery, ok := <-consumer.getDelivery():
				if !ok {
					continue Loop
				}
				task := &AmqpTask{
					Broker:   consumer.broker,
					Consumer: consumer,
					Delivery: delivery,
				}
				//handle situation when two or more consumers process one queue with different routing keys
				if delivery.ConsumerTag == consumer.options.ConsumerTag {
					consumer.getPool().getRequestChan() <- task
				} else {
					cName := amqpConsumerName(delivery.Exchange, delivery.RoutingKey)
					anotherConsumer, ok := consumer.broker.getConsumer(cName)
					if !ok {
						Logger.Warning(fmt.Sprintf("have no consumer in registry with name=[%s]", cName))
						continue Loop
					}
					//send to pool
					anotherConsumer.getPool().getRequestChan() <- task
				}
			case <-consumer.getNotify():
				break Loop
			case <-consumer.ctx.Done():
				if err := consumer.Destroy(); err != nil {
					Logger.Fatal(err.Error())
				}
				return
			}
		}
		consumer.getNotify() <- nil
	}()
}

func (consumer *amqpConsumer) getState() int {
	consumer.RLock()
	defer consumer.RUnlock()
	return consumer.state
}

func (consumer *amqpConsumer) setState(state int) {
	consumer.Lock()
	consumer.state = state
	consumer.Unlock()
}

func amqpConsumerName(exchange string, routingKey string) string {
	if exchange == "" {
		exchange = "default"
	}
	if routingKey == AmqpDirectReplyKey {
		routingKey = "reply"
	}
	return fmt.Sprintf("%s.%s", exchange, routingKey)
}
