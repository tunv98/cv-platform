package port

import (
	"cv-platform/internal/domain"
)

type CVRepository interface {
	Create(cv *domain.CV) error
	Update(cv *domain.CV) error
	FindByID(id string) (*domain.CV, error)
	List(limit int, cursor string) ([]domain.CV, string, error)
}
