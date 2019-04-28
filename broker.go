package mqb

//go:generate mockgen -source=broker.go -destination=mock_broker_test.go -package=mqb BrokerInterface ConsumerInterface

import (
	"context"
	"time"
)

//BrokerInterface interface
type BrokerInterface interface {
	configure(brokerOptions) error
	getNotify() chan error
	serve() error
	stop() error
	newConsumer(ctx context.Context, iConsumerOptions interface{}, pool poolInterface) (ConsumerInterface, error)
	newProducer(iProducerOptions interface{}) error
	Publish(iRequest interface{}) error
	startConsumer(name string) (err error)
}

//ConsumerInterface interface
type ConsumerInterface interface {
	Start() error
	Stop() error
	Destroy() error
	getState() int
	setState(int)
	setName(name string)
	getName() string
	getPool() poolInterface
	poolStart() error
	poolStop()
}

type brokerOptions struct {
	url              string
	maxReconnects    int
	reconnectTimeout time.Duration
}

var defaultBrokerOptions = brokerOptions{
	url:              "amqp://",
	maxReconnects:    100,
	reconnectTimeout: 10 * time.Second,
}

//BrokerOption broker option
type BrokerOption func(options *brokerOptions)

//URL broker url option
func URL(url string) BrokerOption {
	return func(o *brokerOptions) {
		o.url = url
	}
}

//MaxReconnects broker max reconnects option
func MaxReconnects(max int) BrokerOption {
	return func(o *brokerOptions) {
		o.maxReconnects = max
	}
}

//ReconnectTimeout broker reconnect timeout
func ReconnectTimeout(timeout time.Duration) BrokerOption {
	return func(o *brokerOptions) {
		o.reconnectTimeout = timeout
	}
}

func configure(broker BrokerInterface, options brokerOptions) error {
	return broker.configure(options)
}

func serve(broker BrokerInterface) error {
	return broker.serve()
}

func stop(broker BrokerInterface) error {
	return broker.stop()
}

func newConsumer(ctx context.Context, broker BrokerInterface, iConsumerOptions interface{}, pool poolInterface) (consumer ConsumerInterface, err error) {
	return broker.newConsumer(ctx, iConsumerOptions, pool)
}

func newProducer(broker BrokerInterface, iProducerOptions interface{}) error {
	return broker.newProducer(iProducerOptions)
}

func publish(broker BrokerInterface, iRequest interface{}) error {
	return broker.Publish(iRequest)
}

func startConsumer(broker BrokerInterface, name string) (err error) {
	return broker.startConsumer(name)
}
