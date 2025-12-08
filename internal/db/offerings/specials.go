package offerings

import (
	"context"
	"database/sql"
)

// Special represents a travel special offering stored in the database.
type Special struct {
	ID       string
	Name     string
	Price    float64
	Currency string
}

// GetActiveSpecials return all active specials from the database.
func GetActiveSpecials(ctx context.Context, db *sql.DB) ([]Special, error) {
	const q = `SELECT id, name, price, currency FROM specials WHERE active = 1 ORDER BY id`

	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var out []Special

	for rows.Next() {
		var s Special
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.Currency); err != nil {
			return nil, err
		}
		out = append(out, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
