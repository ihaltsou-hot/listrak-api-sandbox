package db

import (
	"context"
	"listrak-api-sandbox/types"
	"time"
)

func GetContactByPhone(ctx context.Context, shortCode int, phoneNumber string) (types.Contact, error) {
	contact := types.Contact{}
	err := Bun.NewSelect().
		Model(&contact).
		Where("phone_number = ?", phoneNumber).
		Where("short_code = ?", shortCode).
		Scan(ctx)

	return contact, err
}

func CreatePendingContact(ctx context.Context, shortCode int, phoneList int, phoneNumber string) (types.Contact, error) {
	contact := types.Contact{
		ShortCode:   shortCode,
		PhoneNumber: phoneNumber,
	}

	_, err := Bun.NewInsert().
		Model(&contact).
		Exec(ctx)

	if err != nil {
		return contact, err
	}

	subscription := types.Subscription{
		ContactId:          contact.ID,
		PhoneList:          phoneList,
		Subscribed:         true,
		PendingDoubleOptIn: true,
		SubscribeDate:      time.Now(),
	}

	_, err = Bun.NewInsert().
		Model(&subscription).
		Exec(ctx)

	if err != nil {
		return contact, err
	}

	return contact, nil
}

func GetSubscriptionsList(ctx context.Context, contact types.Contact) ([]types.Subscription, error) {
	var subscriptions []types.Subscription
	err := Bun.NewSelect().
		Model(&subscriptions).
		Where("contact_id = ?", contact.ID).
		Scan(ctx)

	return subscriptions, err
}

func GetContactSubscriptionByPhoneList(ctx context.Context, contactId int, phoneList int) (types.Subscription, error) {
	var subscription types.Subscription
	err := Bun.NewSelect().
		Model(&subscription).
		Where("contact_id = ?", contactId).
		Where("phone_list = ?", phoneList).
		Scan(ctx)

	return subscription, err
}

func UpdateContactSubscription(ctx context.Context, subscription types.Subscription) error {
	_, err := Bun.NewUpdate().
		Model(&subscription).
		WherePK().
		Exec(ctx)

	return err
}

func GetContactsDto(ctx context.Context) ([]types.ContactDto, error) {
	var contacts []types.Contact
	err := Bun.NewSelect().
		Model(&contacts).
		Scan(ctx)

	contactsDto := make([]types.ContactDto, 0)
	for _, contact := range contacts {
		subscriptions, err := GetSubscriptionsList(ctx, contact)
		if err != nil {
			return nil, err
		}

		contactsDto = append(contactsDto, types.ContactDto{
			Contact:       contact,
			Subscriptions: subscriptions,
		})
	}

	return contactsDto, err
}
