package kafka

import (
	"encoding/json"
	"log/slog"
	"time"

	"context"

	"github.com/IBM/sarama"
)

const retryNumber = 5

type Consumer interface {
	Start() error
}

// ConsumerHandler 消费者
type ConsumerHandler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, data T) error
}

func (h ConsumerHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h ConsumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	allMsg := claim.Messages()
	for msg := range allMsg {
		var data T
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			slog.Error("Unmarshal Message Fail",
				"err", err.Error(), "topic", msg.Topic, "offset", msg.Offset, "partition", msg.Partition)
			continue
		}
		done := true
		if err := h.fn(msg, data); err != nil {
			// TODO 重试
			done = false
			for range retryNumber {
				if err := h.fn(msg, data); err == nil {
					done = true
					break
				}
			}
			slog.Error("Consume Message Fail",
				"err", err.Error(), "topic", msg.Topic, "offset", msg.Offset, "partition", msg.Partition)
		}
		if done {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

func NewConsumerHandler[T any](fn func(msg *sarama.ConsumerMessage, data T) error) *ConsumerHandler[T] {
	return &ConsumerHandler[T]{
		fn: fn,
	}
}

// BatchConsumerHandler 消费者
type BatchConsumerHandler[T any] struct {
	fn func(allMsg []*sarama.ConsumerMessage, allData []T) error
}

const (
	batchSize = 10
	duration  = time.Second
)

func (h BatchConsumerHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h BatchConsumerHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h BatchConsumerHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("Consumer claim start")
	allMsg := claim.Messages()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), duration)
		done := false
		sendMsg := make([]*sarama.ConsumerMessage, 0, batchSize)
		sendData := make([]T, 0, batchSize)
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-allMsg:
				if !ok {
					cancel()
					return nil // 消费者被关闭
				}

				var data T
				if err := json.Unmarshal(msg.Value, &data); err != nil {
					slog.Error("Data Unmarshal Fail",
						"Error", err.Error(), "topic", msg.Topic, "offset", msg.Offset, "partition", msg.Partition)
					continue
				}
				sendMsg = append(sendMsg, msg)
				sendData = append(sendData, data)
			}
		}
		_ = h.fn(sendMsg, sendData)
		for _, msg := range sendMsg {
			session.MarkMessage(msg, "")
		}
		cancel()
	}
}

func NewBatchConsumerHandler[T any](fn func(allMsg []*sarama.ConsumerMessage, allData []T) error) *BatchConsumerHandler[T] {
	return &BatchConsumerHandler[T]{
		fn: fn,
	}
}
