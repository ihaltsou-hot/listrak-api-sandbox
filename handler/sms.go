package handler

import (
	"github.com/go-chi/chi/v5"
	"listrak-api-sandbox/db"
	"net/http"
	"strconv"
)

func SmsCreateContact(w http.ResponseWriter, r *http.Request) error {

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
	if err != nil {
		responseData := map[string]interface{}{
			"status":  http.StatusNotFound,
			"error":   "ERROR_UNABLE_TO_LOCATE_RESOURCE",
			"message": "Unable to locate a resource associated with the shortCodeId and phoneNumber supplied.",
		}
		Respond(w, responseData, http.StatusNotFound)
		return nil
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

	return nil
}
