package mqb

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"syscall"
	"testing"
)

type ServerTestSuite struct {
	suite.Suite
	OKUrl        string
	MCtrl        *gomock.Controller
	Broker       *MockBrokerInterface
	Server       *Server
	Pool         *MockpoolInterface
	Consumer     *MockConsumerInterface
	TestErr      error
	ServerCancel context.CancelFunc
}

func (suite *ServerTestSuite) SetupTest() {
	suite.MCtrl = gomock.NewController(suite.T())
	defer suite.MCtrl.Finish()
	suite.Broker = NewMockBrokerInterface(suite.MCtrl)
	suite.OKUrl = "amqp://"
	suite.Consumer = NewMockConsumerInterface(suite.MCtrl)
	suite.Pool = NewMockpoolInterface(suite.MCtrl)
	suite.TestErr = fmt.Errorf("test error")
	ctx, cancel := context.WithCancel(context.Background())
	suite.ServerCancel = cancel
	suite.Server = &Server{
		broker: suite.Broker,
		ctx:    ctx,
		osSig:  make(chan os.Signal, 1),
	}
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (suite *ServerTestSuite) TestNewServerAmqp() {
	server, err := NewServer(context.Background(), URL(suite.OKUrl))
	ass := assert.New(suite.T())
	if ass.Nil(err) {
		ass.NotNil(server)
	}
}

func (suite *ServerTestSuite) TestNewServerWithUnknownBrokerURLErr() {
	url := "unknown://url"
	_, err := NewServer(context.Background(), URL(url))
	assert.NotNil(suite.T(), err)
}

func (suite *ServerTestSuite) TestServerWithLogger() {
	ass := assert.New(suite.T())
	logger := NewLoggerLog(DebugLevel)
	suite.Server.WithLogger(logger)
	ass.Equal(logger, Logger)
	ass.Equal(DebugLevel, Logger.GetLevel())

	logger = NewLoggerZap(WarnLevel)
	suite.Server.WithLogger(logger)
	ass.Equal(logger, Logger)
	ass.Equal(WarnLevel, Logger.GetLevel())
}

func (suite *ServerTestSuite) TestServerNewConsumerServerSide() {
	rPCallback := func(task *AmqpTask) error { return nil }
	suite.Broker.EXPECT().newConsumer(gomock.Any(), gomock.Any(), gomock.Any()).Return(suite.Consumer, nil)
	consumer, err := suite.Server.NewConsumer(context.Background(), AmqpConsumerOptions{}, rPCallback, CallbackMaxWorkers(10), CallbackArg(1))

	ass := assert.New(suite.T())
	ass.Nil(err)
	ass.Equal(suite.Consumer, consumer)
}

func (suite *ServerTestSuite) TestServerNewConsumerServerSideWithErr() {
	rPCallback := func(task *AmqpTask) error { return nil }
	suite.Broker.EXPECT().newConsumer(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, suite.TestErr)
	_, err := suite.Server.NewConsumer(context.Background(), AmqpConsumerOptions{}, rPCallback, CallbackMaxWorkers(10))
	assert.NotNil(suite.T(), err)
}

func (suite *ServerTestSuite) TestServerNewConsumerClientSide() {
	rPCallback := func(task *AmqpTask) error { return nil }
	suite.Server.started = true

	suite.Broker.EXPECT().newConsumer(gomock.Any(), gomock.Any(), gomock.Any()).Return(suite.Consumer, nil)
	suite.Consumer.EXPECT().getPool().Return(suite.Pool)
	suite.Consumer.EXPECT().getName()
	suite.Pool.EXPECT().start().Return(nil)
	suite.Broker.EXPECT().startConsumer(gomock.Any()).Return(nil)

	consumer, err := suite.Server.NewConsumer(context.Background(), AmqpConsumerOptions{}, rPCallback, CallbackMaxWorkers(10))
	ass := assert.New(suite.T())
	ass.Nil(err)
	ass.NotNil(consumer)
}

func (suite *ServerTestSuite) TestServerNewConsumerClientSideWithStartPoolErr() {
	rPCallback := func(task *AmqpTask) error { return nil }
	suite.Server.started = true

	suite.Broker.EXPECT().newConsumer(gomock.Any(), gomock.Any(), gomock.Any()).Return(suite.Consumer, nil)
	suite.Consumer.EXPECT().getPool().Return(suite.Pool)
	suite.Pool.EXPECT().start().Return(suite.TestErr)

	consumer, err := suite.Server.NewConsumer(context.Background(), AmqpConsumerOptions{}, rPCallback, CallbackMaxWorkers(10))
	ass := assert.New(suite.T())
	ass.NotNil(err)
	ass.Nil(consumer)
}

func (suite *ServerTestSuite) TestServerNewConsumerClientSideWithStartConsumerErr() {
	rPCallback := func(task *AmqpTask) error { return nil }
	suite.Server.started = true

	suite.Broker.EXPECT().newConsumer(gomock.Any(), gomock.Any(), gomock.Any()).Return(suite.Consumer, nil)
	suite.Consumer.EXPECT().getPool().Return(suite.Pool)
	suite.Consumer.EXPECT().getName()
	suite.Pool.EXPECT().start().Return(nil)
	suite.Broker.EXPECT().startConsumer(gomock.Any()).Return(suite.TestErr)

	consumer, err := suite.Server.NewConsumer(context.Background(), AmqpConsumerOptions{}, rPCallback, CallbackMaxWorkers(10))
	ass := assert.New(suite.T())
	ass.NotNil(err)
	ass.Nil(consumer)
}

func (suite *ServerTestSuite) TestServerNewProducerWithErr() {
	pOpts := &AmqpProducerOptions{}
	suite.Broker.EXPECT().newProducer(pOpts).Return(suite.TestErr)

	err := suite.Server.NewProducer(&AmqpProducerOptions{})
	ass := assert.New(suite.T())
	ass.NotNil(err)
}

func (suite *ServerTestSuite) TestListenAndServe() {
	suite.Broker.EXPECT().serve()
	suite.Broker.EXPECT().stop()
	suite.Broker.EXPECT().getNotify()
	go func() {
		suite.ServerCancel()
	}()
	assert.Nil(suite.T(), suite.Server.ListenAndServe())
}

func (suite *ServerTestSuite) TestListenAndServeWithBrokerServeErr() {
	suite.Broker.EXPECT().serve().Return(suite.TestErr)
	assert.NotNil(suite.T(), suite.Server.ListenAndServe())
}

func (suite *ServerTestSuite) TestSignalHandlerWithBrokerStopErr() {
	suite.Broker.EXPECT().serve()
	suite.Broker.EXPECT().stop().Return(suite.TestErr)
	suite.Broker.EXPECT().getNotify()
	go func() {
		suite.ServerCancel()
	}()
	assert.NotNil(suite.T(), suite.Server.ListenAndServe())
}

func (suite *ServerTestSuite) TestSignalHandlerWithBrokerNotifyErr() {
	suite.Broker.EXPECT().serve()
	suite.Broker.EXPECT().stop()
	ch := make(chan error)
	suite.Broker.EXPECT().getNotify().Return(ch)
	stopChan := make(chan bool)
	go func() {
		ch <- suite.TestErr
		stopChan <- true
	}()
	assert.NotNil(suite.T(), suite.Server.ListenAndServe())
	<-stopChan
}

func (suite *ServerTestSuite) TestSignalHandlerWithOnSig() {
	suite.Broker.EXPECT().serve()
	suite.Broker.EXPECT().stop()
	suite.Broker.EXPECT().getNotify()
	stopChan := make(chan bool)
	go func() {
		suite.Server.osSig <- syscall.SIGINT
		stopChan <- true
	}()
	assert.Nil(suite.T(), suite.Server.ListenAndServe())
	<-stopChan
}

func (suite *ServerTestSuite) TestSignalHandlerWithOnSigBrokerStopErr() {
	suite.Broker.EXPECT().serve()
	suite.Broker.EXPECT().stop().Return(suite.TestErr)
	suite.Broker.EXPECT().getNotify()
	stopChan := make(chan bool)
	go func() {
		suite.Server.osSig <- syscall.SIGINT
		stopChan <- true
	}()
	assert.NotNil(suite.T(), suite.Server.ListenAndServe())
	<-stopChan
}

func (suite *ServerTestSuite) TestPublish() {
	suite.Broker.EXPECT().Publish(gomock.Any())
	assert.Nil(suite.T(), suite.Server.Publish(1))
}

func (suite *ServerTestSuite) TestPublishErr() {
	suite.Broker.EXPECT().Publish(gomock.Any()).Return(suite.TestErr)
	assert.NotNil(suite.T(), suite.Server.Publish(1))
}
