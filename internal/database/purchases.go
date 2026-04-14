package database

import (
	"context"
	"list-application-v2/internal/models"

	"github.com/jackc/pgx/v5"
)

func DeletePurchaseByID(conn *pgx.Conn, id int) error {
	query := `
		DELETE FROM purchases
		WHERE id = $1
	`
	_, err := conn.Exec(context.Background(), query, id)
	return err
}

func EditPurchase(conn *pgx.Conn, p models.Purchase) error {
	query := `
		UPDATE purchases
		SET image_url = $1, cost = $2, item_url = $3, item_name = $4
		WHERE id = $5
	`
	_, err := conn.Exec(context.Background(), query,
		p.ImageURL,
		p.Cost,
		p.ItemURL,
		p.ItemName,
		p.ID,
	)
	return err
}

func SavePurchase(conn *pgx.Conn, p models.Purchase) error {
	query := `
		INSERT INTO purchases (image_url, cost, item_url, item_name, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := conn.Exec(context.Background(), query,
		p.ImageURL,
		p.Cost,
		p.ItemURL,
		p.ItemName,
		p.UserID,
	)
	return err
}

func GetAllPurchases(conn *pgx.Conn) ([]models.Purchase, error) {
	rows, err := conn.Query(context.Background(), "SELECT id, item_name, cost, item_url FROM purchases")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Purchase
	for rows.Next() {
		var p models.Purchase
		err := rows.Scan(&p.ID, &p.ItemName, &p.Cost, &p.ItemURL)
		if err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, nil
}
