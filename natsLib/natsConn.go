package natsLib

import (
	"errors"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap/zapcore"
	"goToolkit/logLib"
	"math/rand"
	"time"
)

func connect(name, user, passwd string, serverList []string) (conn *nats.Conn, err error) {
	options := nats.GetDefaultOptions()
	options.PingInterval = 20 * time.Second
	options.MaxPingsOut = 3

	options.CustomReconnectDelayCB = func(attempts int) time.Duration {
		logLib.Zap().Debug("Reconnect attempts has been tried.",
			zapcore.Field{Key: "method", String: "natsLib.connect"},
			zapcore.Field{Key: "attempts", Integer: int64(attempts)})
		sleepMS := rand.Intn(3000)
		return time.Duration(sleepMS) * time.Millisecond
	}

	options.ClosedCB = func(conn *nats.Conn) {
		logLib.Zap().Debug("Connection is closed.", zapcore.Field{Key: "url", String: conn.Opts.Url})
	}

	options.DisconnectedErrCB = func(conn *nats.Conn, err error) {
		logLib.Zap().Debug("Connection disconnected.", zapcore.Field{Key: "url", String: conn.Opts.Url})
	}

	options.ReconnectedCB = func(conn *nats.Conn) {
		logLib.Zap().Debug("Reconnect to server.", zapcore.Field{Key: "url", String: conn.ConnectedUrl()})
	}

	options.DiscoveredServersCB = func(conn *nats.Conn) {
		logLib.Zap().Debug("Discover new server.",
			zapcore.Field{Key: "method", String: "natsLib.connect"},
			zapcore.Field{Key: "known servers", Interface: conn.Servers()},
			zapcore.Field{Key: "discovered servers", Interface: conn.DiscoveredServers()})
	}

	options.AsyncErrorCB = func(conn *nats.Conn, sub *nats.Subscription, err error) {
		logLib.Zap().Debug("Async error occurred.",
			zapcore.Field{Key: "method", String: "natsLib.connect"},
			zapcore.Field{Key: "url", String: conn.Opts.Url},
			zapcore.Field{Key: "subject", String: sub.Subject},
			zapcore.Field{Key: "error", String: err.Error()})

		if errors.Is(err, nats.ErrSlowConsumer) {
			pendingMsgNum, _, err := sub.Pending()
			if err != nil {
				logLib.Zap().Debug("Couldn't get pending message nums.",
					zapcore.Field{Key: "method", String: "natsLib.connect"},
					zapcore.Field{Key: "error", String: err.Error()})
				return
			}
			logLib.Zap().Debug("Falling behind with pending messages on subject.",
				zapcore.Field{Key: "method", String: "natsLib.connect"},
				zapcore.Field{Key: "subject", String: sub.Subject},
				zapcore.Field{Key: "pendingMsgNum", Integer: int64(pendingMsgNum)})
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
