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
}

type TransactionDTO struct {
	ID          *int64       `json:"transaction_id"`
	Version     *int         `json:"version"`
	User        *UserDTO     `json:"user"`
	Category    *CategoryDTO `json:"category"`
	Description *string      `json:"description"`
	Amount      *float64     `json:"amount"`
	CreatedAt   *time.Time   `json:"created_at"`
}

func (t *Transaction) ToDTO() *TransactionDTO {
	dto := &TransactionDTO{}

	if t.ID != 0 {
		dto.ID = &t.ID
	}

	if t.Version != 0 {
		dto.Version = &t.Version
	}

	if t.User != nil {
		dto.User = t.User.ToDTO()
	}

	if t.Category != nil {
		dto.Category = t.Category.ToDTO()
	}

	if t.Description != "" {
		dto.Description = &t.Description
	}

	if t.Amount != 0 {
		dto.Amount = &t.Amount
	}

	dto.CreatedAt = &t.CreatedAt

	return dto
}

func (t *TransactionDTO) ToModel() *Transaction {
	transaction := &Transaction{}

	if t.ID != nil {
		transaction.ID = *t.ID
	}
	if t.Version != nil {
		transaction.Version = *t.Version
	}
	if t.User != nil {
		transaction.User = t.User.ToModel()
	}
	if t.Category != nil {
		transaction.Category = t.Category.ToModel()
	}
	if t.Description != nil {
		transaction.Description = *t.Description
	}
	if t.Amount != nil {
		transaction.Amount = *t.Amount
	}

	return transaction
}

func (t *TransactionDTO) ToDTOUpdateTransaction(transaction *Transaction) {
	if t.ID != nil {
		transaction.ID = *t.ID
	}

	if t.Version != nil {
		transaction.Version = *t.Version
	}

	if t.User != nil {
		transaction.User = t.User.ToModel()
	}

	if t.Category != nil {
		transaction.Category = t.Category.ToModel()
	}

	if t.Description != nil {
		transaction.Description = *t.Description
	}

	if t.Amount != nil {
		transaction.Amount = *t.Amount
	}
}

func (m TransactionModel) GetAllByUserAndCategory(description string, userID int64, categoryID int64, startDate, endDate time.Time, filters Filters) ([]*Transaction, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, created_at, deleted, version, user_id, category_id, description, amount
	FROM transactions
	WHERE (to_tsvector('simple', description) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND user_id = $2 AND deleted = false AND category_id = $3
	AND ($4 IS NULL OR (t.created_at >= $4))
	AND ($5 IS NULL OR (t.created_at <= $5))
	ORDER BY %s %s, id ASC
	LIMIT $6 OFFSET $7
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{description, userID, categoryID, startDate, endDate, filters.limit(), filters.offset()}
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
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		transactions = append(transactions, &transaction)
	}

	metaData := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return transactions, metaData, nil
}

func (m TransactionModel) GetAllByUser(description string, userID int64, startDate, endDate time.Time, categoryType TypeCategoria, filters Filters) ([]*Transaction, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), 
		t.id, 
		t.created_at, 
		t.deleted, 
		t.version, 
		t.user_id, 
		t.category_id, 
		t.description, 
		t.amount
	FROM transactions t
	INNER JOIN categories c ON c.id = t.category_id
	WHERE (to_tsvector('simple', t.description) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND t.user_id = $2 
	AND t.deleted = false
	AND ($3 IS NULL OR (t.created_at >= $3))
	AND ($4 IS NULL OR (t.created_at <= $4))
	AND ($5 IS NULL OR c.type = $5)
	ORDER BY %s %s, t.id ASC
	LIMIT $6 OFFSET $7
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		description,
		userID,
		startDate,
		endDate,
		categoryType,
		filters.limit(),
		filters.offset(),
	}

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
	SELECT id, created_at, deleted, version, user_id, category_id, description, amount
	FROM transactions
	WHERE id = $1 AND user_id = $2 AND deleted = false
	`

	var tx Transaction
	tx.User = &User{}
	tx.Category = &Category{}

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
	)
	VALUES ($1, $2, $3, $4)
	RETURNING id,created_at, version
	`

	args := []any{
		transaction.User.ID,
		transaction.Category.ID,
		transaction.Description,
		transaction.Amount,
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
		version = version + 1
	WHERE 
		id = $5
		AND user_id = $6
		AND deleted = false 
		AND version = $7
	RETURNING version
	`

	args := []any{
		transaction.User.ID,
		transaction.Category.ID,
		transaction.Description,
		transaction.Amount,
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
	v.Check(transaction.Amount != 0, "amount", "must be provided")
}
