package mq

import (
	rc "github.com/slyerr/econfig/core/rest/client"
	"github.com/slyerr/econfig/core/utils"
	"github.com/slyerr/econfig/server/store/client"
	"github.com/slyerr/verifier"
	"go.uber.org/zap"
)

type ConfigPushMsg struct {
	Key   string
	Value string
}

type ConfigPushMQ struct {
	ch      chan ConfigPushMsg
	rc      *rc.RestClient
	started bool
}

func NewConfigPushMQ() *ConfigPushMQ {
	rc := rc.NewRestClient()
	ch := make(chan ConfigPushMsg, 10)
	return &ConfigPushMQ{ch: ch, rc: rc, started: false}
}

func (mq *ConfigPushMQ) Start() {
	if mq.started {
		return
	}

	go func() {
		for msg := range mq.ch {
			msg := msg
			go mq.consume(msg)
		}
	}()
	mq.started = true
}

func (mq *ConfigPushMQ) Produce(msg ConfigPushMsg) {
	if !mq.started {
		panic("econfig server's config push mq not start")
	}

	mq.ch <- msg
}

func (mq *ConfigPushMQ) consume(msg ConfigPushMsg) {
	zap.S().Debugf("econfig server's config push msg: %+v", msg)

	key, err := utils.CheckConfigKey(msg.Key)
	if err != nil {
		return
	}

	value := utils.CleanConfigValue(msg.Value)

	cc, err := client.Store().Get(key)
	if err != nil {
		zap.S().Errorf("econfig server's config push get [%v] client info error: %+v", key, err)
		return
	}
	if len(cc) == 0 {
		return
	}

	for _, c := range cc {
		url := utils.CompleUrl(c.Host, c.PushUrl)
		if err := verifier.S().NotBlankN(url, "url"); err != nil {
			continue
		}

		go func() {
			if _, err := mq.rc.Put(url, value); err != nil {
				zap.S().Errorf("econfig server's config push error: %+v", err)
			}
		}()
	}
}
