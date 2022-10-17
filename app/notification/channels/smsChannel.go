package channels

import (
	"archtecture/app/env"
	"archtecture/app/notification"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func NewSmsFromEnv(env *env.Env) *SmsChannel {
	return NewSmsChannel(NewSmsOptions(env))
}

type SmsOptions struct {
	InfobipUsername string
	InfobipPassword string
	TermiiApiKey    string
	From            string
	TermiiUri       string
	InfobipUri      string
}

func NewSmsOptions(env *env.Env) *SmsOptions {
	return &SmsOptions{
		InfobipUri:      env.InfobipUri,
		InfobipUsername: env.InfobipUsername,
		InfobipPassword: env.InfobipPassword,
		TermiiApiKey:    env.TermiiApiKey,
		TermiiUri:       env.TermiiUri,
		From:            env.SmsFrom,
	}
}

type SmsChannel struct {
	options *SmsOptions
}

func NewSmsChannel(options *SmsOptions) *SmsChannel {
	return &SmsChannel{
		options: options,
	}
}

func (s *SmsChannel) Send(notifiable notification.Notifiable, message notification.Message) error {
	smsNotifiable := notifiable.(SmsNotifiable)
	return s.SendSms(message, smsNotifiable)
}

func (s *SmsChannel) SendSms(message notification.Message, smsNotifiable SmsNotifiable) error {
	sms := message.(SmsMessage).ToSms(smsNotifiable)
	to := smsNotifiable.RouteNotificationForSms()

	err := s.smsClientFor(to).
		send(sms, to)

	if err != nil {
		return err
	}

	return nil
}

func (s *SmsChannel) smsClientFor(to string) smsClient {
	clients := map[string]smsClient{
		"NG":      newTermii(s.options),
		"default": newInfobip(s.options),
	}

	return clients[s.detectCountry(to)]
}

func (s *SmsChannel) detectCountry(to string) string {
	if strings.HasPrefix(to, "+234") || strings.HasPrefix(to, "234") {
		return "NG"
	}

	return "default"
}

type smsBaseClient struct {
	options *SmsOptions
}

type infobip struct {
	*smsBaseClient
}

func newInfobip(options *SmsOptions) *infobip {
	return &infobip{smsBaseClient: &smsBaseClient{options: options}}
}

func (i *infobip) send(text string, to string) error {
	options := i.smsBaseClient.options
	a := fiber.AcquireAgent()

	request := a.Request()
	request.Header.SetMethod(fiber.MethodGet)
	request.SetRequestURI(options.InfobipUri)

	queryString := "username=" + options.InfobipUsername +
		"&password=" + options.InfobipPassword +
		"&from=" + options.From + "&to=" + to + "&text=" + text
	a.QueryString(queryString)

	if err := a.Parse(); err != nil {
		return err
	}

	return nil
}

type termii struct {
	*smsBaseClient
}

func newTermii(options *SmsOptions) *termii {
	return &termii{smsBaseClient: &smsBaseClient{options: options}}
}

func (i *termii) send(text string, to string) error {
	options := i.smsBaseClient.options
	a := fiber.AcquireAgent()
	request := a.Request()
	request.Header.SetMethod(fiber.MethodPost)
	request.SetRequestURI(options.TermiiUri)

	a.JSON(fiber.Map{
		"api_key": options.TermiiApiKey,
		"from":    options.From,
		"to":      to,
		"sms":     text,
		"type":    "plain",
		"channel": "generic",
	})

	if err := a.Parse(); err != nil {
		return err
	}

	return nil
}
