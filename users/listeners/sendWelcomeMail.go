package listeners

import (
	"archtecture/app/notification/channels"
	"archtecture/users/logic"
)

type channel interface {
	SendMail(channels.MailNotifiable, *channels.MailMessageData) error
}

type SendWelcomeMail struct {
	email channel
}

func NewSendWelcomeMail(email channel) *SendWelcomeMail {
	return &SendWelcomeMail{email: email}
}

func (s *SendWelcomeMail) Handle(payload interface{}) {
	event := payload.(logic.UserRegistered)
	user := event.User()
	message := &channels.MailMessageData{
		To:      user.RouteNotificationForMail(),
		Subject: "Welcome to Architecture",
		Body:    "Hi " + user.GetFirstname() + ", Welcome to Architecture",
	}
	_ = s.email.SendMail(user, message)
}
