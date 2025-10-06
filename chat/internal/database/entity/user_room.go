package entity

import "time"

type UserRoom struct {
	UserID   int64     `db:"user_id"`
	RoomID   int64     `db:"room_id"`
	JoinedAt time.Time `db:"joined_at"`
}
