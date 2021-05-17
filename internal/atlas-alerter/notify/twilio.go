package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/thingspect/atlas/pkg/alog"
)

const (
	lookupURL = "https://lookups.twilio.com/v1/PhoneNumbers/%s?Type=carrier"
	smsURL    = "https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json"
)

// twilio contains fields and methods of a Twilio client.
type twilio struct {
	keySID     string
	accountSID string
	keySecret  string
	phone      string
}

// twilioError represents a Twilio error message as returned from an API call.
type twilioError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error returns an error as a string and implements the error interface.
func (te *twilioError) Error() string {
	return fmt.Sprintf("%d - %s", te.Code, te.Message)
}

// lookup represents a response from the Lookup API.
type lookup struct {
	Carrier struct {
		Type string `json:"type"`
	} `json:"carrier"`
}

// lookupCarrier calls the Lookup API to retrieve information about a phone
// number's carrier. This function does not require a populated
// twilio.accountSID or twilio.phone.
func (t *twilio) lookupCarrier(ctx context.Context,
	phone string) (*lookup, error) {
	// Create request.
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(lookupURL,
		phone), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(t.keySID, t.keySecret)

	// Send request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("lookupCarrier resp.Body.Close: %v", err)
		}
	}()

	d := json.NewDecoder(resp.Body)

	// Read response and decode.
	if resp.StatusCode >= 400 {
		te := &twilioError{}
		if err = d.Decode(te); err != nil {
			return nil, err
		}

		return nil, te
	}

	l := &lookup{}
	if err = d.Decode(l); err != nil {
		return nil, err
	}

	return l, nil
}

// sendSMS calls the Message API to send an SMS message.
func (t *twilio) sendSMS(ctx context.Context, to, body string) error {
	// Create request.
	vals := url.Values{}
	vals.Set("From", t.phone)
	vals.Set("To", to)
	vals.Set("Body", body)
	r := strings.NewReader(vals.Encode())

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf(smsURL,
		t.accountSID), r)
	if err != nil {
		return err
	}
	req.SetBasicAuth(t.keySID, t.keySecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("lookupCarrier resp.Body.Close: %v", err)
		}
	}()

	// Read response and decode.
	if resp.StatusCode >= 400 {
		te := &twilioError{}
		d := json.NewDecoder(resp.Body)
		if err = d.Decode(te); err != nil {
			return err
		}

		return te
	}

	// Response body is not useful, discard it to drain connection.
	if _, err = io.Copy(io.Discard, resp.Body); err != nil {
		return err
	}

	return nil
}
