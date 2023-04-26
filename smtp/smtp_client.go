package smtp

import (
	"crypto/tls"
	"time"

	mail "github.com/xhit/go-simple-mail/v2"
)

// Config microservice config
type Config struct {
	MailService MailService
}

// MailService config
type MailService struct {
	URL            string
	From           string
	Host           string
	Port           int
	Username       string
	Password       string
	KeepAlive      bool
	ConnectTimeout time.Duration
	SendTimeout    time.Duration
}

// MailData for send email
type MailData struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

// SMTPClient interface
type SMTPClient interface {
	SendMail(mail *MailData) error
}

type smtpClient struct {
	cfg *Config
}

// NewSmtpClient constructor
func NewSmtpClient(cfg *Config) *smtpClient {
	return &smtpClient{cfg: cfg}
}

// NewEmailSMTPClient connect to mail server and returns SMTP client
func (s *smtpClient) getConn() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()

	// SMTP Server
	server.Host = s.cfg.MailService.Host
	server.Port = s.cfg.MailService.Port
	server.Username = s.cfg.MailService.Username
	server.Password = s.cfg.MailService.Password
	server.ConnectTimeout = s.cfg.MailService.ConnectTimeout * time.Second
	server.SendTimeout = s.cfg.MailService.SendTimeout * time.Second
	server.KeepAlive = false
	server.Encryption = mail.EncryptionTLS
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	server.Authentication = mail.AuthPlain

	return server.Connect()
}

// SendMail send simple email with text message
func (s *smtpClient) SendMail(mailData *MailData) error {
	conn, err := s.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	msg := mail.NewMSG()
	msg.SetFrom(mailData.From)
	msg.AddTo(mailData.To)
	msg.SetSubject(mailData.Subject)
	msg.SetBody(mail.TextPlain, mailData.Content)

	return msg.Send(conn)
}
