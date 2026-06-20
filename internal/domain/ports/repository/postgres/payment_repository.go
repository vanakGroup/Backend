package postgres

import "github.com/Vanak/Backend/internal/domain/entities"

type PaymentRepository interface {
	Create(*entities.MockPayment) error
}
