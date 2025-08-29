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
	ID        int64          `json:"category_id"`
	CreatedAt *time.Time     `json:"-"`
	Name      *string        `json:"name"`
	Tipo      *TypeCategoria `json:"tipo"`
	Color     *string        `json:"color"`
	User      *UserDTO       `json:"-"`
	Version   *int           `json:"version"`
}

type CategoryModel struct {
	DB *sql.DB
}

func (c *Category) ToDTO() *CategoryDTO {
	var createdAt *time.Time
	if !c.CreatedAt.IsZero() {
		createdAt = &c.CreatedAt
	}
	var name *string
	if c.Name != "" {
		name = &c.Name
	}
	var tipo *TypeCategoria
	if c.Type != 0 {
		tipo = &c.Type
	}
	var color *string
	if c.Color != "" {
		color = &c.Color
	}
	var user *UserDTO
	if c.User != nil {
		user = c.User.ToDTO()
	}
	var version *int
	if c.Version != 0 {
		version = &c.Version
	}

	return &CategoryDTO{
		ID:        c.ID,
		CreatedAt: createdAt,
		Name:      name,
		Tipo:      tipo,
		Color:     color,
		User:      user,
		Version:   version,
	}
}

func (c *CategoryDTO) ToModel() *Category {
	var createdAt time.Time
	if c.CreatedAt != nil {
		createdAt = *c.CreatedAt
	}

	var name string
	if c.Name != nil {
		name = *c.Name
	}

	var tipo TypeCategoria
	if c.Tipo != nil {
		tipo = *c.Tipo
	}

	var color string
	if c.Color != nil {
		color = *c.Color
	}

	var user *User
	if c.User != nil {
		user = c.User.ToModel()
	}

	var version int
	if c.Version != nil {
		version = *c.Version
	}

	return &Category{
		ID:        c.ID,
		CreatedAt: createdAt,
		Name:      name,
		Type:      tipo,
		Color:     color,
		User:      user,
		Version:   version,
	}
}

func (c *CategoryDTO) ToDTOUpdateCategory(category *Category) *Category {
	if c.CreatedAt != nil {
		category.CreatedAt = *c.CreatedAt
	}

	if c.Name != nil {
		category.Name = *c.Name
	}

	if c.Tipo != nil {
		category.Type = *c.Tipo
	}

	if c.Color != nil {
		category.Color = *c.Color
	}

	if c.User != nil {
		category.User = c.User.ToModel()
	}

	if c.Version != nil {
		category.Version = *c.Version
	}

	return category
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

func (m CategoryModel) Delete(id int64, userID int64) error {
	query := `
	UPDATE from categories
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

func ValidateCategory(v *validator.Validator, category *Category) {
	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(category.Type.String() != "", "type", "must be provided")
	v.Check(category.Type.String() == "Unknown", "type", "invalid type")
	v.Check(category.Color != "", "color", "must be provided")
}
