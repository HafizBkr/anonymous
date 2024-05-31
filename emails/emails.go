package emails

import (
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
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

func (sender Sender) SendMail(dest []string, subject, bodyMessage string) {
	header := make(map[string]string)
	header["From"] = sender.User
	header["To"] = strings.Join(dest, ",")
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=utf-8"
	header["Content-Transfer-Encoding"] = "quoted-printable"

	message := ""
	for key, value := range header {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	var encodedMessage bytes.Buffer
	finalMessage := quotedprintable.NewWriter(&encodedMessage)
	finalMessage.Write([]byte(bodyMessage))
	finalMessage.Close()

	message += "\r\n" + encodedMessage.String()

	err := smtp.SendMail(SMTPServer+":587",
		smtp.PlainAuth("", sender.User, sender.Password, SMTPServer),
		sender.User, dest, []byte(message))

	if err != nil {
		fmt.Printf("smtp error: %s", err)
		return
	}

	fmt.Println("Mail sent successfully!")
}

func (sender Sender) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {
	header := make(map[string]string)
	header["From"] = sender.User
	header["To"] = strings.Join(dest, ",")
	header["Subject"] = subject
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

func (sender Sender) SendVerificationEmail(dest []string, token string) {
	verificationLink := "http://yourdomain.com/verify?token=" + token
	subject := "Please verify your email address"
	bodyMessage := fmt.Sprintf("Click the following link to verify your email address: <a href=\"%s\">%s</a>", verificationLink, verificationLink)
	message := sender.WriteEmail(dest, "text/html", subject, bodyMessage)
	sender.SendMail(dest, subject, message)
}
