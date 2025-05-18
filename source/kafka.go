package source

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// kafka数据源构建

type kafkaSource[E any] struct {
	*kafka.Consumer

	// 配置
	brokers []string `desc:"节点"`
	topics  []string `desc:"主题"`

	groupId          string `desc:"消费者组"`
	enableAutoCommit bool   `desc:"是否自动提交"`
}

func KafkaSourceBuilder[E any]() *kafkaSource[E] {
	return &kafkaSource[E]{
		enableAutoCommit: true, // 默认自动提交
	}
}

// 消息确认
func (source *kafkaSource[E]) CommitMessage(event *kafka.Message) (topicPartitions []kafka.TopicPartition, err error) {
	if !source.enableAutoCommit {
		topicPartitions, err = source.Consumer.CommitMessage(event)
	}

	return
}

func (source *kafkaSource[E]) SetEnableAutoCommit(enableAutoCommit bool) *kafkaSource[E] {
	source.enableAutoCommit = enableAutoCommit
	return source
}

func (source *kafkaSource[E]) SetBrokers(brokers []string) *kafkaSource[E] {
	source.brokers = brokers
	return source
}

func (source *kafkaSource[E]) SetTopics(topics []string) *kafkaSource[E] {
	source.topics = topics
	return source
}

func (source *kafkaSource[E]) SetGroupId(groupId string) *kafkaSource[E] {
	source.groupId = groupId
	return source
}

func (source *kafkaSource[E]) Build() *kafkaSource[E] {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": source.brokers, // 连接的broker
		"group.id":          source.groupId, // 消费者组

		// 指定kafka返回的数据最小/最大字节数
		"fetch.min.bytes": 1,                // 默认为1字节
		"fetch.max.bytes": 50 * 1024 * 1024, // 默认为50MB

		// 探活
		"session.timeout.ms":    10000, // 会话超市时间
		"heartbeat.interval.ms": 3000,  // 心跳间隔
	})
	if err != nil {
		panic(err)
	}

	// 设置消费主题
	if err = consumer.SubscribeTopics(source.topics, nil); err != nil {
		panic(err)
	}

	return &kafkaSource[E]{
		Consumer: consumer,
	}
}

type KafkaMessage[E any] struct {
	Event   *kafka.Message `desc:"Kafka原始消息"`
	Message E              `desc:"解析后的消息"`
}

func (source *kafkaSource[E]) ReadMessage() <-chan KafkaMessage[E] {
	queue := make(chan KafkaMessage[E], 100)

	go func() {
		defer close(queue)

		for {
			if event, err := source.Consumer.ReadMessage(-1); err != nil {
				// 记录日志
			} else {
				// 解析消息
				var msg E
				if err := json.Unmarshal(event.Value, &msg); err != nil {
					// 记录日志
					continue
				}

				queue <- KafkaMessage[E]{
					Event:   event,
					Message: msg,
				}
			}
		}

	}()

	return queue
}
