package database

import (
	"context"
	"list-application-v2/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

func DeleteUserByID(conn *pgx.Conn, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := conn.Exec(context.Background(), query, id)
	return err

}

func EditUser(conn *pgx.Conn, p models.User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, password = $3, picture_url = $4
		WHERE id = $5
	`
	_, err := conn.Exec(context.Background(), query,
		p.Username,
		p.Email,
		p.Password,
		p.PictureUrl,
		p.ID,
	)
	return err
}

func SaveNewUser(conn *pgx.Conn, p models.User) error {
	query := `
		INSERT INTO users (username, email, password, created_at, picture_url)
		VALUES ($1, $2, $3, $4, $5)
	`
	p.CreatedAt = time.Now()
	_, err := conn.Exec(context.Background(), query,
		p.Username,
		p.Email,
		p.Password,
		p.CreatedAt,
		"",
	)
	return err
}

func GetAllUsers(conn *pgx.Conn) ([]models.User, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, username, email, password, created_at FROM users")
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

func GetUserByUsername(conn *pgx.Conn, username string) (models.User, error) {
	var u models.User
	query := `SELECT id, username, password FROM users WHERE username = $1`

	err := conn.QueryRow(context.Background(), query, username).Scan(&u.ID, &u.Username, &u.Password)
	return u, err
}

func GetUserByID(conn *pgx.Conn, userID int) (models.User, error) {
	var u models.User
	query := `SELECT id, username, picture_url, email, created_at FROM users WHERE id = $1`

	err := conn.QueryRow(context.Background(), query, userID).Scan(&u.ID, &u.Username, &u.PictureUrl, &u.Email, &u.CreatedAt)
	return u, err
}
