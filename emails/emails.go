package emails

import (
	"anonymous/encryption"
	"bytes"
	"fmt"
	"mime/quotedprintable"
	"net/smtp"
	"strings"
)

const (
	SMTPServer           = "smtp.gmail.com"
	SMTPPort             = "587"
	MIMEVersionHeader    = "1.0"
	ContentTypePlainText = "text/plain; charset=utf-8"
	ContentTypeHTML      = "text/html; charset=utf-8"
	ContentTransferEnc   = "quoted-printable"
)

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

func (sender Sender) SendMail(dest []string, subject, message string) error {
	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", SMTPServer, SMTPPort),
		smtp.PlainAuth("", sender.User, sender.Password, SMTPServer),
		sender.User, dest, []byte(message),
	)

	if err != nil {
		return fmt.Errorf("smtp error: %w", err)
	}

	fmt.Println("Mail sent successfully!")
	return nil
}

func (sender Sender) WriteEmail(dest []string, contentType, subject, bodyMessage string) string {
	headers := map[string]string{
		"From":                     sender.User,
		"To":                       strings.Join(dest, ","),
		"Subject":                  subject,
		"MIME-Version":             MIMEVersionHeader,
		"Content-Type":             contentType,
		"Content-Transfer-Encoding": ContentTransferEnc,
		"Content-Disposition":      "inline",
	}

	message := constructMessage(headers, bodyMessage)
	return message
}

func constructMessage(headers map[string]string, bodyMessage string) string {
	var message strings.Builder
	for key, value := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	var encodedMessage bytes.Buffer
	encoder := quotedprintable.NewWriter(&encodedMessage)
	encoder.Write([]byte(bodyMessage))
	encoder.Close()

	message.WriteString("\r\n" + encodedMessage.String())
	return message.String()
}

func (sender *Sender) WriteHTMLEmail(dest []string, subject, bodyMessage string) string {
	return sender.WriteEmail(dest, ContentTypeHTML, subject, bodyMessage)
}

func (sender *Sender) WritePlainEmail(dest []string, subject, bodyMessage string) string {
	return sender.WriteEmail(dest, ContentTypePlainText, subject, bodyMessage)
}

// Vérifie si une chaîne est un email chiffré et la déchiffre si nécessaire
func ensureDecryptedEmail(email string) (string, error) {
	ee, err := encryption.NewEmailEncryption()
	if err != nil {
		return "", fmt.Errorf("erreur lors de l'initialisation du chiffrement: %w", err)
	}
	
	// Si l'email semble être chiffré, essayons de le déchiffrer
	if ee.IsEncrypted(email) {
		decryptedEmail, err := ee.DecryptEmail(email)
		if err != nil {
			return "", fmt.Errorf("erreur lors du déchiffrement de l'email: %w", err)
		}
		return decryptedEmail, nil
	}
	
	// Sinon, on suppose que c'est déjà un email en clair
	return email, nil
}

func (sender Sender) SendVerificationEmail(dest []string, token string) error {
	// S'assurer que les emails ne sont pas chiffrés
	decryptedDest := make([]string, len(dest))
	for i, email := range dest {
		decrypted, err := ensureDecryptedEmail(email)
		if err != nil {
			return err
		}
		decryptedDest[i] = decrypted
	}
	
	verificationLink := fmt.Sprintf("https://anonymous-4jef.onrender.com/auth/verify-email?token=%s", token)
	subject := "Please verify your email address"
	bodyMessage := fmt.Sprintf("Click the following link to verify your email address: <a href=\"%s\">%s</a>", verificationLink, verificationLink)
	message := sender.WriteHTMLEmail(decryptedDest, subject, bodyMessage)
	return sender.SendMail(decryptedDest, subject, message)
}

func (sender Sender) SendPasswordResetEmail(dest []string, token string) error {
	// S'assurer que les emails ne sont pas chiffrés
	decryptedDest := make([]string, len(dest))
	for i, email := range dest {
		decrypted, err := ensureDecryptedEmail(email)
		if err != nil {
			return err
		}
		decryptedDest[i] = decrypted
	}
	
    resetLink := fmt.Sprintf("https://anonymous-4jef.onrender.com/reset-password?token=%s", token)
    subject := "Password Reset Request"
    bodyMessage := fmt.Sprintf(`
        <h1>Password Reset</h1>
        <p>You requested a password reset. Click the link below to reset your password:</p>
        <p><a href="%s">Reset Password</a></p>
        <p>If you did not request a password reset, please ignore this email.</p>
        <p>This link will expire in 24 hours.</p>
    `, resetLink)
    
    message := sender.WriteHTMLEmail(decryptedDest, subject, bodyMessage)
    return sender.SendMail(decryptedDest, subject, message)
}