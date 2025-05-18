package sink

import "testing"

func TestKafkaSink(t *testing.T) {
	// 创建Kafka数据源
	sink := KafkaSinkBuilder().
		SetBrokers([]string{"localhost:9092"}).
		SetAcks(AcksLeader).
		Build()

	// 发送消息
	err := sink.Produce("test-topic", "Hello, Kafka!")
	if err != nil {
		t.Fatalf("Failed to produce message: %v", err)
	}
}
