package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"meus_gastos/internal/validator"
	"time"
)

type Category struct {
	ID          int64
	CreatedAt   time.Time
	Name        string
	Description string
	Color       string
	User        *User
	Deleted     bool
	Version     int
}

type CategoryDTO struct {
	ID          int64     `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	User        *UserDTO  `json:"user"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (c *Category) ToDTO() *CategoryDTO {
	return &CategoryDTO{
		ID:          c.ID,
		CreatedAt:   c.CreatedAt,
		Name:        c.Name,
		Description: c.Description,
		Color:       c.Color,
		User:        c.User.ToDTO(),
	}
}

func (c *CategoryDTO) ToModel() *Category {
	return &Category{
		ID:          c.ID,
		CreatedAt:   c.CreatedAt,
		Name:        c.Name,
		Description: c.Description,
		Color:       c.Color,
		User:        c.User.ToModel(),
	}
}

func (m CategoryModel) Insert(category *Category) error {
	query := `
	INSERT INTO categories (name, description, color, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`

	args := []any{
		category.Name,
		category.Description,
		category.Color,
		category.User.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m CategoryModel) GetByID(id int64, userID int64) (*Category, error) {
	query := `
	SELECT category_id, created_at, name, description, color, user_id,version
	FROM categories
	WHERE id = $1 AND user_id = $2 AND deleted = false
	`

	category := Category{
		User: &User{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Name,
		&category.Description,
		&category.Color,
		&category.User.ID,
		&category.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &category, nil
}

func (m CategoryModel) GetAll(name string, userID int64, filters Filters) ([]*Category, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), category_id, created_at, name, description, color, user_id,version
	FROM categories
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
	AND user_id = $2 AND deleted = false
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{name, userID, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	categories := []*Category{}

	for rows.Next() {
		category := Category{
			User: &User{},
		}

		err := rows.Scan(
			&totalRecords,
			&category.ID,
			&category.CreatedAt,
			&category.Name,
			&category.Description,
			&category.Color,
			&category.User.ID,
			&category.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metaData := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return categories, metaData, nil
}

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(category.Description != "", "description", "must be provided")
	v.Check(len(category.Description) <= 500, "description", "must not be more than 500 bytes long")
	v.Check(category.Color != "", "color", "must be provided")
}
