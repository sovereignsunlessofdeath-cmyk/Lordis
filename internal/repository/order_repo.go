package repository

import (
    "database/sql"

    "lordis/internal/models"
)

type OrderRepo struct{
    DB *sql.DB
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
    return &OrderRepo{DB: db}
}

// Create adds a new order to the orders table (ensure you have migrated it).
func (r *OrderRepo) Create(username, itemID string, qty int) (int, error) {
    var id int
    // This query assumes an `orders` table exists. If you don't use Postgres for orders,
    // keep using the JSON store or adapt accordingly.
    err := r.DB.QueryRow(`INSERT INTO orders (username, item_id, quantity, status, created_at) VALUES ($1,$2,$3,$4,now()) RETURNING id`, username, itemID, qty, "pending").Scan(&id)
    return id, err
}

// ListByUser retrieves orders for a specific username.
func (r *OrderRepo) ListByUser(username string) ([]models.Order, error) {
    rows, err := r.DB.Query(`SELECT id, username, item_id, quantity, status, created_at FROM orders WHERE username=$1 ORDER BY created_at DESC`, username)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var out []models.Order
    for rows.Next() { 
        var o models.Order
        if err := rows.Scan(&o.ID, &o.Username, &o.ItemID, &o.Quantity, &o.Status, &o.CreatedAt); err != nil {
            return nil, err
        }
        out = append(out, o)
    }
    return out, nil
}
