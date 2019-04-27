package mqb

import (
	"github.com/streadway/amqp"
	"time"
)

//AmqpPublishing alias for "github.com/streadway/amqp" `amqp.Publishing` type
type AmqpPublishing = amqp.Publishing

//AmqpTable alias for "github.com/streadway/amqp" `amqp.Table` type
type AmqpTable = amqp.Table

//AmqpHeaders alias for "github.com/streadway/amqp" `amqp.Table` type
type AmqpHeaders = amqp.Table

//AmqpDelivery alias for "github.com/streadway/amqp" `amqp.Delivery` type
type AmqpDelivery = amqp.Delivery

const (
	//AmqpDirectReplyKey key for direct reply pattern https://www.rabbitmq.com/direct-reply-to.html
	AmqpDirectReplyKey = "amq.rabbitmq.reply-to"

	//AmqpExchangeDirect type of rabbitmq exchange - `direct`
	AmqpExchangeDirect = amqp.ExchangeDirect

	//AmqpExchangeFanout type of rabbitmq exchange - `funout`
	AmqpExchangeFanout = amqp.ExchangeFanout

	//AmqpExchangeTopic type of rabbitmq exchange - `topic`
	AmqpExchangeTopic = amqp.ExchangeTopic
)

//AmqpChannelOptions amqp channel brokerOptions
type AmqpChannelOptions struct {
	Exchange     string
	ExchangeType string
	Durable      bool
	AutoDelete   bool
	Internal     bool
	NoWait       bool
	Args         AmqpTable
}

//AmqpQosOptions amqp qos brokerOptions
type AmqpQosOptions struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

//AmqpQueueOptions amqp queue brokerOptions
type AmqpQueueOptions struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       AmqpTable
}

//AmqpBindingOptions amqp binding brokerOptions
type AmqpBindingOptions struct {
	Key    string
	Queue  *AmqpQueueOptions
	NoWait bool
	Args   AmqpTable
}

//AmqpConsumerOptions amqp amqpConsumer brokerOptions
type AmqpConsumerOptions struct {
	ConsumerTag    string
	AutoAck        bool
	Exclusive      bool
	NoLocal        bool
	NoWait         bool
	Args           AmqpTable
	MsgLimit       int //when limit will be reached amqpConsumer will be destroyed
	ChannelOptions *AmqpChannelOptions
	QosOptions     *AmqpQosOptions
	BindingOptions *AmqpBindingOptions
}

//AmqpProducerOptions alias to ChannelOptions
type AmqpProducerOptions = AmqpChannelOptions

//AmqpRequest amqp request for Publish
type AmqpRequest struct {
	Headers         AmqpTable
	Timestamp       time.Time
	Exchange        string
	Key             string
	ContentType     string
	ContentEncoding string
	CorrelationID   string
	ReplyTo         string
	Expiration      string
	MessageID       string
	Type            string
	UserID          string
	AppID           string
	Body            []byte
	DeliveryMode    uint8
	Priority        uint8
	Mandatory       bool
	Immediate       bool
}

func (request *AmqpRequest) convertToAmqpPublishing() AmqpPublishing {
	return AmqpPublishing{
		Headers:         request.Headers,
		ContentType:     request.ContentType,
		ContentEncoding: request.ContentEncoding,
		DeliveryMode:    request.DeliveryMode,
		Priority:        request.Priority,
		CorrelationId:   request.CorrelationID,
		ReplyTo:         request.ReplyTo,
		Expiration:      request.Expiration,
		MessageId:       request.MessageID,
		Timestamp:       request.Timestamp,
		Type:            request.Type,
		UserId:          request.UserID,
		AppId:           request.AppID,
		Body:            request.Body,
	}
}

//AmqpTask task structure for amqp callbacks
type AmqpTask struct {
	Broker   AmqpBrokerInterface
	Consumer AmqpConsumerInterface
	Delivery AmqpDelivery
}

//Ack task acknowledge
func (task *AmqpTask) Ack(multiply bool) error {
	return task.Delivery.Ack(multiply)
}

//Nack negative task acknowledge
func (task *AmqpTask) Nack(multiply bool, requeue bool) error {
	return task.Delivery.Nack(multiply, requeue)
}

//Reject reject task
func (task *AmqpTask) Reject(requeue bool) error {
	return task.Delivery.Reject(requeue)
}
