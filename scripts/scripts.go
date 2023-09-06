package scripts

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
)

// Publisher структура для публикации данных в NATS Streaming.
type Publisher struct {
	sc      stan.Conn
	subject string
}

// NewPublisher создает новый экземпляр Publisher.
func NewPublisher(clusterID, clientID, natsURL, subject string) (*Publisher, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	return &Publisher{
		sc:      sc,
		subject: subject,
	}, nil
}

// Publish отправляет данные в канал NATS Streaming.
func (p *Publisher) Publish(data interface{}) error {
	// Преобразовываем данные в JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Публикуем данные в канал
	err = p.sc.Publish(p.subject, jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Close закрывает соединение с NATS Streaming.
func (p *Publisher) Close() {
	p.sc.Close()
}
