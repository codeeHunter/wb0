// internal/cache/cache.go

package cache

import (
	"database/sql"
	"encoding/json"
	"example/model"
	"sync"
	"time"
)

// Cache представляет in-memory кэш для данных заказов.
type Cache struct {
	mu                 sync.RWMutex
	cache              map[string]interface{}
	expiration         map[string]time.Time // Добавленный словарь для отслеживания времени истечения кэша
	expirationDuration time.Duration        // Время истечения кэша
}

// NewCache создает новый экземпляр кэша.
func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]interface{}),
	}
}

func RestoreCacheFromDB(db *sql.DB, cache *Cache) error {
	rows, err := db.Query("SELECT order_uid, data FROM orders")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var orderUID string
		var orderData []byte

		if err := rows.Scan(&orderUID, &orderData); err != nil {
			return err
		}

		// Распаковываем данные и сохраняем их в кэше
		var order model.Order
		if err := json.Unmarshal(orderData, &order); err != nil {
			return err
		}

		cache.Set(orderUID, &order)
	}

	return nil
}

// Get получает данные из кэша по ключу.
func (c *Cache) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache[key]
}

// Set сохраняет данные в кэше по ключу.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}
