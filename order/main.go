package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type OrderItem struct {
	ItemId    int `bson:"itemId"`
	Quantity  int `bson:"quantity"`
	UnitPrice int `bson:"unitPrice"`
}

type OrderEvent struct {
	ID           string
	CustomerId   string
	CustomerName string
	OrderItem    []OrderItem
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// RabbitMQサーバーへの接続
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// チャネルの作成
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 送信先キューの宣言
	q, err := ch.QueueDeclare(
		"order", // 名前
		false,   // 永続的
		false,   // 未使用の場合削除
		false,   // 排他的
		false,   // 待機なし
		nil,     // 引数
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderEvent := OrderEvent{
		ID:           "test_id",
		CustomerId:   "test",
		CustomerName: "customer name",
		OrderItem:    []OrderItem{},
	}

	body, err2 := json.Marshal(orderEvent)
	failOnError(err2, "Failed to marshal")

	// メッセージをキューにパブリッシュ
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
