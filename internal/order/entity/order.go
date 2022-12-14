package entity

import "errors"

// Saves the data on a database and gets the total value after the tax calculation
type OrderRepositoryInterface interface {
	Save(order *Order) error
	GetTotal() (int, error)
}

// Defines the database structure
type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrder(id string, price float64, tax float64) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}

	err := order.IsValid()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) CalculateFinalPrice() error {
	o.FinalPrice = o.Price + o.Tax

	err := o.IsValid()
	if err != nil {
		return err
	}
	return nil
}

func (o Order) IsValid() error {
	if o.ID == "" {
		return errors.New("Invalid ID")
	}

	if o.Price == 0 {
		return errors.New("Invalid price")
	}

	if o.Tax == 0 {
		return errors.New("Invalid tax")
	}
	return nil
}
