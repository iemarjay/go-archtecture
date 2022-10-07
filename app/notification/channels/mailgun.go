package channels

import (
	"archtecture/app/env"
	"archtecture/app/notification"
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

func NewMailgunFromEnv(env *env.Env) {
	NewMailGun(NewMailgunOptions(env))
}

type MailgunOptions struct {
	Domain     string
	PrivateKey string
}

func NewMailgunOptions(env *env.Env) *MailgunOptions {
	return &MailgunOptions{
		Domain:     env.MailgunDomain,
		PrivateKey: env.MailgunPrivateKey,
	}
}

type MailGun struct {
	options *MailgunOptions
}

func NewMailGun(options *MailgunOptions) *MailGun {
	return &MailGun{options: options}
}

func (mg *MailGun) Send(notifiable notification.Notifiable, message notification.Message) error {
	messageData := message.(mailMessage).ToMail(notifiable)
	user := notifiable.(mailNotifiable)

	to := messageData.To
	if to == "" {
		to = user.RouteNotificationForMail()
	}

	client := mailgun.NewMailgun(mg.options.Domain, mg.options.PrivateKey)
	clientMessage := client.NewMessage(messageData.From, messageData.Subject, messageData.Body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := client.Send(ctx, clientMessage)

	return err
}
