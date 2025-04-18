package postgres

import (
	"database/sql"
	"fmt"
	"time"
)

type PostgresStore struct {
	db        *sql.DB
	tableName string
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{
		db:        db,
		tableName: "sessions",
	}
}

// Get gets the value for the given key.
// `nil, nil` is returned when the key does not exist
func (p *PostgresStore) Get(id string) ([]byte, error) {
	var data []byte
	err := p.db.QueryRow(
		"SELECT data FROM sessions WHERE id = $1 AND expiry > NOW()",
		id,
	).Scan(&data)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return data, nil
}

// Set stores the given value for the given key along
// with an expiration value, 0 means no expiration.
// Empty key or value will be ignored without an error.
func (p *PostgresStore) Set(id string, data []byte, expiry time.Duration) error {
	expiryTime := time.Now().Add(expiry)

	_, err := p.db.Exec(
		`INSERT INTO sessions (id, data, expiry) 
         VALUES ($1, $2, $3) 
         ON CONFLICT (id) DO UPDATE 
         SET data = EXCLUDED.data, 
             expiry = EXCLUDED.expiry`,
		id, data, expiryTime,
	)

	if err != nil {
		return fmt.Errorf("failed to set session: %w", err)
	}

	return nil
}

// Delete deletes the value for the given key.
// It returns no error if the storage does not contain the key,
func (p *PostgresStore) Delete(id string) error {
	_, err := p.db.Exec("DELETE FROM sessions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// Reset resets the storage and delete all keys.
func (p *PostgresStore) Reset() error {
	_, err := p.db.Exec("DELETE FROM sessions")
	if err != nil {
		return fmt.Errorf("failed to reset sessions: %w", err)
	}

	return nil
}

// Close closes the storage and will stop any running garbage
// collectors and open connections.
func (p *PostgresStore) Close() error {
	if p.db != nil {
		if err := p.db.Close(); err != nil {
			return fmt.Errorf("failed to close db: %w", err)
		}
	}

	return nil
}

// GC deletes expired sessions from the database.
func (p *PostgresStore) GC() error {
	_, err := p.db.Exec("DELETE FROM sessions WHERE expiry <= NOW()")
	if err != nil {
		return fmt.Errorf("failed to garbage collect sessions: %w", err)
	}

	return nil
}
