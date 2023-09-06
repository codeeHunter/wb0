package nats

import (
	"database/sql"
	"example/internal/cache"
	"example/internal/database"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
)

// ConnectToNATS подключается к серверу NATS Streaming.
func ConnectToNATS(clusterID, clientID, natsURL string) (stan.Conn, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}
	return sc, nil
}

// SubscribeToChannel подписывается на канал NATS Streaming.
func SubscribeToChannel(sc stan.Conn, channelName string, db *sql.DB, cache *cache.Cache) {
	_, err := sc.Subscribe(channelName, func(msg *stan.Msg) {
		// Обработка полученного сообщения
		orderData, _ := database.ProcessData(msg.Data)

		// Запись данных в PostgreSQL
		if err := database.SaveOrder(orderData, db); err != nil {
			log.Printf("Error saving order to PostgreSQL: %v", err)

			// Важно: В случае ошибки, не обновляйте данные в кэше, чтобы избежать несогласованных данных.
			return
		}

		// Обновление данных в кэше
		cache.Set(orderData.OrderUID, orderData)
	})
	if err != nil {
		log.Fatalf("Error subscribing to channel: %v", err)
	}
	fmt.Println("Subscribed to channel:", channelName)
}
