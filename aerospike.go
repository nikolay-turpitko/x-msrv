package main

import (
	"errors"
	"time"

	as "github.com/aerospike/aerospike-client-go"
)

// aspkConfig represents Aerospike configuration read from file.
type aspkConfig struct {
	Address struct {
		Host string
		Port int
	}
	Namespace string
	Set       string
}

// checkAspkConfig validates aspkConfig and sets defaults if necessary.
func checkAspkConfig(cfg *aspkConfig) error {
	if cfg.Address.Host == "" {
		return errors.New("aerospike/address/host is required")
	}
	if cfg.Address.Port == 0 {
		return errors.New("aerospike/address/port is required")
	}
	if cfg.Namespace == "" {
		return errors.New("aerospike/namespace is required")
	}
	if cfg.Set == "" {
		return errors.New("aerospike/set is required")
	}
	return nil
}

// createAspkClient creates Aerospike client.
func createAspkClient(cfg *aspkConfig) (*as.Client, error) {
	return as.NewClient(cfg.Address.Host, cfg.Address.Port)
}

// saveMsg saves message to DB.
func saveMsg(
	cfg *aspkConfig,
	client *as.Client,
	msg map[string]interface{}) error {
	// Because DB schema is not specified, and per sentence #7 of original
	// requirements, will use current time as record ID.
	id := time.Now().UnixNano()
	key, err := as.NewKey(cfg.Namespace, cfg.Set, id)
	if err != nil {
		return err
	}
	bins := make([]*as.Bin, 0, len(msg))
	bins = append(bins, as.NewBin("gid", id))
	for k, v := range msg {
		bins = append(bins, as.NewBin(k, v))
	}
	return client.PutBins(nil, key, bins...)
}
