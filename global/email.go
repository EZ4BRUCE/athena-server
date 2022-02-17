package global

import (
	"time"

	"github.com/EZ4BRUCE/athena-server/pkg/email"
)

var (
	MailerPool chan *email.Mailer
)

func GetMailer() *email.Mailer {
	for {
		select {
		case mailer := <-MailerPool:
			return mailer
		case <-time.After(time.Second):
		}
	}
}

func PutMailer(mailer *email.Mailer) {
	select {
	case MailerPool <- mailer:
		return
	default:
		Logger.Errorf("mailerPool is full...")
	}
}
