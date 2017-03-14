package main

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync/atomic"
	"syscall"
	"time"

	as "github.com/aerospike/aerospike-client-go"
	"github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

var configFilePath = flag.String("cfg", "/etc/x-msrv/x-msrv.yml", "path to config file")

// parseConfig parses configuration file and returns validated configuration.
func parseConfig() (time.Duration, *nsqConfig, *aspkConfig, error) {
	viper.AddConfigPath(filepath.Dir(*configFilePath))
	base := filepath.Base(*configFilePath)
	ext := filepath.Ext(base)
	viper.SetConfigName(base[:len(base)-len(ext)])
	err := viper.ReadInConfig()
	if err != nil {
		return 0, nil, nil, err
	}
	var st string
	err = viper.UnmarshalKey("timeout", &st)
	if err != nil {
		return 0, nil, nil, err
	}
	timeout, err := time.ParseDuration(st)
	if err != nil {
		return 0, nil, nil, err
	}
	nsqCfg := &nsqConfig{}
	err = viper.UnmarshalKey("nsq", nsqCfg)
	if err != nil {
		return 0, nil, nil, err
	}
	err = checkNSQConfig(nsqCfg)
	if err != nil {
		return 0, nil, nil, err
	}
	aspkCfg := &aspkConfig{}
	err = viper.UnmarshalKey("aerospike", aspkCfg)
	if err != nil {
		return 0, nil, nil, err
	}
	err = checkAspkConfig(aspkCfg)
	if err != nil {
		return 0, nil, nil, err
	}
	return timeout, nsqCfg, aspkCfg, nil
}

// msgHandler implements nsq.Handler interface, processes incoming messages.
type msgHandler struct {
	counter    int64
	maxMsgs    int64
	gotMsgs    chan struct{}
	aspkCfg    *aspkConfig
	aspkClient *as.Client
}

func (h *msgHandler) HandleMessage(m *nsq.Message) error {
	// Since handler can be invoked concurrently, atomic operation used.
	// When enough messages received, handler sends signal to main goroutine
	// to stop consumer and finish execution. Some handlers can be still invoked
	// concurrently after that and they will see counter grater than allowed
	// number of messages to process. In that case, handler returns error to
	// signal NSQ consumer to re-queue message.
	c := atomic.AddInt64(&h.counter, 1)
	switch {
	case c == h.maxMsgs:
		h.gotMsgs <- struct{}{}
	case c > h.maxMsgs:
		return errors.New("enough is enough")
	}

	// Unmarshal to map[string]interface{} validates JSON messages.
	// For known message format, struct with tags can be used.
	var msg map[string]interface{}
	if err := json.Unmarshal(m.Body, &msg); err != nil {
		log.Printf("Error unmarshaling message `%s`, error: `%v`", m.Body, err)
		// Return no error, because there is no sence to receive this message again.
		return nil
	}

	return saveMsg(h.aspkCfg, h.aspkClient, msg)
}

func main() {
	log.Println("Service x-msrv started")
	defer log.Println("Service x-msrv exited")
	flag.Parse()

	timeout, nsqCfg, aspkCfg, err := parseConfig()
	if err != nil {
		log.Fatal("Configuration error: ", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	maxMsgsChan := make(chan struct{}, 1)

	aspkClient, err := createAspkClient(aspkCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer aspkClient.Close()

	// Registering message handlers, which can be invoked concurrency.
	h := &msgHandler{
		counter:    0,
		maxMsgs:    int64(nsqCfg.MaxMsgs),
		gotMsgs:    maxMsgsChan,
		aspkCfg:    aspkCfg,
		aspkClient: aspkClient,
	}
	consumer, err := createNSQConsumer(nsqCfg, h)
	if err != nil {
		// defer won't be executed
		aspkClient.Close()
		log.Fatal(err)
	}

	// Service designed to be executed by external timer, so it should complete
	// faster, than scheduler start it again. But NSQ consumer blocks waiting
	// for messages. So, let's use timeout to finish execution.
	timer := time.After(timeout)

	// Handling termination signal, signal about enough messages received and
	// awaiting customer termination.
	for {
		select {
		case <-consumer.StopChan:
			log.Println("Service x-msrv stopped")
			return
		case <-sigChan:
			log.Println("Service x-msrv caught termination signal")
			consumer.Stop()
		case <-maxMsgsChan:
			log.Println("Service x-msrv received enough messages (per configuration)")
			consumer.Stop()
		case <-timer:
			log.Println("Service x-msrv timed out")
			consumer.Stop()
		}
	}
}
