package db

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sensor-manager-service/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(connString string) (*Repository, error) {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Repository{pool: pool}, nil
}

func (r *Repository) Close() {
	r.pool.Close()
}

func (r *Repository) ListSensors(ctx context.Context, limit, offset int, houseID, typeID *int64, status *string) ([]model.Sensor, int, error) {
	where := []string{"1=1"}
	args := []any{}
	argPos := 1

	if houseID != nil {
		where = append(where, fmt.Sprintf("house_id = $%d", argPos))
		args = append(args, *houseID)
		argPos++
	}
	if typeID != nil {
		where = append(where, fmt.Sprintf("type_id = $%d", argPos))
		args = append(args, *typeID)
		argPos++
	}
	if status != nil {
		where = append(where, fmt.Sprintf("status = $%d", argPos))
		args = append(args, *status)
		argPos++
	}

	countQuery := "SELECT COUNT(*) FROM sensors WHERE " + strings.Join(where, " AND ")
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, house_id, name, type_id, location, serial_number, value, unit, status, last_updated, created_at
		FROM sensors
		WHERE ` + strings.Join(where, " AND ") + fmt.Sprintf(" ORDER BY id LIMIT $%d OFFSET $%d", argPos, argPos+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]model.Sensor, 0)
	for rows.Next() {
		var s model.Sensor
		if err := rows.Scan(
			&s.ID,
			&s.HouseID,
			&s.Name,
			&s.TypeID,
			&s.Location,
			&s.SerialNumber,
			&s.Value,
			&s.Unit,
			&s.Status,
			&s.LastUpdated,
			&s.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		items = append(items, s)
	}

	return items, total, rows.Err()
}

func (r *Repository) GetSensor(ctx context.Context, id int64) (*model.Sensor, error) {
	query := `
		SELECT id, house_id, name, type_id, location, serial_number, value, unit, status, last_updated, created_at
		FROM sensors
		WHERE id = $1
	`
	var s model.Sensor
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID,
		&s.HouseID,
		&s.Name,
		&s.TypeID,
		&s.Location,
		&s.SerialNumber,
		&s.Value,
		&s.Unit,
		&s.Status,
		&s.LastUpdated,
		&s.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) CreateSensor(ctx context.Context, in model.SensorCreate) (*model.Sensor, error) {
	query := `
		INSERT INTO sensors (house_id, name, type_id, type, location, serial_number, value, unit, status, last_updated, created_at)
		VALUES ($1, $2, $3, (SELECT name FROM sensor_types WHERE id = $3), $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, house_id, name, type_id, location, serial_number, value, unit, status, last_updated, created_at
	`
	var s model.Sensor
	if err := r.pool.QueryRow(ctx, query,
		in.HouseID,
		in.Name,
		in.TypeID,
		in.Location,
		in.SerialNumber,
		in.Value,
		in.Unit,
		in.Status,
	).Scan(
		&s.ID,
		&s.HouseID,
		&s.Name,
		&s.TypeID,
		&s.Location,
		&s.SerialNumber,
		&s.Value,
		&s.Unit,
		&s.Status,
		&s.LastUpdated,
		&s.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) UpdateSensor(ctx context.Context, id int64, in model.SensorUpdate) (*model.Sensor, error) {
	setParts := []string{"last_updated = NOW()"}
	args := []any{}
	idx := 1

	if in.Name != nil {
		setParts = append(setParts, fmt.Sprintf("name = $%d", idx))
		args = append(args, *in.Name)
		idx++
	}
	if in.TypeID != nil {
		setParts = append(setParts, fmt.Sprintf("type_id = $%d", idx))
		args = append(args, *in.TypeID)
		setParts = append(setParts, fmt.Sprintf("type = (SELECT name FROM sensor_types WHERE id = $%d)", idx))
		idx++
	}
	if in.Location != nil {
		setParts = append(setParts, fmt.Sprintf("location = $%d", idx))
		args = append(args, *in.Location)
		idx++
	}
	if in.SerialNumber != nil {
		setParts = append(setParts, fmt.Sprintf("serial_number = $%d", idx))
		args = append(args, *in.SerialNumber)
		idx++
	}
	if in.Value != nil {
		setParts = append(setParts, fmt.Sprintf("value = $%d", idx))
		args = append(args, *in.Value)
		idx++
	}
	if in.Unit != nil {
		setParts = append(setParts, fmt.Sprintf("unit = $%d", idx))
		args = append(args, *in.Unit)
		idx++
	}
	if in.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", idx))
		args = append(args, *in.Status)
		idx++
	}

	query := `
		UPDATE sensors
		SET ` + strings.Join(setParts, ", ") + fmt.Sprintf(`
		WHERE id = $%d
		RETURNING id, house_id, name, type_id, location, serial_number, value, unit, status, last_updated, created_at
	`, idx)
	args = append(args, id)

	var s model.Sensor
	if err := r.pool.QueryRow(ctx, query, args...).Scan(
		&s.ID,
		&s.HouseID,
		&s.Name,
		&s.TypeID,
		&s.Location,
		&s.SerialNumber,
		&s.Value,
		&s.Unit,
		&s.Status,
		&s.LastUpdated,
		&s.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) DeleteSensor(ctx context.Context, id int64) error {
	ct, err := r.pool.Exec(ctx, "DELETE FROM sensors WHERE id = $1", id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ApplyTelemetry(ctx context.Context, id int64, value float64, unit, status string, ts time.Time) error {
	query := `
		UPDATE sensors
		SET value = $1, unit = $2, status = $3, last_updated = $4
		WHERE id = $5
	`
	_, err := r.pool.Exec(ctx, query, value, unit, status, ts, id)
	return err
}

var ErrNotFound = fmt.Errorf("not found")
