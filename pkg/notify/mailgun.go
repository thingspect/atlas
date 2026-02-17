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
	"github.com/thingspect/atlas/pkg/consterr"
)

const (
	emailURL                  = "https://api.mailgun.net/v3/%s/messages"
	errMailgun consterr.Error = "mailgun"
)

// mailgun contains fields and methods of a Mailgun client.
type mailgun struct {
	domain string
	apiKey string
}

// mailgunError represents a Mailgun response as returned from an API call.
type mailgunError struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// Error returns an error as a string and implements the error interface.
func (te *mailgunError) Error() string {
	return te.Message
}

// sendEmail calls the Messages API to send an email.
func (t *mailgun) sendEmail(
	ctx context.Context, from, to, subject, body string,
) error {
	// Create request.
	vals := url.Values{}
	vals.Set("from", from)
	vals.Set("to", to)
	vals.Set("subject", subject)
	vals.Set("text", body)
	r := strings.NewReader(vals.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf(emailURL, t.domain), r)
	if err != nil {
		return err
	}
	req.SetBasicAuth("api", t.apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send request.
	client := &http.Client{}
	//nolint:gosec // Built from constants and configuration.
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			logger := alog.FromContext(ctx)
			logger.Errorf("sendEmail resp.Body.Close: %v", err)
		}
	}()

	// Read response and decode.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		te := &mailgunError{}
		// Handle Mailgun mixing JSON and plain text responses.
		if err = json.Unmarshal(respBody, te); err != nil {
			return fmt.Errorf("%w: %d - %s", errMailgun, resp.StatusCode,
				respBody)
		}

		return te
	}

	return nil
}
