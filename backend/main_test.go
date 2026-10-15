package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/OldPanda/green-tiles/backend/services"
)

func assertResponse(t *testing.T, response APIResponse, status int) ResponseMessage {
	t.Helper()
	if response.StatusCode != status {
		t.Fatalf("status = %d, want %d", response.StatusCode, status)
	}
	if response.Headers.ContentType != "application/json" || response.Headers.AccessControlAllowOrigin != "*" {
		t.Errorf("unexpected response headers: %+v", response.Headers)
	}
	var message ResponseMessage
	if err := json.Unmarshal([]byte(response.Body), &message); err != nil {
		t.Fatalf("invalid JSON response %q: %v", response.Body, err)
	}
	return message
}

func TestHandleLambdaEventInvalidRequest(t *testing.T) {
	for _, body := range []string{`invalid`, `[]`, `{"queryStringParameters":{"username":42}}`} {
		t.Run(body, func(t *testing.T) {
			response, err := HandleLambdaEvent(context.Background(), json.RawMessage(body))
			if err != nil {
				t.Fatal(err)
			}
			message := assertResponse(t, response, http.StatusBadRequest)
			if message.Err != "Invalid request" || message.Data != "" {
				t.Errorf("unexpected message: %+v", message)
			}
		})
	}
}

func TestHandleLambdaEventSuccess(t *testing.T) {
	want := &services.ContributionResult{
		Login: "OldPanda", AvatarURL: "https://example.com/avatar.png",
		ContributionCalendars: []services.Calendar{{Year: 2025, Total: 2, Weeks: []services.Week{
			{Days: []services.Day{{Date: "2025-01-01", Count: 2, Weekday: 3, Level: "FIRST_QUARTILE"}}},
		}}},
	}
	response, err := handleLambdaEvent(json.RawMessage(`{"queryStringParameters":{"username":"OldPanda"}}`),
		func(username string) (*services.ContributionResult, error) {
			if username != "OldPanda" {
				t.Errorf("username = %q", username)
			}
			return want, nil
		})
	if err != nil {
		t.Fatal(err)
	}
	message := assertResponse(t, response, http.StatusOK)
	if message.Err != "" {
		t.Errorf("unexpected error: %s", message.Err)
	}
	var got services.ContributionResult
	if err := json.Unmarshal([]byte(message.Data), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(&got, want) {
		t.Errorf("contributions = %+v, want %+v", got, want)
	}
}

func TestHandleLambdaEventFetchError(t *testing.T) {
	for _, body := range []string{`{}`, `{"queryStringParameters":{"username":"missing"}}`} {
		t.Run(body, func(t *testing.T) {
			response, err := handleLambdaEvent(json.RawMessage(body), func(string) (*services.ContributionResult, error) {
				return nil, errors.New("GitHub said \"not found\"\nretry later")
			})
			if err != nil {
				t.Fatal(err)
			}
			message := assertResponse(t, response, http.StatusBadRequest)
			if message.Err != "GitHub said \"not found\"\nretry later" || message.Data != "" {
				t.Errorf("unexpected message: %+v", message)
			}
		})
	}
}

func TestResponseMessages(t *testing.T) {
	for _, text := range []string{"", "message", "quotes \" and newline\n and Unicode 熊猫"} {
		t.Run(text, func(t *testing.T) {
			for _, test := range []struct {
				body string
				want ResponseMessage
			}{
				{getDataMessage(text), ResponseMessage{Data: text}},
				{getErrMessage(text), ResponseMessage{Err: text}},
			} {
				var got ResponseMessage
				if err := json.Unmarshal([]byte(test.body), &got); err != nil {
					t.Fatal(err)
				}
				if got != test.want {
					t.Errorf("message = %+v, want %+v", got, test.want)
				}
			}
		})
	}
}
