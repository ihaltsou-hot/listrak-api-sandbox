package types

import "time"

type Subscription struct {
	ID                 int       `bun:"id,pk,autoincrement"`
	ContactId          int       `bun:"contact_id,notnull"`
	PhoneList          int       `bun:"phone_list,notnull"`
	Subscribed         bool      `bun:"subscribed,notnull"`
	PendingDoubleOptIn bool      `bun:"pending_double_opt_in,notnull"`
	SubscribeDate      time.Time `bun:"subscribe_date"`
	CreatedAt          time.Time `bun:"created_at,default:now()"`
	UpdatedAt          time.Time `bun:"updated_at,default:now()"`
}
