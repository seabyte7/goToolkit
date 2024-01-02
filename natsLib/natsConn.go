package natsLib

import (
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/seabyte7/goToolkit/logLib"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

func connect(name, user, passwd string, serverList []string) (conn *nats.Conn, err error) {
	options := nats.GetDefaultOptions()
	options.PingInterval = 20 * time.Second
	options.MaxPingsOut = 3

	options.CustomReconnectDelayCB = func(attempts int) time.Duration {
		logLib.Zap().Debug("Reconnect attempts has been tried.",
			zap.String("method", "natsLib.connect"),
			zap.Int64("attempts", int64(attempts)))

		sleepMS := rand.Intn(3000)
		return time.Duration(sleepMS) * time.Millisecond
	}

	options.ClosedCB = func(conn *nats.Conn) {
		logLib.Zap().Debug("Connection is closed.", zap.String("url", conn.Opts.Url), zap.Any("servers", conn.Servers()))
	}

	options.DisconnectedErrCB = func(conn *nats.Conn, err error) {
		logLib.Zap().Debug("Connection disconnected.", zap.String("url", conn.Opts.Url), zap.Any("servers", conn.Servers()))
	}

	options.ReconnectedCB = func(conn *nats.Conn) {
		logLib.Zap().Debug("Reconnect to server.", zap.String("url", conn.Opts.Url), zap.Any("servers", conn.Servers()))
	}

	options.DiscoveredServersCB = func(conn *nats.Conn) {
		logLib.Zap().Debug("Discover new server.",
			zap.String("method", "natsLib.connect"),
			zap.Any("known servers", conn.Servers()),
			zap.Any("discovered servers", conn.DiscoveredServers()))
	}

	options.AsyncErrorCB = func(conn *nats.Conn, sub *nats.Subscription, err error) {
		logLib.Zap().Debug("Async error occurred.",
			zap.String("method", "natsLib.connect"),
			zap.String("url", conn.Opts.Url),
			zap.String("subject", sub.Subject),
			zap.String("error", err.Error()))

		if errors.Is(err, nats.ErrSlowConsumer) {
			pendingMsgNum, _, err := sub.Pending()
			if err != nil {
				logLib.Zap().Debug("Couldn't get pending message nums.",
					zap.String("method", "natsLib.connect"),
					zap.String("subject", sub.Subject),
					zap.String("error", err.Error()))
				return
			}
			logLib.Zap().Debug("Falling behind with pending messages on subject.",
				zap.String("method", "natsLib.connect"),
				zap.String("url", conn.Opts.Url),
				zap.Any("servers", conn.Servers()),
				zap.String("subject", sub.Subject),
				zap.Int("pendingMsgNum", pendingMsgNum))
		}
	}

	options.Name = name
	options.User = user
	options.Password = passwd
	options.Servers = serverList

	conn, err = options.Connect()
	if err != nil {
		return
	}

	return
}
