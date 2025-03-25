package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/donus-turkiye/backend/domain"
	"github.com/donus-turkiye/backend/pkg/config"
)

type PgRepository struct {
	db *sql.DB
}

func NewPgRepository(cfg *config.AppConfig) (*PgRepository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &PgRepository{db: db}, nil
}

func (p *PgRepository) CreateUser(ctx context.Context, user *domain.User) (int, error) {
	var userID int

	err := p.db.QueryRowContext(
		ctx,
		`INSERT INTO users (full_name, mail, password_hash, role_id, tel_no, adress, coordinate) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING user_id`,
		user.FullName, user.Email, user.Password, user.RoleId, user.TelNumber, user.Address, user.Coordinate,
	).Scan(&userID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, nil
}

func (p *PgRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := p.db.QueryRowContext(
		ctx,
		`SELECT user_id, full_name, mail, password_hash, role_id, tel_no, adress, coordinate, wallet, total_recycle_count
		 FROM users WHERE mail = $1`,
		email,
	).Scan(&user.Id, &user.FullName, &user.Email, &user.Password, &user.RoleId, &user.TelNumber, &user.Address, &user.Coordinate, &user.Wallet, &user.TotalRecycleCount)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}
