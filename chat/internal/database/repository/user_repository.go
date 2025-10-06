package repository

import (
	"context"
	"database/sql"
	"playgoround/chat/internal/database/entity"
	"playgoround/chat/internal/database/utils"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	query, args, err := sq.Insert("users").
		Columns("username", "email", "created_at", "updated_at").
		Values(user.Username, user.Email, user.CreatedAt, user.UpdatedAt).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	user.ID = id
	return user, nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*entity.User, error) {
	query, args, err := sq.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var user entity.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query, args, err := sq.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var user entity.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	query, args, err := sq.Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("updated_at", user.UpdatedAt).
		Where(sq.Eq{"id": user.ID}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}
	if rows == 0 {
		return nil, utils.NewNotFoundError(query, args, nil)
	}

	return user, nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(ctx context.Context, id int64) (*entity.User, error) {
	query, args, err := sq.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var user entity.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	deleteQuery, deleteArgs, err := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	_, err = r.db.ExecContext(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	return &user, nil
}

// List retrieves all users with pagination
func (r *UserRepository) List(ctx context.Context, limit, offset uint64) ([]*entity.User, error) {
	query, args, err := sq.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		OrderBy("id DESC").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, utils.NewSystemError(query, args, err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return users, nil
}
