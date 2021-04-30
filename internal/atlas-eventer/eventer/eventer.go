// Package eventer provides functions used to run the Eventer service.
package eventer

//go:generate mockgen -source eventer.go -destination mock_ruler_test.go -package eventer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/api/go/common"
	"github.com/thingspect/atlas/internal/atlas-eventer/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/event"
	"github.com/thingspect/atlas/pkg/dao/rule"
	"github.com/thingspect/atlas/pkg/queue"
)

const (
	ServiceName = "eventer"
)

// ruler defines the methods provided by a rule.DAO.
type ruler interface {
	ListByTags(ctx context.Context, orgID string, attr string,
		deviceTags []string) ([]*common.Rule, error)
}

// eventer defines the methods provided by an event.DAO.
type eventer interface {
	Create(ctx context.Context, event *api.Event) error
}

// Eventer holds references to the database and message broker connections.
type Eventer struct {
	ruleDAO  ruler
	eventDAO eventer

	evQueue      queue.Queuer
	vOutSub      queue.Subber
	eOutPubTopic string
}

// New builds a new Eventer and returns a reference to it and an error value.
func New(cfg *config.Config) (*Eventer, error) {
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

	// Prime the queue before subscribing to allow for discovery by nsqlookupd.
	if err = nsq.Prime(cfg.NSQSubTopic); err != nil {
		return nil, err
	}

	// Subscribe to the topic.
	vOutSub, err := nsq.Subscribe(cfg.NSQSubTopic)
	if err != nil {
		return nil, err
	}

	return &Eventer{
		ruleDAO:  rule.NewDAO(pg),
		eventDAO: event.NewDAO(pg),

		evQueue:      nsq,
		vOutSub:      vOutSub,
		eOutPubTopic: cfg.NSQPubTopic,
	}, nil
}

// Serve starts the message eventers.
func (ev *Eventer) Serve(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go ev.eventMessages()
	}

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	if err := ev.vOutSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve ev.vOutSub.Unsubscribe: %v", err)
	}
	ev.evQueue.Disconnect()
}
