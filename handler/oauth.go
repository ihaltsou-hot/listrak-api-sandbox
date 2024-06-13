package handler

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
)

func generateBearerToken(length int) (string, error) {
	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode the random bytes to base64
	token := base64.RawURLEncoding.EncodeToString(randomBytes)

	return token, nil
}

func CreateOAuthToken(w http.ResponseWriter, r *http.Request) error {
	token, err := generateBearerToken(128) // Adjust the length as needed
	if err != nil {
		RespondWithError(w, err)
		return err
	}

	response := map[string]interface{}{
		"access_token": token,
		"token_type":   "bearer",
		"expires_in":   3599,
	}

	Respond(w, response, http.StatusOK)
	return nil
}
