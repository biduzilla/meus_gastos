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
	ID        int64
	CreatedAt time.Time
	Name      string
	Type      TypeCategoria
	Color     string
	User      *User
	Deleted   bool
	Version   int
}

type TypeCategoria int

const (
	Receita TypeCategoria = iota + 1
	Despesa
)

func (t TypeCategoria) String() string {
	switch t {
	case Receita:
		return "Receita"
	case Despesa:
		return "Despesa"
	default:
		return "Unknown"
	}
}

type CategoryDTO struct {
	ID        int64         `json:"category_id"`
	CreatedAt time.Time     `json:"created_at"`
	Name      string        `json:"name"`
	Tipo      TypeCategoria `json:"tipo"`
	Color     string        `json:"color"`
	User      *UserDTO      `json:"user"`
	Version   int           `json:"version"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (c *Category) ToDTO() *CategoryDTO {
	return &CategoryDTO{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		Name:      c.Name,
		Tipo:      c.Type,
		Color:     c.Color,
		User:      c.User.ToDTO(),
		Version:   c.Version,
	}
}

func (c *CategoryDTO) ToModel() *Category {
	return &Category{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		Name:      c.Name,
		Type:      c.Tipo,
		Color:     c.Color,
		User:      c.User.ToModel(),
		Version:   c.Version,
	}
}

func (m CategoryModel) Insert(category *Category) error {
	query := `
	INSERT INTO categories (name, type, color, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{
		category.Name,
		category.Type,
		category.Color,
		category.User.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m CategoryModel) GetByID(id int64, userID int64) (*Category, error) {
	query := `
	SELECT id, created_at, name, type, color, user_id,version
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
		&category.Type,
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
	SELECT count(*) OVER(), id, created_at, name, type, color, user_id,version
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
			&category.Type,
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

func (m CategoryModel) Update(category *Category, userID int64) error {
	query := `
	UPDATE categories
	SET 
		name = $1, 
		type = $2, 
		color = $3, 
		version = version + 1
	WHERE 
		id = $4 
		AND user_id = $5 
		AND deleted = false 
		AND version = $6
	RETURNING version
	`

	args := []any{
		category.Name,
		category.Type,
		category.Color,
		category.ID,
		userID,
		category.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&category.Version)

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

func (m CategoryModel) Delete(category *Category, userID int64) error {
	query := `
	UPDATE from categories
	SET
		deleted = true
	WHERE
		id = $1
		AND user_id = $2
		AND deleted = false
		AND version = $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, category.ID, userID, category.Version)

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

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(category.Type.String() != "", "type", "must be provided")
	v.Check(category.Type.String() == "Unknown", "type", "invalid type")
	v.Check(category.Color != "", "color", "must be provided")
}
