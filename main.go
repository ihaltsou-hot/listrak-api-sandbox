package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"listrak-api-sandbox/db"
	"listrak-api-sandbox/handler"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	if err := initEverything(); err != nil {
		slog.Error("error of the initialization", "err", err)
		os.Exit(1)
	}

	router := chi.NewMux()

	router.Post("/sms/v1/ShortCode/{shortCode:[0-9]+}/PhoneList/{phoneList:[0-9]+}/Contact", handler.Make(handler.SmsCreateContact))
	router.Get("/sms/v1/ShortCode/{shortCode:[0-9]+}/Contact/{phoneNumber:[0-9]{11}}", handler.Make(handler.SmsGetContact))
	router.Post("/sms/v1/ShortCode/{shortCode:[0-9]+}/Contact/{phoneNumber:[0-9]{11}}/PhoneList/{phoneList:[0-9]+}", handler.Make(handler.SmsSubscribeContact))
	router.Delete("/sms/v1/ShortCode/{shortCode:[0-9]+}/ContactUnsubscribe/{phoneNumber:[0-9]{11}}/PhoneList/{phoneList:[0-9]+}", handler.Make(handler.SmsUnsubscribeContact))
	router.Get("/sms/v1/ShortCode/{shortCode:[0-9]+}/Contact/{phoneNumber:[0-9]{11}}/PhoneList", handler.Make(handler.SmsGetContactListCollection))

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := db.Init(); err != nil {
		return err
	}

	return nil
}
