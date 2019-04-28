package mqb

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type BrokerTestSuite struct {
	suite.Suite
}

func TestBrokerTestSuite(t *testing.T) {
	suite.Run(t, new(BrokerTestSuite))
}

func (suite *BrokerTestSuite) TestOptions() {
	wantUrl := "test dummy url"
	wantReconnects := 32
	wantReconnectTimeout := time.Second * 3
	opts := defaultBrokerOptions
	URL(wantUrl)(&opts)
	MaxReconnects(wantReconnects)(&opts)
	ReconnectTimeout(wantReconnectTimeout)(&opts)
	assert.Equal(suite.T(), wantUrl, opts.url)
	assert.Equal(suite.T(), wantReconnects, opts.maxReconnects)
	assert.Equal(suite.T(), wantReconnectTimeout, opts.reconnectTimeout)

}
