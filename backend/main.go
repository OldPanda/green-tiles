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
	ContentType              string `json:"Content-Type,omitempty"`
	AccessControlAllowOrigin string `json:"Access-Control-Allow-Origin,omitempty"`
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
	return handleLambdaEvent(eventJSON, services.GetAllContributions)
}

func handleLambdaEvent(eventJSON json.RawMessage, fetchContributions func(string) (*services.ContributionResult, error)) (APIResponse, error) {
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
				ContentType:              "application/json",
				AccessControlAllowOrigin: "*",
			},
		}, nil
	}

	username := request.QueryStringParams.Username
	log.Info().
		Str("Username", username).
		Msg("Fetching contribution records")
	contributions, err := fetchContributions(username)
	if err != nil {
		log.Error().
			Str("Err", err.Error()).
			Msg("Failed to fetch contributions from GitHub")
		return APIResponse{
			StatusCode: http.StatusBadRequest,
			Body:       getErrMessage(err.Error()),
			Headers: Headers{
				ContentType:              "application/json",
				AccessControlAllowOrigin: "*",
			},
		}, nil
	}

	// ContributionResult contains only JSON-safe strings, integers and slices.
	respBytes, _ := json.Marshal(contributions)

	return APIResponse{
		StatusCode: http.StatusOK,
		Body:       getDataMessage(string(respBytes)),
		Headers: Headers{
			ContentType:              "application/json",
			AccessControlAllowOrigin: "*",
		}}, nil
}

func getDataMessage(msg string) string {
	r := &ResponseMessage{
		Data: msg,
	}
	// ResponseMessage contains only strings, so encoding cannot fail.
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

func getErrMessage(msg string) string {
	r := &ResponseMessage{
		Err: msg,
	}
	// ResponseMessage contains only strings, so encoding cannot fail.
	bytes, _ := json.Marshal(r)
	return string(bytes)
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
