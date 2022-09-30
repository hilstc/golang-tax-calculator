package usecase

import (
	"github.com/hilstc/golang-tax-calculator/internal/order/entity"
)

// Returns the total value of the price + tax
type GetTotalOutputDTO struct {
	Total int
}

// Receives the OrderRepositoryInterface for the usecase
type GetTotalUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

// Creates a new usecase using the repository as a parameter
func NewGetTotalUseCase(orderRepository entity.OrderRepositoryInterface) *GetTotalUseCase {
	return &GetTotalUseCase{OrderRepository: orderRepository}
}

// Execute does not send any parameter, so it does not need an initial DTO,
// but it needs an exit DTO (GetTotalOutputDTO) for the return
func (c *GetTotalUseCase) Execute() (*GetTotalOutputDTO, error) {
	// Gets the total value from the repository
	total, err := c.OrderRepository.GetTotal()
	// If there is an error, returns an empty DTO and an error message
	if err != nil {
		return nil, err
	}
	// If there is no error, returns the total from the DTO and an empty error message
	return &GetTotalOutputDTO{Total: total}, nil
}
