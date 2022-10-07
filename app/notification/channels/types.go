package channels

import "archtecture/app/notification"

type MailMessageData struct {
	From    string
	To      string
	Subject string
	Body    string
}

type mailMessage interface {
	notification.Message
	ToMail(notification.Notifiable) MailMessageData
}

type mailNotifiable interface {
	notification.Notifiable
	RouteNotificationForMail() string
}

type smsMessage interface {
	notification.Message
	ToSms(notification.Notifiable) string
}

type smsNotifiable interface {
	notification.Notifiable
	RouteNotificationForSms() string
}

type smsClient interface {
	send(text string, to string) error
}
