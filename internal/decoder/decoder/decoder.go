// Package decoder provides functions used to run the Decoder service.
package decoder

//go:generate mockgen -source decoder.go -destination mock_devicer_test.go -package decoder

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/decoder/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/device"
	"github.com/thingspect/atlas/pkg/decode/registry"
	"github.com/thingspect/atlas/pkg/queue"
)

const (
	ServiceName = "decoder"
)

// devicer defines the methods provided by a device.DAO.
type devicer interface {
	ReadByUniqID(ctx context.Context, uniqID string) (*common.Device, error)
}

// Decoder holds references to the database and message broker connections.
type Decoder struct {
	devDAO devicer
	reg    *registry.Registry

	decQueue    queue.Queuer
	dInSub      queue.Subber
	vInPubTopic string
}

// New builds a new Decoder and returns a reference to it and an error decue.
func New(cfg *config.Config) (*Decoder, error) {
	// Set up database connection.
	pg, err := dao.NewPgDB(cfg.PgURI)
	if err != nil {
		return nil, err
	}

	// Build the NSQ connection for consuming and publishing.
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, cfg.NSQLookupAddrs,
		cfg.NSQSubChannel, queue.DefaultNSQRequeueDelay)
	if err != nil {
		return nil, err
	}

	// Subscribe to the topic.
	dInSub, err := nsq.Subscribe(cfg.NSQSubTopic)
	if err != nil {
		return nil, err
	}

	return &Decoder{
		devDAO: device.NewDAO(pg),
		reg:    registry.New(),

		decQueue:    nsq,
		dInSub:      dInSub,
		vInPubTopic: cfg.NSQPubTopic,
	}, nil
}

// Serve starts the message decoders.
func (dec *Decoder) Serve(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go dec.decodeMessages()
	}

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	if err := dec.dInSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve dec.subber.Unsubscribe: %v", err)
	}
	dec.decQueue.Disconnect()
}
