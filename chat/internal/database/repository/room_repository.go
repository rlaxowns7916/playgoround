package repository

import (
	"context"
	"database/sql"
	"playgoround/chat/internal/database/entity"
	"playgoround/chat/internal/database/utils"

	sq "github.com/Masterminds/squirrel"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

// Create inserts a new room
func (r *RoomRepository) Create(ctx context.Context, room *entity.Room) (*entity.Room, error) {
	query, args, err := sq.Insert("rooms").
		Columns("name", "created_at", "updated_at").
		Values(room.Name, room.CreatedAt, room.UpdatedAt).
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

	room.ID = id
	return room, nil
}

// FindByID retrieves a room by ID
func (r *RoomRepository) FindByID(ctx context.Context, id int64) (*entity.Room, error) {
	query, args, err := sq.Select("id", "name", "created_at", "updated_at").
		From("rooms").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var room entity.Room
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&room.ID,
		&room.Name,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return &room, nil
}

// Update updates a room
func (r *RoomRepository) Update(ctx context.Context, room *entity.Room) (*entity.Room, error) {
	query, args, err := sq.Update("rooms").
		Set("name", room.Name).
		Set("updated_at", room.UpdatedAt).
		Where(sq.Eq{"id": room.ID}).
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

	return room, nil
}

// Delete deletes a room by ID
func (r *RoomRepository) Delete(ctx context.Context, id int64) (*entity.Room, error) {
	query, args, err := sq.Select("id", "name", "created_at", "updated_at").
		From("rooms").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var room entity.Room
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&room.ID,
		&room.Name,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	deleteQuery, deleteArgs, err := sq.Delete("rooms").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	_, err = r.db.ExecContext(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	return &room, nil
}

// List retrieves all rooms with pagination
func (r *RoomRepository) List(ctx context.Context, limit, offset uint64) ([]*entity.Room, error) {
	query, args, err := sq.Select("id", "name", "created_at", "updated_at").
		From("rooms").
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

	var rooms []*entity.Room
	for rows.Next() {
		var room entity.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, utils.NewSystemError(query, args, err)
		}
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return rooms, nil
}

// AddUserToRoom adds a user to a room (creates UserRoom relationship)
func (r *RoomRepository) AddUserToRoom(ctx context.Context, userRoom *entity.UserRoom) (*entity.UserRoom, error) {
	query, args, err := sq.Insert("user_rooms").
		Columns("user_id", "room_id", "joined_at").
		Values(userRoom.UserID, userRoom.RoomID, userRoom.JoinedAt).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return userRoom, nil
}

// RemoveUserFromRoom removes a user from a room
func (r *RoomRepository) RemoveUserFromRoom(ctx context.Context, userID, roomID int64) (*entity.UserRoom, error) {
	query, args, err := sq.Select("user_id", "room_id", "joined_at").
		From("user_rooms").
		Where(sq.Eq{"user_id": userID, "room_id": roomID}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var userRoom entity.UserRoom
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&userRoom.UserID,
		&userRoom.RoomID,
		&userRoom.JoinedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	deleteQuery, deleteArgs, err := sq.Delete("user_rooms").
		Where(sq.Eq{"user_id": userID, "room_id": roomID}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	_, err = r.db.ExecContext(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	return &userRoom, nil
}

// GetRoomUsers retrieves all users in a room
func (r *RoomRepository) GetRoomUsers(ctx context.Context, roomID int64) ([]*entity.User, error) {
	query, args, err := sq.Select("u.id", "u.username", "u.email", "u.created_at", "u.updated_at").
		From("users u").
		Join("user_rooms ur ON u.id = ur.user_id").
		Where(sq.Eq{"ur.room_id": roomID}).
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

// GetUserRooms retrieves all rooms for a user
func (r *RoomRepository) GetUserRooms(ctx context.Context, userID int64) ([]*entity.Room, error) {
	query, args, err := sq.Select("r.id", "r.name", "r.created_at", "r.updated_at").
		From("rooms r").
		Join("user_rooms ur ON r.id = ur.room_id").
		Where(sq.Eq{"ur.user_id": userID}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}
	defer rows.Close()

	var rooms []*entity.Room
	for rows.Next() {
		var room entity.Room
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.CreatedAt,
			&room.UpdatedAt,
		)
		if err != nil {
			return nil, utils.NewSystemError(query, args, err)
		}
		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return rooms, nil
}
