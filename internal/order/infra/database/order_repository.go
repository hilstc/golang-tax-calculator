package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hilstc/golang-tax-calculator/internal/order/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) GetTotal() (int, error) {
	var total int
	// Scans the memory address of "total" and adds the new value to the same variable
	err := r.Db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&total)

	// Returns 0 and an error message in case of errors
	if err != nil {
		return 0, err
	}
	// If there are no errors, returns the value of "total" and returns a null error message
	return total, nil
}
