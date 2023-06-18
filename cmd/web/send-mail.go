package main

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"log"
	"my/gomodule/internal/models"
	"time"
)

func listenForMail() {
	go func() {
		for {
			m := <-app.MailChan
			sendMsg(m)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025

	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println("Connect err:", err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To)
	email.SetSubject(m.Subject)

	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Mail sent!")
	}
}
