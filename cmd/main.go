package main

import (
	"encoding/json"
	"example/internal/cache"
	"example/internal/database"
	"example/internal/nats"
	"example/internal/server"
	"example/model"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Config - структура для хранения настроек приложения из файла конфигурации.
type Config struct {
	NATS     NATSConfig     `json:"nats"`
	Database DatabaseConfig `json:"database"`
	Server   ServerConfig   `json:"server"`
}

// NATSConfig - структура для настроек NATS.
type NATSConfig struct {
	ClusterID string `json:"cluster_id"`
	ClientID  string `json:"client_id"`
	URL       string `json:"url"`
	Channel   string `json:"channel"`
}

// DatabaseConfig - структура для настроек базы данных.
type DatabaseConfig struct {
	URL string `json:"url"`
}

// ServerConfig - структура для настроек HTTP-сервера.
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	// Загрузка конфигурации
	config := loadConfig()

	// Создание соединения с NATS Streaming
	sc, err := nats.ConnectToNATS(config.NATS.ClusterID, config.NATS.ClientID, config.NATS.URL)

	if err != nil {
		log.Fatalf("Failed to connect to NATS Streaming: %v", err)
	}
	defer sc.Close()
	//// Создание и использование Publisher
	//publisher, err := scripts.NewPublisher(config.NATS.ClusterID, "your-publisher-client", config.NATS.URL, "your-subject")
	//if err != nil {
	//	log.Fatalf("Ошибка создания публикатора: %v", err)
	//}
	//defer publisher.Close()
	//
	//// Пример публикации данных
	//go func() {
	//	for {
	//		var data = generateData() // Замените на свою функцию для генерации данных
	//
	//		if err := publisher.Publish(data); err != nil {
	//			log.Printf("Ошибка при публикации данных: %v", err)
	//		} else {
	//			fmt.Printf("Отправлено сообщение: %v\n", data)
	//		}
	//
	//		time.Sleep(5 * time.Second) // Отправка данных каждые 5 секунд
	//	}
	//}()
	//

	// Создание соединения с PostgreSQL
	db, err := database.NewDB(config.Database.URL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Создание кэша
	newCache := cache.NewCache()

	// Запуск сервиса подписки на канал NATS Streaming
	go nats.SubscribeToChannel(sc, config.NATS.Channel, db, newCache)

	// Запуск HTTP-сервера для отображения данных
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		orderUID := r.URL.Query().Get("id")
		if orderUID == "" {
			http.Error(w, "Order ID is required", http.StatusBadRequest)
			return
		}

		// Попытка получить данные из кэша
		order := newCache.Get(orderUID)
		if order != nil {
			server.RenderOrderHTML(w, order)
			return
		}

		// Если данные отсутствуют в кэше, попытка получить из БД
		order, err := database.GetOrder(orderUID, db)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching order: %v", err), http.StatusInternalServerError)
			return
		}

		// Сохранение данных в кэше
		newCache.Set(orderUID, order)

		server.RenderOrderHTML(w, order)
	})

	http.HandleFunc("/create-order", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса, должен быть POST
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method. Use POST.", http.StatusMethodNotAllowed)
			return
		}

		// Распаковываем JSON-данные из тела запроса
		var orderData model.Order // Предположим, что у вас есть структура Order

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&orderData); err != nil {
			http.Error(w, fmt.Sprintf("Error decoding JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Сохраняем полученный заказ в базе данных
		if err := database.SaveOrder(&orderData, db); err != nil {
			http.Error(w, fmt.Sprintf("Error saving order to PostgreSQL: %v", err), http.StatusInternalServerError)
			return
		}
		// Генерируем JSON-ответ
		response := struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		}{
			Code:    http.StatusCreated,
			Message: "Order created successfully",
		}

		// Сериализуем структуру ответа в JSON
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error encoding JSON response: %v", err), http.StatusInternalServerError)
			return
		}

		// Устанавливаем Content-Type заголовок в application/json
		w.Header().Set("Content-Type", "application/json")

		// Пишем JSON-ответ в тело HTTP-ответа
		w.WriteHeader(http.StatusCreated)
		w.Write(jsonResponse)
	})

	serverAddr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	fmt.Printf("HTTP Server listening on %s...\n", serverAddr)
	log.Fatal(server.StartHTTPServer(serverAddr))
}

func loadConfig() *Config {
	config := Config{}

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Failed to decode config: %v", err)
	}

	return &config
}

func generateData() string {
	value := "hello"

	return value
}
