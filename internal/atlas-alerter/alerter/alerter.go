// Package alerter provides functions used to run the Alerter service.
package alerter

//go:generate mockgen -source alerter.go -destination mock_alarmer_test.go -package alerter

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/thingspect/atlas/internal/atlas-alerter/config"
	"github.com/thingspect/atlas/pkg/alog"
	"github.com/thingspect/atlas/pkg/cache"
	"github.com/thingspect/atlas/pkg/dao"
	"github.com/thingspect/atlas/pkg/dao/alarm"
	"github.com/thingspect/atlas/pkg/dao/alert"
	"github.com/thingspect/atlas/pkg/dao/org"
	"github.com/thingspect/atlas/pkg/dao/user"
	"github.com/thingspect/atlas/pkg/notify"
	"github.com/thingspect/atlas/pkg/queue"
	"github.com/thingspect/proto/go/api"
)

// ServiceName provides consistent naming, including logs and metrics.
const ServiceName = "alerter"

// orger defines the methods provided by an org.DAO.
type orger interface {
	Read(ctx context.Context, orgID string) (*api.Org, error)
}

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
	orgDAO   orger
	alarmDAO alarmer
	userDAO  userer
	aleDAO   alerter
	cache    cache.Cacher

	aleQueue queue.Queuer
	eOutSub  queue.Subber

	notify notify.Notifier
}

// New builds a new Alerter and returns a reference to it and an error value.
func New(cfg *config.Config) (*Alerter, error) {
	// Set up database connection.
	pgRW, err := dao.NewPgDB(cfg.PgRwURI)
	if err != nil {
		return nil, err
	}

	pgRO, err := dao.NewPgDB(cfg.PgRoURI)
	if err != nil {
		return nil, err
	}

	// Set up cache connection.
	redis, err := cache.NewRedis(cfg.RedisHost + ":6379")
	if err != nil {
		return nil, err
	}

	// Set up Notifier. Allow a mock for local usage, but warn loudly.
	var n notify.Notifier
	if cfg.AppAPIKey == "" || cfg.SMSKeySecret == "" || cfg.EmailAPIKey == "" {
		alog.Error("New notify secrets not found, using notify.NewFake()")
		n = notify.NewFake()
	} else {
		n = notify.New(redis, cfg.AppAPIKey, cfg.SMSKeyID, cfg.SMSAccountID,
			cfg.SMSKeySecret, cfg.SMSPhone, cfg.EmailDomain, cfg.EmailAPIKey)
	}

	// Build the NSQ connection for consuming.
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
	eOutSub, err := nsq.Subscribe(cfg.NSQSubTopic)
	if err != nil {
		return nil, err
	}

	return &Alerter{
		orgDAO:   org.NewDAO(pgRW, pgRO),
		alarmDAO: alarm.NewDAO(pgRW, pgRO),
		userDAO:  user.NewDAO(pgRW, pgRO),
		aleDAO:   alert.NewDAO(pgRW, pgRO),
		cache:    redis,

		aleQueue: nsq,
		eOutSub:  eOutSub,

		notify: n,
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
