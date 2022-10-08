package messages

import (
	"archtecture/app/notification"
	"archtecture/app/notification/channels"
)

type channel interface {
	Send(notification.Notifiable, notification.Message) error
}

type user interface {
	GetFirstname() string
}

type Welcome struct {
	sms   channel
	email channel
}

func NewWelcome(sms channel, email channel) *Welcome {
	return &Welcome{sms: sms, email: email}
}

func (w *Welcome) Channels(notifiable notification.Notifiable) []notification.Channel {
	return []notification.Channel{w.sms, w.email}
}

func (w *Welcome) ToMail(notifiable channels.MailNotifiable) *channels.MailMessageData {
	return &channels.MailMessageData{
		To:      notifiable.RouteNotificationForMail(),
		Subject: "Welcome to Architecture",
		Body:    "Hi " + notifiable.(user).GetFirstname() + ", Welcome to Architecture",
	}
}

func (w *Welcome) ToSms(notifiable channels.SmsNotifiable) string {
	return "Hi " + notifiable.(user).GetFirstname() + ", Welcome to Architecture"
}
