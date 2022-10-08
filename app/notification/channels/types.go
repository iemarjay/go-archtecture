package channels

import "archtecture/app/notification"

type MailMessageData struct {
	From    string
	To      string
	Subject string
	Body    string
}

type MailMessage interface {
	notification.Message
	ToMail(MailNotifiable) *MailMessageData
}

type MailNotifiable interface {
	notification.Notifiable
	RouteNotificationForMail() string
}

type SmsMessage interface {
	notification.Message
	ToSms(SmsNotifiable) string
}

type SmsNotifiable interface {
	notification.Notifiable
	RouteNotificationForSms() string
}

type smsClient interface {
	send(text string, to string) error
}
