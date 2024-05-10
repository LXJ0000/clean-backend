package event

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/LXJ0000/clean-backend/utils/kafka"

	"github.com/IBM/sarama"
	"github.com/LXJ0000/clean-backend/internal/domain"
)

type ReadEvent struct {
	UserID int64
	PostID int64
}

type SyncReadEventProducer struct {
	producer sarama.SyncProducer
}

func NewSyncProducer(producer sarama.SyncProducer) *SyncReadEventProducer {
	return &SyncReadEventProducer{producer: producer}
}

func (s SyncReadEventProducer) ProduceReadEvent(c context.Context, event ReadEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	msg := sarama.ProducerMessage{
		Topic: "post_read",
		Value: sarama.ByteEncoder(data),
	}
	slog.Debug("消息组装完毕... 准备发送噜~")
	_, _, err = s.producer.SendMessage(&msg)
	return err
}

// BatchSyncReadEventConsumer 批量消费者
type BatchSyncReadEventConsumer struct {
	client sarama.Client
	repo   domain.InteractionRepository
}

func NewBatchSyncReadEventConsumer(client sarama.Client, repo domain.InteractionRepository) *BatchSyncReadEventConsumer {
	return &BatchSyncReadEventConsumer{client: client, repo: repo}
}

func (c *BatchSyncReadEventConsumer) Start() error {
	config, err := sarama.NewConsumerGroupFromClient("interactive", c.client)
	if err != nil {
		return err
	}
	go func() {
		if err := config.Consume(context.Background(), []string{"post_read"}, kafka.NewBatchConsumerHandler[ReadEvent](c.Consumer)); err != nil {
			slog.Error("Consumer Start Fail", "topic", []string{"post_read"})
		}
	}()
	return nil
}

func (c *BatchSyncReadEventConsumer) Consumer(msg []*sarama.ConsumerMessage, event []ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	biz := make([]string, 0, len(event))
	id := make([]int64, 0, len(event))
	for _, e := range event {
		biz = append(biz, domain.BizPost)
		id = append(id, e.PostID)
		//go func() { // 批量性能较好，一个事务与多个事务的区别
		//	if err := c.repo.IncrReadCount(ctx, domain.BizPost, e.PostID); err != nil {
		//		slog.Error("Consumer Message Fail", "Error", err.Error(), "biz", domain.BizPost, "post_id", e.PostID)
		//	}
		//}()
	}
	return c.repo.BatchIncrReadCount(ctx, biz, id)
}
