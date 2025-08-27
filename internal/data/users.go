package data

import (
	"context"
	"database/sql"
	"errors"
	"meus_gastos/internal/validator"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Email     string
	Password  password
	Phone     string
	Activated bool
	Cod       int
	Version   int
	Deleted   bool
}

type UserDTO struct {
	ID        int64  `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Activated bool   `json:"activated"`
	Cod       int    `json:"cod"`
}

type UserSaveDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) ToDTO() *UserDTO {
	return &UserDTO{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Activated: u.Activated,
		Cod:       u.Cod,
	}
}

func (u *UserDTO) ToModel() *User {
	return &User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		Activated: u.Activated,
		Cod:       u.Cod,
	}
}

func (u *UserSaveDTO) ToModel() (*User, error) {
	user := &User{
		Name:  u.Name,
		Email: u.Email,
		Phone: u.Phone,
	}

	err := user.Password.Set(u.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m UserModel) GetByCodAndEmail(cod int, email string) (*User, error) {
	query := `
	SELECT id, created_at, name, phone, email,cod, password_hash, activated, version
	FROM users
	WHERE email = $1 AND deleted = false AND cod = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, email, cod).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Cod,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetByID(ID int64) (*User, error) {
	query := `
	SELECT id, created_at, name, phone, email, cod, password_hash, activated, version
	FROM users
	WHERE email = %1 AND deleted = false
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, ID).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Cod,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (name, email, phone,cod, password_hash, activated,deleted)
	VALUES ($1, $2, $3, $4, $5, $6,false)
	RETURNING id, created_at, version
	`

	args := []any{
		user.Name,
		user.Email,
		user.Phone,
		user.Cod,
		user.Password.hash,
		user.Activated,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, created_at, name, phone, email, cod, password_hash, activated, version
	FROM users
	WHERE email = $1 AND deleted = false
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Phone,
		&user.Email,
		&user.Cod,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) UpdateCodByEmail(user *User) error {
	query := `
	UDATE users SET
	cod = $1
	WHERE id = $1 AND version = $2
	RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, user.ID, user.Cod).Scan(
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

func (m UserModel) Update(user *User) error {
	query := `
	UPDATE users SET 
	name = $1, email = $2, cod = $3, phone = $4, password_hash = $5,
	activated = $6,version = version + 1
	WHERE id = $7 AND version = $8
	RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Cod,
		user.Phone,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Version,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m UserModel) Delete(user *User) error {
	query := `
	UPDATE users set
	deleted = true
	where id = $1 AND version = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, user.ID, user.Version)
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

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(user.Phone != "", "phone", "must be provided")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
