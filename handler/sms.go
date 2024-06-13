package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"listrak-api-sandbox/db"
	"net/http"
	"strconv"
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

	return nil
}

func SmsUnsubscribeContact(w http.ResponseWriter, r *http.Request) error {

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
			if subscription.Subscribed {
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
