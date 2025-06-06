package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/donus-turkiye/backend/domain"
	"github.com/donus-turkiye/backend/pkg/config"
)

type PgRepository struct {
	db           *sql.DB
	SessionStore *PostgresStore
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

	return &PgRepository{
		db:           db,
		SessionStore: NewPostgresStore(db),
	}, nil
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

func (p *PgRepository) GetUserByTelNumber(ctx context.Context, telNumber string) (*domain.User, error) {
	var user domain.User

	err := p.db.QueryRowContext(
		ctx,
		`SELECT user_id, full_name, mail, password_hash, role_id, tel_no, adress, coordinate, wallet, total_recycle_count
		 FROM users WHERE tel_no = $1`,
		telNumber,
	).Scan(&user.Id, &user.FullName, &user.Email, &user.Password, &user.RoleId, &user.TelNumber, &user.Address, &user.Coordinate, &user.Wallet, &user.TotalRecycleCount)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by tel number: %w", err)
	}

	return &user, nil
}

func (p *PgRepository) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	err := p.db.QueryRowContext(
		ctx,
		`SELECT user_id, full_name, mail, role_id, tel_no, adress, coordinate, wallet, total_recycle_count
		 FROM users WHERE user_id = $1`,
		id,
	).Scan(&user.Id, &user.FullName, &user.Email, &user.RoleId, &user.TelNumber, &user.Address, &user.Coordinate, &user.Wallet, &user.TotalRecycleCount)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (p *PgRepository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	var categories []domain.Category

	rows, err := p.db.QueryContext(ctx, `SELECT category_id, waste_type, unit_type FROM category`)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var category domain.Category
		if err := rows.Scan(&category.CategoryId, &category.WasteType, &category.UnitType); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over categories: %w", err)
	}
	if len(categories) == 0 {
		return nil, fmt.Errorf("no categories found")
	}

	return categories, nil
}
