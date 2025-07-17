package product

import (
	"4-order-api/pkg/req"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string         `json:"name" gorm:"uniqueIndex:idx_name" validate:"required"`
	Description string         `json:"description"`
	Images      pq.StringArray `json:"images" gorm:"type:text[]" validate:"min=1,dive,required"`
}

func (p *Product) Validate() error {
	return req.IsValid(*p)
}

func NewProduct(name, description string, images []string) (*Product, error) {
	product := &Product{
		Name:        name,
		Description: description,
		Images:      images,
	}
	if err := product.Validate(); err != nil {
		return nil, err
	}
	return product, nil
}
