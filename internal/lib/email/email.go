package email

import (
	"github.com/joho/godotenv"
	"log"
	"net/smtp"
	"os"
)

const (
	path    = "./config/vars.env"
	host    = "smtp.gmail.com"
	port    = "587"
	address = host + ":" + port
	subject = "Subject: Stay here!\r\n" + "\r\n" + "You've missed a lot of new stuff. Comeback and get new knowledge!\r\n"
)

var (
	email    string
	password string
)

func InitEmail() {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("[FATAL] invalid path: %s", path)
	}

	email = os.Getenv("EMAIL")
	if email == "" {
		log.Fatalf("[FATAL] EMAIL is empty")
	}
	password = os.Getenv("EMAIL_PASSWORD")
	if password == "" {
		log.Fatalf("[FATAL] EMAIL_PASSWORD is empty")
	}
}

func Send(userEmail string) error {
	var message = []byte("To: " + userEmail + "\r\n" + subject)
	to := []string{userEmail}

	auth := smtp.PlainAuth("", email, password, host)

	err := smtp.SendMail(address, auth, email, to, message)
	if err != nil {
		return err
	}

	return nil
}
