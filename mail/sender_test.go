package mail

import (
	"github.com/stretchr/testify/require"
	"purebank/db/util"
	"testing"
)

func TestSendEmailWithGmail(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A Test Mail from Apple"
	contact := `
 <h1>Hello Chit Su Wai </h1>   
<p> This is a test message from <a href="https://www.apple.com/macbook-air/">  Apple Website</a></p>

`
	to := []string{"kyawkyaw.thar14@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, contact, to, nil, nil, attachFiles)

	require.NoError(t, err)
}
