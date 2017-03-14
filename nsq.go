package main

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	nsq "github.com/nsqio/go-nsq"
)

// nsqConfig represents NSQ configuration read from file.
type nsqConfig struct {
	NSQDTCPAddrs     []string `mapstructure:"nsqd-tcp-address"`
	LookupdHTTPAddrs []string `mapstructure:"lookupd-http-address"`
	Topic            string
	Channel          string
	MaxInFlight      int `mapstructure:"max-in-flight"`
	MaxMsgs          int `mapstructure:"max-messages"`
}

// checkNSQConfig validates nsqConfig and sets defaults if necessary.
func checkNSQConfig(cfg *nsqConfig) error {
	if cfg.Channel == "" {
		rand.Seed(time.Now().UnixNano())
		cfg.Channel = fmt.Sprintf("tail%06d#ephemeral", rand.Int()%999999)
	}
	if cfg.Topic == "" {
		return errors.New("nsq/topic is required")
	}
	if len(cfg.NSQDTCPAddrs) == 0 && len(cfg.LookupdHTTPAddrs) == 0 {
		return errors.New("nsq/nsqd-tcp-address or nsq/lookupd-http-address required")
	}
	if len(cfg.NSQDTCPAddrs) > 0 && len(cfg.LookupdHTTPAddrs) > 0 {
		return errors.New("use nsq/nsqd-tcp-address or nsq/lookupd-http-address, not both")
	}
	return nil
}

// createNSQConsumer connects to NSQ and returns nsq.Consumer.
func createNSQConsumer(raw *nsqConfig, handler nsq.Handler) (*nsq.Consumer, error) {
	cfg := nsq.NewConfig()
	cfg.UserAgent = "x-msrv"
	cfg.MaxInFlight = raw.MaxInFlight
	consumer, err := nsq.NewConsumer(raw.Topic, raw.Channel, cfg)
	if err != nil {
		return nil, err
	}
	consumer.AddConcurrentHandlers(handler, runtime.NumCPU())
	err = consumer.ConnectToNSQDs(raw.NSQDTCPAddrs)
	if err != nil {
		return nil, err
	}
	err = consumer.ConnectToNSQLookupds(raw.LookupdHTTPAddrs)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
