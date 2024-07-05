package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"listrak-api-sandbox/db"
	"listrak-api-sandbox/view/sms"
)

type ContactRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

func SmsCreateContact(w http.ResponseWriter, r *http.Request) error {
	shortCodeStr := chi.URLParam(r, "shortCode")
	shortCode, err := strconv.Atoi(shortCodeStr)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	phoneListSrt := chi.URLParam(r, "phoneList")
	phoneList, err := strconv.Atoi(phoneListSrt)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	var contactReq ContactRequest
	if err := json.NewDecoder(r.Body).Decode(&contactReq); err != nil {
		RespondWithError(w, err)
		return err
	}

	ctx := r.Context()
	_, err = db.GetContactByPhone(ctx, shortCode, contactReq.PhoneNumber)
	if err != nil && err.Error() != "sql: no rows in result set" {
		RespondWithError(w, err)
		return err
	}

	if err == nil {
		Respond(
			w,
			map[string]interface{}{
				"status":  http.StatusBadRequest,
				"error":   "ERROR_PENDING_PHONE_NUMBER",
				"message": "The phone number provided is already pending subscription through double opt in.",
			},
			http.StatusBadRequest,
		)
		return nil
	}

	contact, err := db.CreatePendingContact(ctx, shortCode, phoneList, contactReq.PhoneNumber)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	responseData := map[string]interface{}{
		"status":     http.StatusCreated,
		"resourceId": contact.PhoneNumber,
	}
	Respond(w, responseData, http.StatusCreated)

	return nil
}

func SmsGetContact(w http.ResponseWriter, r *http.Request) error {
	shortCodeStr := chi.URLParam(r, "shortCode")
	phoneNumber := chi.URLParam(r, "phoneNumber")

	shortCode, err := strconv.Atoi(shortCodeStr)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	ctx := r.Context()
	contact, err := db.GetContactByPhone(ctx, shortCode, phoneNumber)
	if err != nil && err.Error() == "sql: no rows in result set" {
		responseData := map[string]interface{}{
			"status":  http.StatusNotFound,
			"error":   "ERROR_UNABLE_TO_LOCATE_RESOURCE",
			"message": "Unable to locate a resource associated with the shortCodeId and phoneNumber supplied.",
		}
		Respond(w, responseData, http.StatusNotFound)
		return nil
	} else if err != nil {
		RespondWithError(w, err)
	}

	responseData := map[string]interface{}{
		"status": 200,
		"data": map[string]interface{}{
			"subscribeDate":   "",
			"unsubscribeDate": nil,
			"phoneNumber":     contact.PhoneNumber,
			"emailAddress":    "",
			"firstName":       "",
			"lastName":        "",
			"birthday":        "",
			"postalCode":      "",
			"optedOut":        false,
		},
	}

	Respond(w, responseData, http.StatusOK)

	return nil
}

func SmsSubscribeContact(w http.ResponseWriter, r *http.Request) error {
	shortCodeStr := chi.URLParam(r, "shortCode")
	shortCode, err := strconv.Atoi(shortCodeStr)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	phoneListSrt := chi.URLParam(r, "phoneList")
	phoneList, err := strconv.Atoi(phoneListSrt)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	phoneNumber := chi.URLParam(r, "phoneNumber")
	if len(phoneNumber) == 0 {
		err = errors.New("phoneNumber is required")
		RespondWithError(w, err)
		return err
	}

	ctx := r.Context()
	contact, err := db.GetContactByPhone(ctx, shortCode, phoneNumber)
	if err != nil && err.Error() == "sql: no rows in result set" {
		responseData := map[string]interface{}{
			"status":  http.StatusNotFound,
			"error":   "ERROR_UNABLE_TO_LOCATE_RESOURCE",
			"message": "Unable to locate a resource associated with the shortCodeId and phoneNumber supplied.",
		}
		Respond(w, responseData, http.StatusNotFound)
		return nil
	} else if err != nil {
		RespondWithError(w, err)
		return err
	}

	subscription, err := db.GetContactSubscriptionByPhoneList(ctx, contact.ID, phoneList)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	if subscription.PendingDoubleOptIn {
		responseData := map[string]interface{}{
			"status":  http.StatusBadRequest,
			"error":   "ERROR_PENDING_PHONE_NUMBER",
			"message": "The phone number provided is already pending subscription through double opt in.",
		}
		Respond(w, responseData, http.StatusBadRequest)
		return nil
	}

	if subscription.Subscribed {
		responseData := map[string]interface{}{
			"status":  http.StatusBadRequest,
			"error":   "ERROR_SUBSCRIBED_PHONE_NUMBER",
			"message": "The phone number provided is already active subscription.",
		}
		Respond(w, responseData, http.StatusBadRequest)
		return nil
	}

	subscription.PendingDoubleOptIn = true
	subscription.Subscribed = true
	subscription.SubscribeDate = time.Now()

	err = db.UpdateContactSubscription(ctx, subscription)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	responseData := map[string]interface{}{
		"status":     http.StatusCreated,
		"resourceId": contact.PhoneNumber,
	}

	Respond(w, responseData, http.StatusCreated)

	return nil
}

func SmsUnsubscribeContact(w http.ResponseWriter, r *http.Request) error {
	shortCodeStr := chi.URLParam(r, "shortCode")
	shortCode, err := strconv.Atoi(shortCodeStr)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	phoneListSrt := chi.URLParam(r, "phoneList")
	phoneList, err := strconv.Atoi(phoneListSrt)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	phoneNumber := chi.URLParam(r, "phoneNumber")
	if len(phoneNumber) == 0 {
		err = errors.New("phoneNumber is required")
		RespondWithError(w, err)
		return err
	}

	ctx := r.Context()
	contact, err := db.GetContactByPhone(ctx, shortCode, phoneNumber)
	if err != nil && err.Error() == "sql: no rows in result set" {
		responseData := map[string]interface{}{
			"status":  http.StatusNotFound,
			"error":   "ERROR_UNABLE_TO_LOCATE_RESOURCE",
			"message": "Unable to locate a resource associated with the shortCodeId and phoneNumber supplied.",
		}
		Respond(w, responseData, http.StatusNotFound)
		return nil
	} else if err != nil {
		RespondWithError(w, err)
		return err
	}

	subscription, err := db.GetContactSubscriptionByPhoneList(ctx, contact.ID, phoneList)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	if subscription.PendingDoubleOptIn || !subscription.Subscribed {
		responseData := map[string]interface{}{
			"status":  http.StatusBadRequest,
			"error":   "ERROR_PHONE_NUMBER_NOT_SUBSCRIBED",
			"message": "The phone number provided is not subscribed to the list.",
		}
		Respond(w, responseData, http.StatusBadRequest)
		return nil
	}

	subscription.Subscribed = false

	err = db.UpdateContactSubscription(ctx, subscription)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	responseData := map[string]interface{}{
		"status":     http.StatusCreated,
		"resourceId": contact.PhoneNumber,
	}

	Respond(w, responseData, http.StatusCreated)

	return nil
}

func SmsGetContactListCollection(w http.ResponseWriter, r *http.Request) error {
	shortCodeStr := chi.URLParam(r, "shortCode")
	phoneNumber := chi.URLParam(r, "phoneNumber")

	shortCode, err := strconv.Atoi(shortCodeStr)
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	ctx := r.Context()
	contact, err := db.GetContactByPhone(ctx, shortCode, phoneNumber)
	if err != nil && err.Error() == "sql: no rows in result set" {
		responseData := map[string]interface{}{
			"status":  http.StatusNotFound,
			"error":   "ERROR_UNABLE_TO_LOCATE_RESOURCE",
			"message": "Unable to locate a resource associated with the shortCodeId and phoneNumber supplied.",
		}
		Respond(w, responseData, http.StatusNotFound)
		return nil
	} else if err != nil {
		RespondWithError(w, err)
		return err
	}

	subscriptions, err := db.GetSubscriptionsList(ctx, contact)
	if len(subscriptions) == 0 {
		responseData := map[string]interface{}{
			"status":         http.StatusOK,
			"nextPageCursor": "",
			"data": []map[string]interface{}{
				{
					"phoneListId":        0,
					"subscriptionState":  "Not Subscribed",
					"subscribeDate":      contact.CreatedAt,
					"pendingDoubleOptIn": false,
					"pendingAgeGate":     false,
				},
			},
		}

		Respond(w, responseData, http.StatusOK)
	} else {
		subscriptionsData := make([]map[string]interface{}, len(subscriptions))
		for i, subscription := range subscriptions {
			subscribedStatus := "Not Subscribed"
			if subscription.Subscribed && !subscription.PendingDoubleOptIn {
				subscribedStatus = "Subscribed"
			}
			subscriptionsData[i] = map[string]interface{}{
				"phoneListId":        subscription.PhoneList,
				"subscriptionState":  subscribedStatus,
				"subscribeDate":      contact.CreatedAt,
				"pendingDoubleOptIn": subscription.PendingDoubleOptIn,
				"pendingAgeGate":     false,
			}
		}

		responseData := map[string]interface{}{
			"status":         http.StatusOK,
			"nextPageCursor": "",
			"data":           subscriptionsData,
		}

		Respond(w, responseData, http.StatusOK)
	}

	return nil
}

func HandleSmsIndex(w http.ResponseWriter, r *http.Request) error {
	contacts, err := db.GetContactsDto(r.Context())
	if err != nil {
		return err
	}

	return render(w, r, sms.Index(contacts))
}

func HandleSubscriptionUpdate(w http.ResponseWriter, r *http.Request) error {
	subscriptionIdParam := chi.URLParam(r, "subscriptionId")
	subscriptionId, err := strconv.Atoi(subscriptionIdParam)
	if err != nil {
		return render(w, r, sms.SubscriptionError(err))
	}

	subscription, err := db.GetContactSubscriptionById(r.Context(), subscriptionId)
	if err != nil {
		return render(w, r, sms.SubscriptionError(err))
	}

	var value bool
	valueParam := chi.URLParam(r, "value")
	if valueParam == "1" {
		value = true
	} else {
		value = false
	}

	fieldNameParam := chi.URLParam(r, "fieldName")
	if fieldNameParam == "subscribed" {
		subscription.Subscribed = value
		if value {
			subscription.PendingDoubleOptIn = true
		} else {
			subscription.PendingDoubleOptIn = false
		}
	} else {
		subscription.PendingDoubleOptIn = value
		if value {
			subscription.Subscribed = true
		}
	}

	err = db.UpdateContactSubscription(r.Context(), subscription)
	if err != nil {
		return render(w, r, sms.SubscriptionError(err))
	}

	return render(w, r, sms.Subscription(subscription))
}

func HandleContactDelete(w http.ResponseWriter, r *http.Request) error {
	contactIdParam := chi.URLParam(r, "contactId")
	contactId, err := strconv.Atoi(contactIdParam)
	if err != nil {
		return render(w, r, sms.ContactError(err))
	}

	contact, err := db.GetContactById(r.Context(), contactId)
	if err != nil {
		return render(w, r, sms.ContactError(err))
	}

	err = db.DeleteContact(r.Context(), contact)
	if err != nil {
		return render(w, r, sms.ContactError(err))
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
