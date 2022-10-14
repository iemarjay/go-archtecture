package listeners

import (
	"archtecture/app/notification"
	"archtecture/users/logic"
	"archtecture/users/messages"
)

type notifier interface {
	Notify(notification.Notifiable) notification.Notifier
	That(notification.Message) error
}

type SendWelcomeNotification struct {
	notifier notifier
	message  *messages.Welcome
}

func NewSendWelcomeNotification(notifier notifier, message *messages.Welcome) *SendWelcomeNotification {
	return &SendWelcomeNotification{notifier: notifier, message: message}
}

func (s *SendWelcomeNotification) Handle(events ...interface{}) {
	for _, event := range events {
		user := event.(*logic.UserRegistered).User()
		go s.notifier.Notify(user).That(s.message)
	}
}
