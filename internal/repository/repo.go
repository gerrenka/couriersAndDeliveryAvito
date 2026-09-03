package repository

import (
	"avito/internal/model"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCourierRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCourierRepository(db *pgxpool.Pool) *PostgresCourierRepository {
	return &PostgresCourierRepository{db: db}
}

func (r *PostgresCourierRepository) CreateCourier(ctx context.Context, name, phone, status, transport_type string) (model.Courier, error) {
	var cour model.Courier
	query := `INSERT INTO couriers (name, phone, status, transport_type) VALUES($1, $2, $3, $4) RETURNING id, name, phone, status, created_at, updated_at, transport_type`
	err := r.db.QueryRow(ctx, query, name, phone, status, transport_type).Scan(&cour.ID, &cour.Name, &cour.Phone, &cour.Status, &cour.CreatedAt, &cour.UpdatedAt, &cour.TransportType)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return model.Courier{}, model.ErrPhoneExists
	}
	if err != nil {
		return model.Courier{}, err
	}
	return cour, nil

}

func (r *PostgresCourierRepository) UpdateCourier(ctx context.Context, id int, name, phone, status, transport_type string) (model.Courier, error) {
	query := `UPDATE couriers SET name = $1, phone = $2, status = $3, updated_at = now(), transport_type = $4 WHERE id = $5 RETURNING id, name, phone, status, created_at, updated_at, transport_type`
	var c model.Courier
	err := r.db.QueryRow(ctx, query, name, phone, status, transport_type, id).Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt, &c.TransportType)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Courier{}, model.ErrCourierNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return model.Courier{}, model.ErrPhoneExists
	}
	if err != nil {
		return model.Courier{}, err
	}
	return c, nil

}

func (r *PostgresCourierRepository) GetById(ctx context.Context, id int) (model.Courier, error) {
	query := `SELECT id, name, phone, status, created_at, updated_at, transport_type FROM couriers WHERE id = $1`
	var cour model.Courier
	err := r.db.QueryRow(ctx, query, id).Scan(&cour.ID, &cour.Name, &cour.Phone, &cour.Status, &cour.CreatedAt, &cour.UpdatedAt, &cour.TransportType)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Courier{}, model.ErrCourierNotFound
	}
	if err != nil {
		return model.Courier{}, err
	}
	return cour, nil
}

func (r *PostgresCourierRepository) GetAllCourier(ctx context.Context) ([]model.Courier, error) {
	query := `SELECT id, name, phone, status, created_at, updated_at, transport_type FROM couriers`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	couriers := []model.Courier{}
	for rows.Next() {
		var c model.Courier
		if err := rows.Scan(&c.ID, &c.Name, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt, &c.TransportType); err != nil {
			return nil, err
		}
		couriers = append(couriers, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return couriers, nil
}

func (r *PostgresCourierRepository) DeleteCourier(ctx context.Context, id int) error {
	cmd, err := r.db.Exec(ctx, "DELETE FROM couriers WHERE id = $1", id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return model.ErrCourierNotFound
	}
	return nil
}

// Блок В — равномерное распределение. Переписать FindAvailable: LEFT JOIN к подсчёту завершённых доставок,
// ORDER BY count ASC LIMIT 1, плюс блокировка строки.

func (r *PostgresCourierRepository) FindAvailable(ctx context.Context, tx pgx.Tx) (model.Courier, error) {
	var cour model.Courier
	query :=
		`SELECT c.id, c.name, c.phone, c.status, c.created_at, c.updated_at, c.transport_type  
	FROM couriers c 
	LEFT JOIN (SELECT courier_id, COUNT(*) AS cnt
		FROM delivery
		WHERE status = 'completed'
		GROUP BY courier_id) d ON c.id = d.courier_id 
	WHERE c.status = 'available' ORDER BY COALESCE(d.cnt, 0), c.id
	LIMIT 1
	FOR UPDATE OF c SKIP LOCKED`
	err := tx.QueryRow(ctx, query).Scan(&cour.ID, &cour.Name, &cour.Phone, &cour.Status, &cour.CreatedAt, &cour.UpdatedAt, &cour.TransportType)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Courier{}, model.ErrCourierNotAvailable
	}
	if err != nil {
		return model.Courier{}, err
	}
	return cour, nil
}

func (r *PostgresCourierRepository) UpdateStatus(ctx context.Context, tx pgx.Tx, id int, status string) error {
	cmd, err := tx.Exec(ctx, `UPDATE couriers SET status = $1 WHERE id = $2`, status, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return model.ErrCourierNotFound
	}
	return nil
}
