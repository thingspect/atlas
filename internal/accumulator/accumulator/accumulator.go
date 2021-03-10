// Package accumulator provides functions used to run the Accumulator service.
package accumulator

//go:generate mockgen -source accumulator.go -destination mock_datapointer_test.go -package accumulator

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/accumulator/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/datapoint"
	"github.com/thingspect/atlas/pkg/queue"
)

const (
	ServiceName = "accumulator"
)

// datapointer defines the methods provided by a datapoint.DAO.
type datapointer interface {
	Create(ctx context.Context, point *common.DataPoint, orgID string) error
}

// Accumulator holds references to the database and message broker connections.
type Accumulator struct {
	dpDAO datapointer

	accQueue queue.Queuer
	vOutSub  queue.Subber
}

// New builds a new Accumulator and returns a reference to it and an error
// value.
func New(cfg *config.Config) (*Accumulator, error) {
	// Set up database connection.
	pg, err := dao.NewPgDB(cfg.PgURI)
	if err != nil {
		return nil, err
	}

	// Build the NSQ connection for consuming.
	nsq, err := queue.NewNSQ(cfg.NSQPubAddr, cfg.NSQLookupAddrs,
		cfg.NSQSubChannel, queue.DefaultNSQRequeueDelay)
	if err != nil {
		return nil, err
	}

	// Subscribe to the topic.
	vOutSub, err := nsq.Subscribe(cfg.NSQSubTopic)
	if err != nil {
		return nil, err
	}

	return &Accumulator{
		dpDAO: datapoint.NewDAO(pg),

		accQueue: nsq,
		vOutSub:  vOutSub,
	}, nil
}

// Serve starts the message accumulators.
func (acc *Accumulator) Serve(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go acc.accumulateMessages()
	}

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	if err := acc.vOutSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve acc.vOutSub.Unsubscribe: %v", err)
	}
	acc.accQueue.Disconnect()
}
