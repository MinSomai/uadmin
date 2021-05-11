package uadmin

import (
	"strings"
	"testing"
	"time"
)

// TestSendEmail is a unit testing function for SendEmail() function
func TestSendEmail(t *testing.T) {
	// This email should be sent
	SendEmail([]string{"user@example.com"}, []string{}, []string{}, "subject", "body")
	time.Sleep(time.Millisecond * 500)
	if receivedEmail == "" {
		t.Errorf("SendEmail didn't send an email")
	}
	if !strings.Contains(receivedEmail, "From: uadmin@example.com") {
		t.Errorf("SendEmail don't have a valid From")
	}
	if !strings.Contains(receivedEmail, "To: user@example.com") {
		t.Errorf("SendEmail don't have a valid To")
	}
	if !strings.Contains(receivedEmail, "Subject: subject") {
		t.Errorf("SendEmail don't have a valid Subject")
	}
	if !strings.Contains(receivedEmail, "body") {
		t.Errorf("SendEmail don't have a valid body")
	}
	receivedEmail = ""

	// Not try sending an email with missing settings
	temp := EmailUsername
	EmailUsername = ""
	err := SendEmail([]string{"user@example.com"}, []string{}, []string{}, "subject", "body")
	if err == nil {
		t.Errorf("SendEmail send an email with missing settings")
	}
	EmailUsername = temp
}
