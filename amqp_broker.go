package mqb

import (
	"context"
	"fmt"
	"sync"
)

const (
	//states
	closed = iota
	opened
	closing
)

//AmqpBrokerInterface broker interface for task structure
type AmqpBrokerInterface interface {
	Publish(iRequest interface{}) error
}

//amqpBrokerInterface amqp broker interface
type amqpBrokerInterface interface {
	BrokerInterface
	deleteConsumer(name string) error
	getInputConnection() amqpConnectionInterface
	getOutputConnection() amqpConnectionInterface
	getConsumer(string) (amqpConsumerInterface, bool)
	getProducer(string) (amqpProducerInterface, bool)
	startSentinel(amqpConnectionInterface) error
	startConsumers() error
	stopConsumers() error
	startProducers() error
	stopProducers() error
}

//amqpBroker structure
type amqpBroker struct {
	options   *brokerOptions          //transport brokerOptions
	inConn    amqpConnectionInterface //incoming connection
	outConn   amqpConnectionInterface //outgoing connection
	consumers sync.Map                //goroutine-safe store for consumers
	producers sync.Map                //goroutine-safe store for producers
	stopChan  chan bool               //transport stop channel
	notify    chan error              //transport notify channel
}

//newAmqpBroker make new broker
func newAmqpBroker() (broker *amqpBroker) {
	broker = &amqpBroker{stopChan: make(chan bool)}
	return
}

//newConsumer make new amqpConsumer
func (broker *amqpBroker) newConsumer(ctx context.Context, iConsumerOptions interface{}, pool poolInterface) (consumer ConsumerInterface, err error) {
	consumerOpts, ok := iConsumerOptions.(*AmqpConsumerOptions)
	if !ok {
		return nil, fmt.Errorf("wrong type of `consumer options`, you must use `*mqb.AmqpConsumerOptions` struct")
	}
	consumer = &amqpConsumer{broker: broker, pool: pool, options: consumerOpts, ctx: ctx}
	name := amqpConsumerName(consumerOpts.ChannelOptions.Exchange, consumerOpts.BindingOptions.Key)
	consumer.setName(name)
	broker.consumers.Store(name, consumer)
	return
}

func (broker *amqpBroker) configure(options brokerOptions) error {
	broker.options = &options
	broker.inConn = newAmqpConnection(broker.options.url)
	broker.outConn = newAmqpConnection(broker.options.url)
	broker.inConn.setMaxReconnects(broker.options.maxReconnects)
	broker.inConn.setReconnectTimeout(broker.options.reconnectTimeout)
	broker.outConn.setMaxReconnects(broker.options.maxReconnects)
	broker.outConn.setReconnectTimeout(broker.options.reconnectTimeout)
	return nil
}

//getNotify get broker notify channel
func (broker *amqpBroker) getNotify() chan error {
	return broker.notify
}

//getInputConnection get broker input amqpConnection
func (broker *amqpBroker) getInputConnection() amqpConnectionInterface {
	return broker.inConn
}

func (broker *amqpBroker) getOutputConnection() amqpConnectionInterface {
	return broker.outConn
}

func (broker *amqpBroker) newProducer(iProducerOptions interface{}) (err error) {
	producerOpts, ok := iProducerOptions.(*AmqpProducerOptions)
	if !ok {
		return fmt.Errorf("wrong type of `producer options`, you must use `*mqb.AmqpProducerOptions` struct")
	}
	name := producerName(producerOpts.Exchange)
	broker.producers.Store(name, &amqpProducer{
		name:        name,
		publishChan: make(chan *AmqpRequest),
		options:     producerOpts,
	})
	return
}

func (broker *amqpBroker) serve() (err error) {
	broker.notify = make(chan error)

	sentinelProcess := func(conn amqpConnectionInterface) {
		for {
			if conn.getState() == closing {
				break
			}
			if e := broker.startSentinel(conn); e != nil {
				broker.getNotify() <- e
			}
		}
	}

	conn := broker.getInputConnection()
	if err = conn.open(); err != nil {
		return
	}
	Logger.Debug("input connection opened")

	go sentinelProcess(conn)

	conn = broker.getOutputConnection()
	if err = conn.open(); err != nil {
		return
	}
	Logger.Debug("output connection opened")

	go sentinelProcess(conn)

	if err = broker.startProducers(); err != nil {
		return
	}

	if err = broker.startConsumers(); err != nil {
		return
	}
	return
}

func (broker *amqpBroker) stop() (err error) {
	if err = broker.stopConsumers(); err != nil {
		return
	}
	if err = broker.inConn.close(); err != nil {
		return
	}
	Logger.Debug("input connection closed")
	if err = broker.stopProducers(); err != nil {
		return
	}
	if err = broker.outConn.close(); err != nil {
		return
	}
	Logger.Debug("output connection closed")
	return
}

func (broker *amqpBroker) Publish(iRequest interface{}) (err error) {
	request, ok := iRequest.(*AmqpRequest)
	if !ok {
		return fmt.Errorf("wrong type of `request`, you must use `*mqb.AmqpRequest` struct")
	}
	iProducer, ok := broker.producers.Load(request.Exchange)
	//If we have producer
	if ok {
		producer := iProducer.(amqpProducerInterface)
		producer.getPublishChan() <- request
		return
	}
	//If we have no producer, Publish from new channel
	ch, err := broker.outConn.channel()
	if err != nil {
		return
	}
	defer func() {
		err = ch.Close()
	}()
	return ch.Publish(request.Exchange, request.Key, request.Mandatory, request.Immediate, request.convertToAmqpPublishing())
}

func (broker *amqpBroker) startSentinel(conn amqpConnectionInterface) (err error) {
	var t, ioType string

	switch conn {
	case broker.getInputConnection():
		t = "consumers"
		ioType = "input"
	case broker.getOutputConnection():
		t = "producers"
		ioType = "output"
	}

	amqpError := <-conn.getNotifyCloser()
	if amqpError == nil {
		return
	}

	conn.setState(closed)

	Logger.Warning(fmt.Sprintf("start reconnection [%s connection] due to an error=[%s]", ioType, amqpError.Error()))
	if err = conn.reconnect(); err != nil {
		Logger.Error(fmt.Sprintf("reconnection process for [%s connection] failed with error=[%s]", ioType, err))
		return
	}

	switch t {
	case "consumers":
		broker.consumers.Range(func(name interface{}, cInterface interface{}) (ok bool) {
			consumer := cInterface.(amqpConsumerInterface)
			err = consumer.Stop()
			if err != nil {
				return false
			}
			return consumer.Start() == nil
		})
	case "producers":
		broker.producers.Range(func(name interface{}, pInterface interface{}) (ok bool) {
			producer := pInterface.(amqpProducerInterface)
			if err = producer.stop(); err != nil {
				return false
			}
			return producer.start() == nil
		})
	}
	if err != nil {
		return
	}

	Logger.Debug(fmt.Sprintf(fmt.Sprintf("%s restarted", t)))
	Logger.Info(fmt.Sprintf(fmt.Sprintf("reconnection process for [%s connection] completed successfully", ioType)))
	return
}

func (broker *amqpBroker) startConsumer(name string) (err error) {
	iConsumer, ok := broker.consumers.Load(name)
	if !ok {
		return fmt.Errorf("can't start consumer: have no consumer with name `%s`", name)
	}
	consumer := iConsumer.(amqpConsumerInterface)
	consumer.setConnection(broker.inConn)
	return consumer.Start()
}

func (broker *amqpBroker) startConsumers() (err error) {
	broker.consumers.Range(func(name interface{}, cInterface interface{}) (ok bool) {
		consumer := cInterface.(amqpConsumerInterface)
		consumer.setConnection(broker.inConn)
		if consumer.getPool() != nil {
			err = consumer.poolStart()
			if err != nil {
				return false
			}
		}
		return consumer.Start() == nil
	})
	return
}

func (broker *amqpBroker) startProducers() (err error) {
	broker.producers.Range(func(name interface{}, pInterface interface{}) (ok bool) {
		producer := pInterface.(amqpProducerInterface)
		producer.setConnection(broker.outConn)
		return producer.start() == nil
	})
	return
}

func (broker *amqpBroker) stopConsumers() (err error) {
	broker.consumers.Range(func(name interface{}, iConsumer interface{}) (ok bool) {
		consumer := iConsumer.(amqpConsumerInterface)
		return consumer.Destroy() == nil
	})
	return
}

func (broker *amqpBroker) stopProducers() (err error) {
	broker.producers.Range(func(name interface{}, iProducer interface{}) (ok bool) {
		producer := iProducer.(amqpProducerInterface)
		return producer.stop() == nil
	})
	return
}

func (broker *amqpBroker) getConsumer(name string) (iConsumer amqpConsumerInterface, ok bool) {
	consumer, loadOk := broker.consumers.Load(name)
	if !loadOk {
		return
	}
	return consumer.(amqpConsumerInterface), true
}

func (broker *amqpBroker) getProducer(name string) (iProducer amqpProducerInterface, ok bool) {
	producer, loadOk := broker.producers.Load(name)
	if !loadOk {
		return
	}
	return producer.(amqpProducerInterface), true
}

func (broker *amqpBroker) deleteConsumer(name string) error {
	if _, ok := broker.consumers.Load(name); !ok {
		return fmt.Errorf("can't delete consumer from store: have no consumer with name `%s`", name)
	}
	broker.consumers.Delete(name)
	return nil
}
