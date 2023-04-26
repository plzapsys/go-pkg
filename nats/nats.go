package nats

import (
	"time"

	"github.com/plzapsys/go-pkg/logger"

	"github.com/nats-io/stan.go"
)

const (
	connectWait        = time.Second * 30
	pubAckWait         = time.Second * 30
	interval           = 10
	maxOut             = 5
	maxPubAcksInflight = 25
)

// Nats config
type Config struct {
	URL       string
	ClusterID string
	ClientID  string
}

func NewNatsConnect(cfg *Config, log logger.ILogger) (stan.Conn, error) {
	return stan.Connect(
		cfg.ClusterID,
		cfg.ClientID,
		stan.ConnectWait(connectWait),
		stan.PubAckWait(pubAckWait),
		stan.NatsURL(cfg.URL),
		stan.Pings(interval, maxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}),
		stan.MaxPubAcksInflight(maxPubAcksInflight),
	)
}
