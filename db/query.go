package db

import (
	"context"
	"listrak-api-sandbox/types"
	"time"
)

func GetContactByPhone(ctx context.Context, shortCode int, phoneNumber string) (*types.Contact, error) {
	contact := &types.Contact{}
	err := Bun.NewSelect().
		Model(contact).
		Where("phone_number = ?", phoneNumber).
		Where("short_code = ?", shortCode).
		Scan(ctx)

	return contact, err
}

func CreatePendingContact(ctx context.Context, shortCode int, phoneList int, phoneNumber string) (*types.Contact, error) {
	contact := &types.Contact{
		ShortCode:   shortCode,
		PhoneNumber: phoneNumber,
	}

	_, err := Bun.NewInsert().
		Model(contact).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	subscription := &types.Subscription{
		ContactId:          contact.ID,
		PhoneList:          phoneList,
		Subscribed:         true,
		PendingDoubleOptIn: true,
		SubscribeDate:      time.Now(),
	}

	_, err = Bun.NewInsert().
		Model(subscription).
		Exec(ctx)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func GetSubscriptionsList(ctx context.Context, contact *types.Contact) ([]types.Subscription, error) {
	var subscriptions []types.Subscription
	err := Bun.NewSelect().
		Model(&subscriptions).
		Where("contact_id = ?", contact.ID).
		Scan(ctx)

	return subscriptions, err
}
