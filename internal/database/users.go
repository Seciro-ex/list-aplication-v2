package database

import (
	"context"
	"list-application-v2/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func DeleteUserByID(ctx context.Context, db *pgxpool.Pool, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := db.Exec(ctx, query, id)
	return err

}

func EditUser(ctx context.Context, db *pgxpool.Pool, p models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, picture_url = $4
		WHERE id = $5
	`
	_, err := db.Exec(ctx, query,
		p.Username,
		p.Email,
		p.Password,
		p.PictureUrl,
		p.ID,
	)
	return err
}

func SaveNewUser(ctx context.Context, db *pgxpool.Pool, p models.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, picture_url)
		VALUES ($1, $2, $3, $4, $5)
	`
	p.CreatedAt = time.Now()
	_, err := db.Exec(ctx, query,
		p.Username,
		p.Email,
		p.Password,
		p.CreatedAt,
		"",
	)
	return err
}

func GetAllUsers(ctx context.Context, db *pgxpool.Pool) ([]models.User, error) {
	rows, err := db.Query(ctx, "SELECT id, username, email, password, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.User
	for rows.Next() {
		var p models.User
		err := rows.Scan(&p.ID, &p.Username, &p.Email, &p.Password, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}

func GetUserByUsername(ctx context.Context, db *pgxpool.Pool, username string) (models.User, error) {
	var u models.User
	query := `SELECT id, username, password FROM users WHERE username = $1`

	err := db.QueryRow(ctx, query, username).Scan(&u.ID, &u.Username, &u.Password)
	return u, err
}

func GetUserByID(ctx context.Context, db *pgxpool.Pool, userID int) (models.User, error) {
	var u models.User
	query := `SELECT id, username, picture_url, email, password, created_at FROM users WHERE id = $1`

	err := db.QueryRow(ctx, query, userID).Scan(
		&u.ID, &u.Username, &u.PictureUrl, &u.Email, &u.Password, &u.CreatedAt,
	)
	return u, err
}
