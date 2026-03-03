package repository

import (
	"database/sql"
	"github.com/hxseqwe/korochki-est/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (login, password_hash, full_name, phone, email, is_admin) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`

	return r.db.QueryRow(query, user.Login, string(hashedPassword), user.FullName,
		user.Phone, user.Email, user.IsAdmin).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) FindByLogin(login string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, login, password_hash, full_name, phone, email, is_admin, created_at, updated_at 
              FROM users WHERE login = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash,
		&user.FullName, &user.Phone, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id, login, full_name, phone, email, is_admin, created_at, updated_at 
              FROM users WHERE id = $1 AND deleted_at IS NULL`

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Login, &user.FullName,
		&user.Phone, &user.Email, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) ValidatePassword(user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

func (r *UserRepository) IsLoginExists(login string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 AND deleted_at IS NULL)`
	err := r.db.QueryRow(query, login).Scan(&exists)
	return exists, err
}
