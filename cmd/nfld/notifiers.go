package main

import (
	"bytes"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"text/template"

	"github.com/ameske/nfl-pickem"
)

type nullNotifier struct{}

func (n nullNotifier) Notify(to string, week int, picks []nflpickem.Pick) error { return nil }

type fsNotifier struct{}

func (n fsNotifier) Notify(to string, week int, picks []nflpickem.Pick) error {
	fd, err := os.Create(fmt.Sprintf("%s-%d.txt", to, week))
	if err != nil {
		return err
	}

	et, err := template.New("email").Parse(emailBody)
	if err != nil {
		return err
	}

	pe := struct {
		To      string
		From    string
		Subject string
		Week    int
		Picks   []nflpickem.Pick
	}{
		to, "debugserver", fmt.Sprintf("Week %d Picks", week), week, picks,
	}

	return et.Execute(fd, pe)
}

type emailNotifier struct {
	auth           smtp.Auth
	sender         string
	smtpServer     string
	smtpServerPort string
	et             *template.Template
}

func NewEmailNotifier(server, sendAsAddress, password string) (nflpickem.Notifier, error) {
	addr, port, err := net.SplitHostPort(server)
	if err != nil {
		return nil, err
	}

	et, err := template.New("email").Parse(emailBody)
	if err != nil {
		return nil, err
	}

	a := smtp.PlainAuth("",
		sendAsAddress,
		password,
		addr,
	)

	return emailNotifier{auth: a, sender: sendAsAddress, smtpServer: addr, smtpServerPort: port, et: et}, nil
}

func (e emailNotifier) Notify(to string, week int, picks []nflpickem.Pick) error {
	pe := struct {
		To      string
		From    string
		Subject string
		Week    int
		Picks   []nflpickem.Pick
	}{
		to, e.sender, fmt.Sprintf("Week %d Picks", week), week, picks,
	}

	var body bytes.Buffer
	err := e.et.Execute(&body, pe)
	if err != nil {
		return err
	}

	fullAddr := fmt.Sprintf("%s:%s", e.smtpServer, e.smtpServerPort)

	return smtp.SendMail(fullAddr, e.auth, e.sender, []string{to}, body.Bytes())
}

var emailBody = `
To: {{.To}}
From: {{.From}}
Subject: {{.Subject}}

Here are the picks that I currently have recorded in my system for Week {{.Week}}!

Please double-check and make sure there are no errors. E-mail me if you find any problems.

{{range .Picks}}
{{.AwayNickname}}/{{.HomeNick}} - {{if eq .Selection 1}}{{.AwayNickname}}{{else}}{{.HomeNickname}}{{end}} {{.Points}}
{{end}}

Good luck!

-Kyle Ames Bot
`
