package db

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/denis-sukhoverkhov/calendar/app/domain/models"
	"github.com/denis-sukhoverkhov/calendar/app/interfaces/repositories"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userDbRepository struct {
	pool *pgxpool.Pool
	sq   sq.StatementBuilderType
}

func NewUserDbRepository(pgPool *pgxpool.Pool) repositories.UserRepository {
	return &userDbRepository{
		pool: pgPool,
		sq:   sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *userDbRepository) FindById(id int) (*models.User, error) {
	query := r.sq.Select("*").From("\"user\"").Where(sq.Eq{"id": id}, sq.Eq{"active": true})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("User.FindById QueryBuilder error %w", err)
	}

	user := &models.User{}
	err = r.pool.QueryRow(context.Background(), sql, args...).Scan(
		&user.Id, &user.FirstName, &user.LastName, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("User.FindById QueryRow error %w", err)
	}
	return user, nil
}

func (r *userDbRepository) FindAll() ([]*models.User, error) {
	query := r.sq.Select("*").From("\"user\"").Where(sq.Eq{"active": true})
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("User.FindAll QueryBuilder error %w", err)
	}

	users := make([]*models.User, 0)
	rows, err := r.pool.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, fmt.Errorf("User.FindAll rows errors %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		err := rows.Scan(
			&user.Id, &user.FirstName, &user.LastName, &user.Active, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("User.FindAll QueryRow error %w", err)
		}
		users = append(users, &user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("User.FindAll all errors %w", err)
	}

	return users, nil
}

func (r *userDbRepository) Store(user models.User) (*models.User, error) {
	query := r.sq.Insert("\"user\"").
		Columns("first_name", "last_name").
		Values(user.FirstName, user.LastName).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("User.FindAll QueryBuilder error %w", err)
	}

	newUser := &models.User{}
	err = r.pool.QueryRow(context.Background(), sql, args...).Scan(
		&newUser.Id, &newUser.FirstName, &newUser.LastName, &newUser.Active, &newUser.CreatedAt, &newUser.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("User.FindById execution error %w", err)
	}

	return newUser, nil
}

func (r *userDbRepository) Delete(id int) error {
	query := r.sq.Delete("\"user\"").Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("User.Delete QueryBuilder error %w", err)
	}

	_, err = r.pool.Exec(context.Background(), sql, args...)
	if err != nil {
		return fmt.Errorf("User.Delete execution error %w", err)
	}

	return nil
}
