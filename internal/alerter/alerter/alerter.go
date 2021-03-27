// Package alerter provides functions used to run the Alerter service.
package alerter

//go:generate mockgen -source alerter.go -destination mock_alarmer_test.go -package alerter

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thingspect/api/go/api"
	"github.com/thingspect/atlas/internal/alerter/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/alarm"
	"github.com/thingspect/atlas/pkg/dao/alert"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/queue"
)

const (
	ServiceName = "alerter"
)

// alarmer defines the methods provided by a alarm.DAO.
type alarmer interface {
	List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string,
		limit int32, ruleID string) ([]*api.Alarm, int32, error)
}

// userer defines the methods provided by a user.DAO.
type userer interface {
	ListByTags(ctx context.Context, orgID string, tags []string) ([]*api.User,
		error)
}

// alerter defines the methods provided by a alert.DAO.
type alerter interface {
	Create(ctx context.Context, alert *api.Alert) error
}

// Alerter holds references to the database and message broker connections.
type Alerter struct {
	alarmDAO alarmer
	userDAO  userer
	alertDAO alerter
	cache    cache.Cacher

	aleQueue queue.Queuer
	eOutSub  queue.Subber
}

// New builds a new Alerter and returns a reference to it and an error value.
func New(cfg *config.Config) (*Alerter, error) {
	// Set up database connection.
	pg, err := dao.NewPgDB(cfg.PgURI)
	if err != nil {
		return nil, err
	}

	// Set up cache connection.
	cache, err := cache.NewRedis(cfg.RedisHost + ":6379")
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
	eOutSub, err := nsq.Subscribe(cfg.NSQSubTopic)
	if err != nil {
		return nil, err
	}

	return &Alerter{
		alarmDAO: alarm.NewDAO(pg),
		userDAO:  user.NewDAO(pg),
		alertDAO: alert.NewDAO(pg),
		cache:    cache,

		aleQueue: nsq,
		eOutSub:  eOutSub,
	}, nil
}

// Serve starts the message alerters.
func (ale *Alerter) Serve(concurrency int) {
	for i := 0; i < concurrency; i++ {
		go ale.alertMessages()
	}

	// Handle graceful shutdown.
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)
	<-exitChan

	alog.Info("Serve received signal, exiting")
	if err := ale.eOutSub.Unsubscribe(); err != nil {
		alog.Errorf("Serve ale.eOutSub.Unsubscribe: %v", err)
	}
	ale.aleQueue.Disconnect()
}
