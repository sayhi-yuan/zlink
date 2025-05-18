package source

import (
	"fmt"
	"strings"

	"github.com/nats-io/nats.go"
)

type NatsSource struct {
	hosts   []string   `desc:"节点"`
	conn    *nats.Conn `desc:"nats core连接"`
	subject string     `desc:"主题"`
}

func (source *NatsSource) SetHots(hosts []string) *NatsSource {
	source.hosts = hosts
	return source
}

func (source *NatsSource) SetSubject(subject string) *NatsSource {
	source.subject = subject
	return source
}

// 构建nats core
func (source *NatsSource) BuildCore() *NatsSource {
	conn, err := nats.Connect(strings.Join(source.hosts, ","))
	if err != nil {
		panic(fmt.Errorf("连接NATS异常, error:[%+v]", err))
	}

	return &NatsSource{
		conn: conn,
	}
}

// 订阅-发布
type NatsSubscribeSource struct {
	*NatsSource
}

type NatsMessage struct {
	Event   *nats.Msg
	Message string
}

func (source *NatsSubscribeSource) ReadMessage(max int) <-chan *nats.Msg {
	queue := make(chan *nats.Msg, max)

	_, err := source.conn.ChanSubscribe(source.subject, queue)
	if err != nil {
		// 错误记录
	}

	return queue
}

func NatsSubscribeSourceBuilder() *NatsSubscribeSource {
	return &NatsSubscribeSource{
		NatsSource: &NatsSource{},
	}
}

// 队列
type NatsQueueSource struct {
	*NatsSource

	name string `desc:"队列名称"`
}

func (source *NatsQueueSource) SetName(name string) *NatsQueueSource {
	source.name = name
	return source
}

func (source *NatsQueueSource) ReadMessage(max int) <-chan *nats.Msg {
	queue := make(chan *nats.Msg, max)

	_, err := source.conn.QueueSubscribe(source.subject, source.name, func(msg *nats.Msg) {
		queue <- msg
	})
	if err != nil {
		// 错误记录
	}

	return queue
}

func NatsQueueSourceBuilder() *NatsQueueSource {
	return &NatsQueueSource{
		NatsSource: &NatsSource{},
	}
}
