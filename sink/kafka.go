package sink

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type kafkaSink struct {
	*kafka.Producer

	// 配置
	brokers []string `desc:"节点"`
	acks    string   `desc:"确认方式"`
}

func KafkaSinkBuilder() *kafkaSink {
	return &kafkaSink{
		acks: AcksLeader, // 默认等待leader确认
	}
}

func (sink *kafkaSink) SetBrokers(brokers []string) *kafkaSink {
	sink.brokers = brokers
	return sink
}

const (
	// 0 - 不等待确认，1 - 等待leader确认，all - 等待所有副本确认
	AcksNone   = "0"   // 不等待确认
	AcksLeader = "1"   // 等待leader确认
	AcksAll    = "all" // 等待所有副本确认
)

func (sink *kafkaSink) SetAcks(acks string) *kafkaSink {
	sink.acks = acks
	return sink
}

func (sink *kafkaSink) Build() *kafkaSink {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": sink.brokers, // 连接的broker

		"go.delivery.reports": false,
		// 重试配置
		"retries":          3,
		"retry.backoff.ms": 100, // 重试间隔

		// 0 - 不等待确认，1 - 等待leader确认，all - 等待所有副本确认
		"acks":      sink.acks, // 等待所有副本确认
		"linger.ms": 10,        // 等带批次加入的时间
	})
	if err != nil {
		panic(err)
	}

	return &kafkaSink{
		Producer: producer,
	}
}

func (sink *kafkaSink) Produce(topic string, data any) (err error) {
	msg, err := json.Marshal(data)
	if err != nil {
		return
	}

	err = sink.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: msg,
	}, nil)
	return
}
