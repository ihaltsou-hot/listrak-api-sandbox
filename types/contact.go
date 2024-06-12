package types

import "time"

type Contact struct {
	ID          int       `bun:"id,pk,autoincrement"`
	ShortCode   int       `bun:"short_code"`
	PhoneNumber string    `bun:"phone_number"`
	CreatedAt   time.Time `bun:"created_at,default:now()"`
}
