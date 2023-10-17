package natsLib

import (
	"github.com/nats-io/nats.go"
)

type NatsMsgHandler func(*nats.Msg) error
