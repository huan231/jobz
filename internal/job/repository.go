package job

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	Add(ctx context.Context, j Job) (bool, error)
	List(ctx context.Context) ([]Job, error)
	CompleteByID(ctx context.Context, id string, s Status, t time.Time) (bool, error)
	DeleteByCronJobID(ctx context.Context, cronJobID string) error
	FindByID(ctx context.Context, id string) (*Job, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) Add(ctx context.Context, j Job) (bool, error) {
	stmt, err := r.db.PrepareContext(
		ctx,
		`INSERT OR IGNORE INTO job (id, cron_job_id, status, created_at, updated_at, completed_at)
    	VALUES (:id, :cronJobID, :status, :createdAt, :updatedAt, :completedAt)`,
	)

	if err != nil {
		return false, err
	}

	var completedAt *int64

	if j.CompletedAt != nil {
		u := j.CompletedAt.UTC().UnixMilli()

		completedAt = &u
	}

	res, err := stmt.ExecContext(
		ctx,
		sql.Named("id", j.ID),
		sql.Named("cronJobID", j.CronJobID),
		sql.Named("status", j.Status),
		sql.Named("createdAt", j.CreatedAt.UTC().UnixMilli()),
		sql.Named("updatedAt", j.UpdatedAt.UTC().UnixMilli()),
		sql.Named("completedAt", completedAt),
	)

	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	return n == 1, nil
}

func (r *repository) List(ctx context.Context) ([]Job, error) {
	stmt, err := r.db.PrepareContext(
		ctx,
		`SELECT id, cron_job_id, status, created_at, updated_at, completed_at FROM job ORDER BY updated_at ASC`,
	)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []Job

	for rows.Next() {
		job, err := scan(rows)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, *job)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *repository) CompleteByID(ctx context.Context, id string, s Status, t time.Time) (bool, error) {
	stmt, err := r.db.PrepareContext(
		ctx,
		`UPDATE job SET status = :status, updated_at = :t, completed_at = :t WHERE id = :id AND completed_at IS NULL`,
	)

	if err != nil {
		return false, err
	}

	res, err := stmt.ExecContext(
		ctx,
		sql.Named("id", id),
		sql.Named("status", s),
		sql.Named("t", t.UTC().UnixMilli()),
	)

	if err != nil {
		return false, err
	}

	n, err := res.RowsAffected()

	if err != nil {
		return false, err
	}

	return n == 1, nil
}

func (r *repository) DeleteByCronJobID(ctx context.Context, cronJobID string) error {
	stmt, err := r.db.PrepareContext(ctx, `DELETE FROM job WHERE cron_job_id = :cronJobID`)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, sql.Named("cronJobID", cronJobID))

	return err
}

func (r *repository) FindByID(ctx context.Context, id string) (*Job, error) {
	stmt, err := r.db.PrepareContext(ctx, `SELECT id, cron_job_id, status, created_at, updated_at, completed_at FROM job WHERE id = :id`)

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, sql.Named("id", id))

	return scan(row)
}

type scannable interface {
	Scan(dest ...any) error
}

func scan(s scannable) (*Job, error) {
	job := Job{}

	var createdAt, updatedAt int64
	var completedAt sql.NullInt64

	if err := s.Scan(&job.ID, &job.CronJobID, &job.Status, &createdAt, &updatedAt, &completedAt); err != nil {
		return nil, err
	}

	job.CreatedAt = time.UnixMilli(createdAt)
	job.UpdatedAt = time.UnixMilli(updatedAt)

	if completedAt.Valid {
		u := time.UnixMilli(completedAt.Int64)

		job.CompletedAt = &u
	}

	return &job, nil
}
