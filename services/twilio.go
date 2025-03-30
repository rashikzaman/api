package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SendSMS(to, message, twilioAccountSid, twilioPhoneNumber, twilioAuthToken string) error {
	// Twilio API endpoint for sending SMS
	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", twilioAccountSid)

	// Set up form data
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", twilioPhoneNumber)
	msgData.Set("Body", message)

	// Create HTTP request
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(msgData.Encode()))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.SetBasicAuth(twilioAccountSid, twilioAuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
