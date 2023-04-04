// Package validator provides functions used to run the Validator service.
package validator

//go:generate mockgen -source validator.go -destination mock_devicer_test.go -package validator

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/atlas-validator/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/queue"
)

// ServiceName provides consistent naming, including logs and metrics.
const ServiceName = "validator"

// deviceExp provides the device DAO device expiration.
const deviceExp = 15 * time.Minute

// devicer defines the methods provided by a device.DAO.
type devicer interface {
	ReadByUniqID(ctx context.Context, uniqID string) (*api.Device, error)
}

// Validator holds references to the database and message broker connections.
type Validator struct {
	devDAO devicer

	valQueue     queue.Queuer
	vInSub       queue.Subber
	vOutPubTopic string
}

// New builds a new Validator and returns a reference to it and an error value.
func New(cfg *config.Config) (*Validator, error) {
	// Set up database connection.
	pg, err := dao.NewPgDB(cfg.PgURI)
	if err != nil {
		return nil, err
	}

	// Set up cache connection.
	redis, err := cache.NewRedis(cfg.RedisHost + ":6379")
	if err != nil {
		return nil, err
	}

	// Build the NSQ connection for consuming and publishing.
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, cfg.NSQLookupAddrs,
		cfg.NSQSubChannel)
	if err != nil {
		return nil, err
	}

	// Prime the queue before subscribing to allow for discovery by nsqlookupd.
	if err = nsq.Prime(cfg.NSQSubTopic); err != nil {
		return nil, err
	}

	// Subscribe to the topic.
	vInSub, err := nsq.Subscribe(cfg.NSQSubTopic)
	if err != nil {
		return nil, err
	}

	return &Validator{
		devDAO: device.NewDAO(pg, redis, deviceExp),

		valQueue:     nsq,
		vInSub:       vInSub,
		vOutPubTopic: cfg.NSQPubTopic,
	}, nil
}

// Serve starts the message validators.
func (val *Validator) Serve(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go val.validateMessages()
	}

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	if err := val.vInSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve val.vInSub.Unsubscribe: %v", err)
	}
	val.valQueue.Disconnect()
}
