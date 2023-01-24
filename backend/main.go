package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/OldPanda/green-tiles/backend/services"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

type Headers struct {
	ContentType string `json:"Content-Type,omitempty"`
}

type APIRequest struct {
	Headers           Headers `json:"headers"`
	QueryStringParams Params  `json:"queryStringParameters"`
}

type Params struct {
	Username string `json:"username"`
}

type APIResponse struct {
	StatusCode int     `json:"statusCode"`
	Headers    Headers `json:"headers"`
	Body       string  `json:"body"`
}

type ResponseMessage struct {
	Err  string `json:"error,omitempty"`
	Data string `json:"data,omitempty"`
}

func HandleLambdaEvent(ctx context.Context, eventJSON json.RawMessage) (APIResponse, error) {
	var request APIRequest
	if err := json.Unmarshal(eventJSON, &request); err != nil {
		log.Error().
			Str("Parameters", string(eventJSON)).
			Str("Err", err.Error()).
			Msg("Failed to parse request")
		return APIResponse{
			StatusCode: http.StatusBadRequest,
			Body:       getErrMessage("Invalid request"),
			Headers: Headers{
				ContentType: "application/json",
			},
		}, nil
	}

	username := request.QueryStringParams.Username
	contributions, err := services.GetAllContributions(username)
	if err != nil {
		log.Error().
			Str("Err", err.Error()).
			Msg("Failed to fetch contributions from GitHub")
		return APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       getErrMessage("Internal server error"),
			Headers: Headers{
				ContentType: "application/json",
			},
		}, nil
	}

	respBytes, err := json.Marshal(contributions)
	if err != nil {
		log.Error().
			Str("Err", err.Error()).
			Msg("Failed to marshal response")
		return APIResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       getErrMessage("Internal server error"),
			Headers: Headers{
				ContentType: "application/json",
			},
		}, nil
	}

	return APIResponse{
		StatusCode: http.StatusOK,
		Body:       getDataMessage(string(respBytes)),
		Headers: Headers{
			ContentType: "application/json",
		}}, nil
}

func getDataMessage(msg string) string {
	r := &ResponseMessage{
		Data: msg,
	}
	if bytes, err := json.Marshal(r); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

func getErrMessage(msg string) string {
	r := &ResponseMessage{
		Err: msg,
	}
	if bytes, err := json.Marshal(r); err != nil {
		return ""
	} else {
		return string(bytes)
	}
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
