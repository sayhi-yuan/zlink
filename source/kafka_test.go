package source

import (
	"testing"
)

func TestKafkaSource(t *testing.T) {
	// 创建Kafka数据源
	source := KafkaSourceBuilder[string]().
		SetBrokers([]string{"localhost:9092"}).
		SetTopics([]string{"test-topic"}).
		SetGroupId("test-group").
		Build()

	// 读取消息
	for msg := range source.ReadMessage() {
		t.Logf("Received message: %s", msg)
	}
}
