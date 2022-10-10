package notification

import "archtecture/app/utils"

type Channel interface {
	Send(Notifiable, Message) error
}

type Message interface {
	Channels(Notifiable) []Channel
}

type Notifiable interface{}

type Notifier interface {
	Notify(Notifiable) Notifier
	That(Message) error
}

type DefaultNotifier struct {
	notifiable Notifiable
}

func NewDefaultNotifier() *DefaultNotifier {
	return &DefaultNotifier{}
}

func NewDefaultNotifierWithNotifiable(notifiable Notifiable) *DefaultNotifier {
	notifier := NewDefaultNotifier()
	notifier.notifiable = notifiable
	return notifier
}

func (n *DefaultNotifier) Notify(notifiable Notifiable) Notifier {
	return &DefaultNotifier{notifiable: notifiable}
}

func (n *DefaultNotifier) That(message Message) error {
	channels := message.Channels(n.notifiable)

	errors := utils.NewErrorBag()
	for _, channel := range channels {
		err := channel.Send(n.notifiable, message)
		errors.Add(channel, err)
	}

	return errors
}
