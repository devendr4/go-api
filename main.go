package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/smtp"
)

type email struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func postEmail(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newEmail email
	err := decoder.Decode(&newEmail)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error"))
	}
	mailErr := sendEmail(newEmail)
	if mailErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(mailErr.Error()))
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func sendEmail(emailData email) (err error) {
	from := "gabrielmirabal18@gmail.com"
	password := "putalocura15"

	to := []string{
		"angel.emo.altuve@gmail.com",
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("To: angel.emo.altuve@gmail.com\r\n" +
		"Subject: " + emailData.Subject + "\r\n" +
		"\r\n" +
		"From: " + emailData.Email + "\r\n" +
		"Name: " + emailData.Name + "\r\n" +
		"Message: " + emailData.Message + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return err
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/new-email", postEmail)

	fmt.Println("Listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
