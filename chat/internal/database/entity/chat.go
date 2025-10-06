package entity

import "time"

type Chat struct {
	ID        int64     `db:"id"`
	RoomID    int64     `db:"room_id"`
	UserID    int64     `db:"user_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}
