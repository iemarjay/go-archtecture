package channels

import (
	"archtecture/app/env"
	"archtecture/app/notification"
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

func NewMailgunFromEnv(env *env.Env) *MailGun {
	return NewMailGun(NewMailgunOptions(env))
}

type MailgunOptions struct {
	Domain     string
	PrivateKey string
	MailFrom   string
}

func NewMailgunOptions(env *env.Env) *MailgunOptions {
	return &MailgunOptions{
		Domain:     env.MailgunDomain,
		PrivateKey: env.MailgunPrivateKey,
		MailFrom:   env.MailFrom,
	}
}

type MailGun struct {
	options *MailgunOptions
}

func NewMailGun(options *MailgunOptions) *MailGun {
	return &MailGun{options: options}
}

func (mg *MailGun) Send(notifiable notification.Notifiable, message notification.Message) error {
	mailNotifiable := notifiable.(MailNotifiable)
	messageData := message.(MailMessage).ToMail(mailNotifiable)

	return mg.SendMail(mailNotifiable, messageData)
}

func (mg *MailGun) SendMail(mailNotifiable MailNotifiable, messageData *MailMessageData) error {
	to := messageData.To
	if to == "" {
		to = mailNotifiable.RouteNotificationForMail()
	}
	from := mg.options.MailFrom
	if messageData.From != "" {
		from = messageData.From
	}

	client := mailgun.NewMailgun(mg.options.Domain, mg.options.PrivateKey)
	clientMessage := client.NewMessage(from, messageData.Subject, messageData.Body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	_, _, err := client.Send(ctx, clientMessage)

	return err
}
