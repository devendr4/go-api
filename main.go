package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

type email struct {
	Name    string
	Email   string
	Subject string
	Msg     string
}

func postEmail(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newEmail email
	err := decoder.Decode(&newEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "status-code": "500"})
	} else {
		mailErr := sendEmail(newEmail)
		if mailErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error", "status-code": "500"})
		} else {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "status-code": "200"})
		}
	}
}

func sendEmail(emailData email) (err error) {
	from := os.Getenv("FROM")
	password := os.Getenv("PASSWORD")

	to := []string{
		os.Getenv("TO"),
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("To: angel.emo.altuve@gmail.com\r\n" +
		"Subject: " + emailData.Subject + "\r\n" +
		"\r\n" +
		"From: " + emailData.Email + "\r\n" +
		"Name: " + emailData.Name + "\r\n" +
		"Message: " + emailData.Msg + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func main() {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Printf("Error loading .env")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/new-email", postEmail).Methods("POST", "OPTIONS")

	handler := cors.Default().Handler(router)

	fmt.Println("Listening on port " + os.Getenv("PORT") + " ...")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
}
