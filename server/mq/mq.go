package mq

var cpMq = NewConfigPushMQ()

func ConfigPush() *ConfigPushMQ {
	return cpMq
}
