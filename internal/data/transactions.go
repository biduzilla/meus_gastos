package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"meus_gastos/internal/validator"
	"time"
)

type TransactionModel struct {
	DB *sql.DB
}

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

func (m TransactionModel) GetAllByUser(description string, userID int64, filters Filters) ([]*Transaction, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, deleted, version, user_id, category_id, description, amount, type
	FROM transactions
	WHERE (to_tsvector('simple', description) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND user_id = $2 AND deleted = false
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{description, userID, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	transactions := []*Transaction{}

	for rows.Next() {
		transaction := Transaction{
			User:     &User{},
			Category: &Category{User: &User{}},
		}
		err := rows.Scan(
			&totalRecords,
			&transaction.ID,
			&transaction.CreatedAt,
			&transaction.Deleted,
			&transaction.Version,
			&transaction.User.ID,
			&transaction.Category.ID,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Type,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		transactions = append(transactions, &transaction)
	}

	metaData := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return transactions, metaData, nil
}

func (m TransactionModel) GetByID(id int64, userID int64) (*Transaction, error) {
	query := `
	SELECT id, created_at, deleted, version, user_id, category_id, description, amount, type
	FROM transactions
	WHERE id = $1 AND user_id = $2 AND deleted = false
	`

	var tx Transaction
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&tx.ID,
		&tx.CreatedAt,
		&tx.Deleted,
		&tx.Version,
		&tx.User.ID,
		&tx.Category.ID,
		&tx.Description,
		&tx.Amount,
		&tx.Type,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &tx, nil
}

func (m TransactionModel) Insert(transaction *Transaction) error {
	query := `
	INSERT INTO transactions ( 
			user_id, 
			category_id, 
			description, 
			amount, 
			type
	)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id,created_at, version
	`

	args := []any{
		transaction.User.ID,
		transaction.Category.ID,
		transaction.Description,
		transaction.Amount,
		transaction.Type,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&transaction.ID,
		&transaction.CreatedAt,
		&transaction.Version,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m TransactionModel) Update(transaction *Transaction, userID int64) error {
	query := `
	UPDATE transactions
	SET user_id = $1, 
		category_id = $2, 
		description = $3, 
		amount = $4, 
		type = $5,
		version = version + 1
	WHERE 
		id = $6 
		AND user_id = $7 
		AND deleted = false 
		AND version = $8
	RETURNING version
	`

	args := []any{
		transaction.User.ID,
		transaction.Category.ID,
		transaction.Description,
		transaction.Amount,
		transaction.Type,
		transaction.ID,
		userID,
		transaction.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&transaction.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m TransactionModel) Delete(id int64, userID int64) error {
	query := `
	UPDATE transactions
	SET 
		deleted = true
	WHERE 
		id = $1 
		AND user_id = $2 
		AND deleted = false
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, userID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateTransaction(v *validator.Validator, transaction *Transaction) {
	v.Check(transaction.User != nil, "user", "must be provided")
	v.Check(transaction.Category != nil, "category", "must be provided")
	v.Check(transaction.Description != "", "description", "must be provided")
	v.Check(len(transaction.Description) <= 500, "description", "must not be more than 500 bytes long")
	v.Check(transaction.Amount > 0, "amount", "must be positive")
	v.Check(transaction.Type != 0, "type", "must be provided")
	v.Check(transaction.Amount != 0, "amount", "must be provided")
}
