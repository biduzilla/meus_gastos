package data

import "time"

type Transaction struct {
	ID          int64
	CreatedAt   time.Time
	Deleted     bool
	Version     int
	User        *User
	Category    *Category
	Description string
	Amount      float64
	Type        TypeCategoria
}

type TransactionDTO struct {
	ID          *int64       `json:"transaction_id"`
	Version     *int         `json:"version"`
	User        *UserDTO     `json:"user"`
	Category    *CategoryDTO `json:"category"`
	Description *string      `json:"description"`
	Amount      *float64     `json:"amount"`
	Type        *string      `json:"type"`
}

func (t *Transaction) toDTO() *TransactionDTO {
	var typeStr string
	if t.Type != 0 {
		typeStr = t.Type.String()
	}

	var id *int64
	if t.ID != 0 {
		id = &t.ID
	}

	var version *int
	if t.Version != 0 {
		version = &t.Version
	}

	var user *UserDTO
	if t.User != nil {
		user = t.User.ToDTO()
	}

	var category *CategoryDTO
	if t.Category != nil {
		category = t.Category.ToDTO()
	}

	var description *string
	if t.Description != "" {
		description = &t.Description
	}

	var amount *float64
	if t.Amount != 0 {
		amount = &t.Amount
	}

	return &TransactionDTO{
		ID:          id,
		Version:     version,
		User:        user,
		Category:    category,
		Description: description,
		Amount:      amount,
		Type:        &typeStr,
	}
}

func (t *TransactionDTO) toModel() *Transaction {
	var id *int64
	if t.ID != nil {
		id = t.ID
	}

	var version *int
	if t.Version != nil {
		version = t.Version
	}

	var user *User
	if t.User != nil {
		user = t.User.ToModel()
	}

	var category *Category
	if t.Category != nil {
		category = t.Category.ToModel()
	}

	var description *string
	if t.Description != nil {
		description = t.Description
	}

	var tipo TypeCategoria
	if t.Type != nil {
		tipo = TypeCategoriaFromString(*t.Type)
	}

	var amount *float64
	if t.Amount != nil {
		amount = t.Amount
	}

	return &Transaction{
		ID:          *id,
		Version:     *version,
		User:        user,
		Category:    category,
		Description: *description,
		Amount:      *amount,
		Type:        tipo,
	}
}
