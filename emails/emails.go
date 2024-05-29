package handler

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

const SMTPServer = "smtp.gmail.com"

type Sender struct {
	User     string
	Password string
}

func NewSender(username, password string) Sender {
	return Sender{
		User:     username,
		Password: password,
	}
}

func (sender Sender) SendMail(Dest []string, Subject, bodyMessage string) {

	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(SMTPServer+":587",
		smtp.PlainAuth("", sender.User, sender.Password, SMTPServer),
		sender.User, Dest, []byte(msg))

	if err != nil {

		fmt.Printf("smtp error: %s", err)
		return
	}

	fmt.Println("Mail sent successfully!")
}

func (sender Sender) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {

	header := make(map[string]string)
	header["From"] = sender.User

	receipient := ""

	for _, user := range dest {
		receipient = receipient + user
	}

	// header["To"] = receipient
	// header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = fmt.Sprintf("%s; charset=\"utf-8\"", contentType)
	header["Content-Transfer-Encoding"] = "quoted-printable"
	header["Content-Disposition"] = "inline"

	message := ""

	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer

	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	return message
}

func (sender *Sender) WriteHTMLEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/html", subject, bodyMessage)
}

func (sender *Sender) WritePlainEmail(dest []string, subject, bodyMessage string) string {

	return sender.WriteEmail(dest, "text/plain", subject, bodyMessage)
}

func Handler(w http.ResponseWriter, r *http.Request) {
  senderEmail := os.Getenv("SENDER_EMAIL")
  appPassword := os.Getenv("APP_PASSWORD")
	sender := NewSender(senderEmail, appPassword)
	receiver := []string{"hamanecisse2@gmail.com"}
	body := sender.WritePlainEmail(receiver, "Alerte", "Seul atteint")
	sender.SendMail(receiver, "Alerte", body)
	w.WriteHeader(http.StatusOK)
	return
}