package database

import (
	"database/sql"
	"encoding/json"
	"example/model"
	_ "github.com/lib/pq"
	"log"
)

// DB представляет соединение с базой данных PostgreSQL.
type DB struct {
	db *sql.DB
}

// NewDB создает новое соединение с базой данных.
func NewDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close закрывает соединение с базой данных.
func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func GetOrder(orderUID string, db *sql.DB) (*model.Order, error) {
	// Выполните SQL-запрос, чтобы получить заказ по orderUID из таблицы orders
	query := "SELECT * FROM orders WHERE order_uid = $1"
	row := db.QueryRow(query, orderUID)

	// Создайте переменную для хранения заказа
	var order model.Order

	// Парсинг данных из результата запроса и заполнение структуры Order
	err := row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.CustomerID,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SMID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, err
	}

	// Получите информацию о доставке из таблицы delivery
	delivery := model.Delivery{}
	deliveryQuery := "SELECT name, phone, zip, city, address, region, email FROM delivery WHERE order_uid = $1"
	err = db.QueryRow(deliveryQuery, orderUID).Scan(
		&delivery.Name,
		&delivery.Phone,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		return nil, err
	}
	order.Delivery = delivery

	// Получите информацию о платеже из таблицы payment
	payment := model.Payment{}
	paymentQuery := "SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payment WHERE order_uid = $1"
	err = db.QueryRow(paymentQuery, orderUID).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}
	order.Payment = payment

	// Получите информацию о товарах (Items) из таблицы order_items
	items := []model.Item{}
	itemsQuery := "SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM order_items WHERE order_uid = $1"
	rows, err := db.Query(itemsQuery, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		err := rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.RID,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NMID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	order.Items = items
	// Верните заказ
	return &order, nil
}

func ProcessData(data []byte) (*model.Order, error) {
	var orderData model.Order

	// Распаковываем JSON-данные в структуру Order
	if err := json.Unmarshal(data, &orderData); err != nil {
		log.Printf("Error unmarshaling JSON data: %v", err)
		return nil, err
	}

	return &orderData, nil
}

// SaveOrder сохраняет заказ в базе данных.
func SaveOrder(order *model.Order, db *sql.DB) error {
	// Начнем транзакцию, чтобы гарантировать целостность данных
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Выполним INSERT операции для каждой таблицы, связанной с заказом

	// Сохранение информации о заказе
	orderQuery := `
		INSERT INTO orders (order_uid, track_number, entry, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (order_uid) DO NOTHING`
	_, err = tx.Exec(orderQuery, order.OrderUID, order.TrackNumber, order.Entry, order.CustomerID, order.DeliveryService, order.Shardkey, order.SMID, order.DateCreated, order.OofShard)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Сохранение информации о доставке
	deliveryQuery := `
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO NOTHING`
	_, err = tx.Exec(deliveryQuery, order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Сохранение информации о платеже
	paymentQuery := `
		INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING`
	_, err = tx.Exec(paymentQuery, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Сохранение информации о товарах
	itemsQuery := `
		INSERT INTO order_items (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	for _, item := range order.Items {
		_, err := tx.Exec(itemsQuery, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale, item.Size, item.TotalPrice, item.NMID, item.Brand, item.Status)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Если все успешно, закрываем транзакцию
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
