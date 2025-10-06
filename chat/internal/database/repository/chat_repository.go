package repository

import (
	"context"
	"database/sql"
	"playgoround/chat/internal/database/entity"
	"playgoround/chat/internal/database/utils"

	sq "github.com/Masterminds/squirrel"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// Create inserts a new chat message
func (r *ChatRepository) Create(ctx context.Context, chat *entity.Chat) (*entity.Chat, error) {
	query, args, err := sq.Insert("chats").
		Columns("room_id", "user_id", "message", "created_at").
		Values(chat.RoomID, chat.UserID, chat.Message, chat.CreatedAt).
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

	chat.ID = id
	return chat, nil
}

// FindByID retrieves a chat message by ID
func (r *ChatRepository) FindByID(ctx context.Context, id int64) (*entity.Chat, error) {
	query, args, err := sq.Select("id", "room_id", "user_id", "message", "created_at").
		From("chats").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var chat entity.Chat
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&chat.ID,
		&chat.RoomID,
		&chat.UserID,
		&chat.Message,
		&chat.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return &chat, nil
}

// GetRoomChats retrieves all chat messages for a room with pagination
func (r *ChatRepository) GetRoomChats(ctx context.Context, roomID int64, limit, offset uint64) ([]*entity.Chat, error) {
	query, args, err := sq.Select("id", "room_id", "user_id", "message", "created_at").
		From("chats").
		Where(sq.Eq{"room_id": roomID}).
		OrderBy("created_at ASC").
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

	var chats []*entity.Chat
	for rows.Next() {
		var chat entity.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.RoomID,
			&chat.UserID,
			&chat.Message,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, utils.NewSystemError(query, args, err)
		}
		chats = append(chats, &chat)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return chats, nil
}

// GetUserChats retrieves all chat messages by a user with pagination
func (r *ChatRepository) GetUserChats(ctx context.Context, userID int64, limit, offset uint64) ([]*entity.Chat, error) {
	query, args, err := sq.Select("id", "room_id", "user_id", "message", "created_at").
		From("chats").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
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

	var chats []*entity.Chat
	for rows.Next() {
		var chat entity.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.RoomID,
			&chat.UserID,
			&chat.Message,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, utils.NewSystemError(query, args, err)
		}
		chats = append(chats, &chat)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return chats, nil
}

// Delete deletes a chat message by ID
func (r *ChatRepository) Delete(ctx context.Context, id int64) (*entity.Chat, error) {
	query, args, err := sq.Select("id", "room_id", "user_id", "message", "created_at").
		From("chats").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	var chat entity.Chat
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&chat.ID,
		&chat.RoomID,
		&chat.UserID,
		&chat.Message,
		&chat.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, utils.NewNotFoundError(query, args, err)
	}
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	deleteQuery, deleteArgs, err := sq.Delete("chats").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	_, err = r.db.ExecContext(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	return &chat, nil
}

// DeleteRoomChats deletes all chat messages in a room
func (r *ChatRepository) DeleteRoomChats(ctx context.Context, roomID int64) ([]*entity.Chat, error) {
	// First, get all chats to return
	selectQuery, selectArgs, err := sq.Select("id", "room_id", "user_id", "message", "created_at").
		From("chats").
		Where(sq.Eq{"room_id": roomID}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(selectQuery, selectArgs, err)
	}

	rows, err := r.db.QueryContext(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, utils.NewSystemError(selectQuery, selectArgs, err)
	}
	defer rows.Close()

	var chats []*entity.Chat
	for rows.Next() {
		var chat entity.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.RoomID,
			&chat.UserID,
			&chat.Message,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, utils.NewSystemError(selectQuery, selectArgs, err)
		}
		chats = append(chats, &chat)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewSystemError(selectQuery, selectArgs, err)
	}

	// Then delete
	deleteQuery, deleteArgs, err := sq.Delete("chats").
		Where(sq.Eq{"room_id": roomID}).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	_, err = r.db.ExecContext(ctx, deleteQuery, deleteArgs...)
	if err != nil {
		return nil, utils.NewSystemError(deleteQuery, deleteArgs, err)
	}

	return chats, nil
}

// GetRecentChats retrieves the most recent chat messages across all rooms
func (r *ChatRepository) GetRecentChats(ctx context.Context, limit uint64) ([]*entity.Chat, error) {
	query, args, err := sq.Select("id", "room_id", "user_id", "message", "created_at").
		From("chats").
		OrderBy("created_at DESC").
		Limit(limit).
		ToSql()
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}
	defer rows.Close()

	var chats []*entity.Chat
	for rows.Next() {
		var chat entity.Chat
		err := rows.Scan(
			&chat.ID,
			&chat.RoomID,
			&chat.UserID,
			&chat.Message,
			&chat.CreatedAt,
		)
		if err != nil {
			return nil, utils.NewSystemError(query, args, err)
		}
		chats = append(chats, &chat)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewSystemError(query, args, err)
	}

	return chats, nil
}
