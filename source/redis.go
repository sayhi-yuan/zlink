package source

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisQueueSource[E any] struct {
	*redis.Client

	queue []string `desc:"队列名称"`

	// ReadMessage
	timeout time.Duration `desc:"超时时间"`

	// ReadMessageOnce
	batchCount int64 `desc:"批量读取数量"`
}

func NewRedisQueueSourceBuilder[E any](client *redis.Client) *redisQueueSource[E] {
	return &redisQueueSource[E]{
		Client: client,

		timeout: 1 * time.Second, // 默认超时时间

		// 默认读取数量
		batchCount: -1, // 一次性读取数据时的默认配置
	}
}

func NewRedisQueueSourceOptionsBuilder[E any](options *redis.Options) *redisQueueSource[E] {
	client := redis.NewClient(options)

	// 检测
	if _, err := client.Ping(context.TODO()).Result(); err != nil {
		panic(err)
	}

	return &redisQueueSource[E]{
		Client: client,
	}
}
func (source *redisQueueSource[E]) SetQueue(queue []string) *redisQueueSource[E] {
	source.queue = queue
	return source
}

func (source *redisQueueSource[E]) Build() *redisQueueSource[E] {
	return source
}

type RedisMessage[E any] struct {
	Queue   string `desc:"队列名称"`
	Message E      `desc:"redis消息"`
	Error   error  `desc:"一次消息读取的error"`
}

func (source *redisQueueSource[E]) SetTimeout(timeout time.Duration) *redisQueueSource[E] {
	source.timeout = timeout
	return source
}

func (source *redisQueueSource[E]) ReadMessage(ctx context.Context, max int64) <-chan RedisMessage[E] {
	ch := make(chan RedisMessage[E], max)

	go func() {
		defer close(ch)

		for {
			dataList, err := source.BRPop(ctx, source.timeout, source.queue...).Result()
			if err != nil {
				ch <- RedisMessage[E]{Error: err}
				continue
			}

			for index, data := range dataList {
				if index%2 == 0 {
					continue
				}

				msg := RedisMessage[E]{
					Queue: dataList[index-1],
				}
				if err := json.Unmarshal([]byte(data), &msg.Message); err != nil {
					msg.Error = err
				}

				ch <- msg
			}
		}
	}()

	return ch
}

// ReadMessageOnce 读取一次消息
// 默认全部读区
func (source *redisQueueSource[E]) SetBatchCount(batchCount int64) *redisQueueSource[E] {
	source.batchCount = batchCount
	return source
}

// 当max为<=0时，一次性读取完
// 否则读取max条数据
func (source *redisQueueSource[E]) ReadMessageOnce(ctx context.Context, max int64) <-chan RedisMessage[E] {
	ch := make(chan RedisMessage[E], max)

	go func() {
		defer close(ch)

		for _, queue := range source.queue {
			dataList, err := source.LRange(ctx, queue, 0, source.batchCount).Result()
			if err != nil {
				ch <- RedisMessage[E]{
					Queue: queue,
					Error: err,
				}
				continue
			}

			for _, data := range dataList {
				msg := RedisMessage[E]{
					Queue: queue,
				}
				if err := json.Unmarshal([]byte(data), &msg.Message); err != nil {
					msg.Error = err
				}

				ch <- msg
			}
		}
	}()

	return ch
}
