package main

import (
	"embed"
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

//go:embed public
var FS embed.FS

func main() {
	if err := initEverything(); err != nil {
		slog.Error("error of the initialization", "err", err)
		os.Exit(1)
	}

	router := chi.NewMux()

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))

	router.Post("/OAuth2/Token", handler.Make(handler.CreateOAuthToken))

	router.Post("/sms/v1/ShortCode/{shortCode:[0-9]+}/PhoneList/{phoneList:[0-9]+}/Contact", handler.Make(handler.SmsCreateContact))
	router.Get("/sms/v1/ShortCode/{shortCode:[0-9]+}/Contact/{phoneNumber:[0-9]{11}}", handler.Make(handler.SmsGetContact))
	router.Post("/sms/v1/ShortCode/{shortCode:[0-9]+}/Contact/{phoneNumber:[0-9]{11}}/PhoneList/{phoneList:[0-9]+}", handler.Make(handler.SmsSubscribeContact))
	router.Delete("/sms/v1/ShortCode/{shortCode:[0-9]+}/ContactUnsubscribe/{phoneNumber:[0-9]{11}}/PhoneList/{phoneList:[0-9]+}", handler.Make(handler.SmsUnsubscribeContact))
	router.Get("/sms/v1/ShortCode/{shortCode:[0-9]+}/Contact/{phoneNumber:[0-9]{11}}/PhoneList", handler.Make(handler.SmsGetContactListCollection))

	router.Get("/sms", handler.Make(handler.HandleSmsIndex))
	router.Put("/sms/subscription/{subscriptionId:[0-9]+}/{fieldName:(subscribed|pending)}/{value:(0|1)}/", handler.Make(handler.HandleSubscriptionUpdate))

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initEverything() error {
	_ = godotenv.Load()

	if err := db.Init(); err != nil {
		return err
	}

	return nil
}
