package db

import (
	"context"
	"listrak-api-sandbox/types"
)

func GetContactByPhone(ctx context.Context, shortCode int, phone string) (types.Contact, error) {
	var contact types.Contact
	err := Bun.NewSelect().
		Model(&contact).
		Where("phone_number = ?", phone).
		Where("short_code = ?", shortCode).
		Scan(ctx)

	return contact, err
}
