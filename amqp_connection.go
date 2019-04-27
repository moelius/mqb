package mqb

import (
	"fmt"
	"github.com/streadway/amqp"
	"io"
	"sync"
	"time"
)

type extAmqpConnectionInterface interface {
	io.Closer
	NotifyClose(chan *amqp.Error) chan *amqp.Error
	Channel() (*amqp.Channel, error)
}

type amqpDialerInterface interface {
	dial(string) (extAmqpConnectionInterface, error)
}

//amqpConnectionInterface amqpConnection interface
type amqpConnectionInterface interface {
	open() error
	close() error
	channel() (extAmqpChannelInterface, error)
	reconnect() error
	getState() int
	setState(int)
	setMaxReconnects(int)
	setReconnectTimeout(time.Duration)
	getNotifyCloser() chan *amqp.Error
	setNotifyCloser(chan *amqp.Error)
}

type amqpDialer struct{}

func (d *amqpDialer) dial(url string) (extAmqpConnectionInterface, error) {
	return amqp.Dial(url)
}

func newAmqpConnection(url string) amqpConnectionInterface {
	return &amqpConnection{url: url, dialer: &amqpDialer{}}
}

type amqpConnection struct {
	sync.RWMutex
	url              string
	dialer           amqpDialerInterface
	amqp             extAmqpConnectionInterface
	state            int
	notify           chan *amqp.Error
	maxReconnects    int
	reconnectTimeout time.Duration
}

func (conn *amqpConnection) setMaxReconnects(max int) {
	conn.maxReconnects = max
}

func (conn *amqpConnection) setReconnectTimeout(timeout time.Duration) {
	conn.reconnectTimeout = timeout
}

func (conn *amqpConnection) getNotifyCloser() chan *amqp.Error {
	conn.RLock()
	defer conn.RUnlock()
	return conn.notify
}

func (conn *amqpConnection) setNotifyCloser(notify chan *amqp.Error) {
	conn.Lock()
	conn.notify = notify
	conn.Unlock()
}

func (conn *amqpConnection) open() (err error) {
	conn.Lock()
	if conn.amqp, err = conn.dialer.dial(conn.url); err != nil {
		conn.Unlock()
		return
	}
	conn.Unlock()
	conn.setNotifyCloser(conn.amqp.NotifyClose(make(chan *amqp.Error)))
	conn.setState(opened)
	return
}

func (conn *amqpConnection) close() (err error) {
	if conn.getState() == closed || conn.getState() == closing {
		return
	}
	conn.setState(closing)
	conn.Lock()
	defer conn.Unlock()
	if err = conn.amqp.Close(); err != nil {
		return
	}
	return
}

func (conn *amqpConnection) channel() (extAmqpChannelInterface, error) {
	if conn.amqp == nil || conn.getState() != opened {
		return nil, fmt.Errorf("connection not opened, can't create new channel")
	}
	return conn.amqp.Channel()
}

func (conn *amqpConnection) reconnect() (err error) {
	if conn.getState() == opened {
		return
	}
	ticker := time.NewTicker(conn.reconnectTimeout)
	reconCounter := 0
	for range ticker.C {
		reconCounter++
		err = conn.open()
		if err != nil {
			Logger.Warning(fmt.Sprintf("reconnection error=[%s]", err))
			if reconCounter == conn.maxReconnects {
				break
			}
			continue
		}
		break
	}
	ticker.Stop()
	return
}

func (conn *amqpConnection) getState() int {
	conn.RLock()
	defer conn.RUnlock()
	return conn.state
}

func (conn *amqpConnection) setState(state int) {
	conn.Lock()
	conn.state = state
	conn.Unlock()
}
