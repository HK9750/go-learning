package main

import (
	"fmt"
)

type Notifier interface {
	Notify(message string) error
}

type EmailNotifier struct {}

func (e *EmailNotifier) Notify(message string) error {
	fmt.Println("EmailNotifier: ",message)
	return nil
}

type SmsNotifier struct {}

func (s *SmsNotifier) Notify(message string) error {
	fmt.Println("SmsNotifier: ",message)
	return nil
}

func Broadcast(notifiers []Notifier, msg string) {
	for _, notifier := range notifiers {
		notifier.Notify(msg)
	}
}

func main() {
	emailNotifier := &EmailNotifier{}
	smsNotifier := &SmsNotifier{}

	Broadcast([]Notifier{emailNotifier, smsNotifier}, "Hello")
}