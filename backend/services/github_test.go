package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func mockGitHub(t *testing.T, transport roundTripFunc) {
	t.Helper()
	previousClient, previousToken := client, github_oauth_token
	client.Transport = transport
	github_oauth_token = "test-token"
	t.Cleanup(func() { client, github_oauth_token = previousClient, previousToken })
}

func jsonResponse(body string) (*http.Response, error) {
	return &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
}

func requestQuery(t *testing.T, request *http.Request) string {
	t.Helper()
	if request.Method != http.MethodPost || request.URL.String() != github_api {
		t.Errorf("unexpected request: %s %s", request.Method, request.URL)
	}
	if request.Header.Get("Content-Type") != "application/json" || request.Header.Get("authorization") != "Bearer test-token" {
		t.Errorf("unexpected headers: %v", request.Header)
	}
	var body struct {
		Query string `json:"query"`
	}
	if err := json.NewDecoder(request.Body).Decode(&body); err != nil {
		t.Fatalf("invalid GraphQL request: %v", err)
	}
	if !strings.Contains(body.Query, `user(login: "OldPanda")`) {
		t.Errorf("missing username in query: %s", body.Query)
	}
	return body.Query
}

func TestGetAllContributionsMissingToken(t *testing.T) {
	mockGitHub(t, func(*http.Request) (*http.Response, error) {
		t.Error("called GitHub without a token")
		return nil, errors.New("unexpected request")
	})
	github_oauth_token = ""
	got, err := GetAllContributions("OldPanda")
	if got != nil || err == nil || err.Error() != "GitHub OAuth token is missing" {
		t.Fatalf("got %v, %v", got, err)
	}
}

func TestGetAllContributionsConcurrentYears(t *testing.T) {
	secondYearRequested := make(chan struct{})
	var mu sync.Mutex
	var queries []string
	mockGitHub(t, func(request *http.Request) (*http.Response, error) {
		query := requestQuery(t, request)
		mu.Lock()
		queries = append(queries, query)
		mu.Unlock()
		switch {
		case strings.Contains(query, "contributionYears"):
			return jsonResponse(`{"data":{"user":{"login":"OldPanda","avatarUrl":"avatar.png","contributionsCollection":{"years":[2025,2024]}}}}`)
		case strings.Contains(query, "2025-01-01T00:00:00Z"):
			// A sequential implementation cannot complete this request before starting 2024.
			select {
			case <-secondYearRequested:
			case <-request.Context().Done():
				return nil, request.Context().Err()
			}
			return jsonResponse(`{"data":{"user":{"contributionsCollection":{"contributionCalendar":{"total":5,"weeks":[{"days":[{"date":"2025-01-01","weekday":3,"contributionCount":5,"level":"FOURTH_QUARTILE"}]}]}}}}}`)
		default:
			close(secondYearRequested)
			return jsonResponse(`{"data":{"user":{"contributionsCollection":{"contributionCalendar":{"total":1,"weeks":[]}}}}}`)
		}
	})
	got, err := GetAllContributions("OldPanda")
	if err != nil {
		t.Fatal(err)
	}
	want := &ContributionResult{
		Login: "OldPanda", AvatarURL: "avatar.png",
		ContributionCalendars: []Calendar{
			{Year: 2025, Total: 5, Weeks: []Week{{Days: []Day{{Date: "2025-01-01", Weekday: 3, Count: 5, Level: "FOURTH_QUARTILE"}}}}},
			{Year: 2024, Total: 1, Weeks: []Week{}},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("contributions = %+v, want %+v", got, want)
	}
	if len(queries) != 3 {
		t.Errorf("requests = %d, want 3", len(queries))
	}
	for _, year := range []string{"2025", "2024"} {
		if !strings.Contains(strings.Join(queries, "\n"), year+"-12-31T00:00:00Z") {
			t.Errorf("missing end date for %s", year)
		}
	}
}

func TestGetAllContributionsNoYears(t *testing.T) {
	for _, body := range []string{
		`{"data":{"user":{"login":"OldPanda","contributionsCollection":{"years":[]}}}}`,
		`{"errors":[],"data":{"user":{"login":"OldPanda","contributionsCollection":{"years":[]}}}}`,
	} {
		t.Run(body, func(t *testing.T) {
			mockGitHub(t, func(*http.Request) (*http.Response, error) { return jsonResponse(body) })
			got, err := GetAllContributions("OldPanda")
			if err != nil || got.Login != "OldPanda" || got.ContributionCalendars == nil || len(got.ContributionCalendars) != 0 {
				t.Fatalf("got %+v, %v", got, err)
			}
		})
	}
}

func TestGetAllContributionsErrors(t *testing.T) {
	for _, test := range []struct{ name, years, details, want string }{
		{"years transport", "transport", "", "offline"},
		{"years GraphQL", `{"errors":[{"message":"user not found"}]}`, "", "user not found"},
		{"details transport", `{"data":{"user":{"contributionsCollection":{"years":[2025]}}}}`, "transport", "offline"},
		{"details GraphQL", `{"data":{"user":{"contributionsCollection":{"years":[2025]}}}}`, `{"errors":[{"message":"rate limited"}]}`, "rate limited"},
	} {
		t.Run(test.name, func(t *testing.T) {
			mockGitHub(t, func(request *http.Request) (*http.Response, error) {
				query := requestQuery(t, request)
				body := test.years
				if !strings.Contains(query, "contributionYears") {
					body = test.details
				}
				if body == "transport" {
					return nil, errors.New("offline")
				}
				return jsonResponse(body)
			})
			got, err := GetAllContributions("OldPanda")
			if got != nil || err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got %+v, %v; want %q", got, err, test.want)
			}
		})
	}
}

func TestGetContributionsDetailsEmptyErrors(t *testing.T) {
	mockGitHub(t, func(request *http.Request) (*http.Response, error) {
		requestQuery(t, request)
		return jsonResponse(`{"errors":[],"data":{"user":{"contributionsCollection":{"contributionCalendar":{"total":7,"weeks":[]}}}}}`)
	})
	got, err := GetContributionsDetailsInAYear("OldPanda", 2025)
	if err != nil || got.Total != 7 || got.Weeks == nil {
		t.Fatalf("got %+v, %v", got, err)
	}
}

type trackedBody struct {
	io.Reader
	closed bool
}

func (b *trackedBody) Close() error { b.closed = true; return nil }

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("body read failed") }

func TestCallGitHubAPIErrors(t *testing.T) {
	for _, test := range []struct {
		name string
		body io.Reader
		want string
	}{
		{"read failure", failingReader{}, "body read failed"},
		{"invalid JSON", strings.NewReader("invalid"), "invalid character"},
		{"null response", strings.NewReader("null"), "GitHub returned an empty response"},
		{"API message", strings.NewReader(`{"message":"Bad credentials"}`), "got message: Bad credentials"},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := &trackedBody{Reader: test.body}
			mockGitHub(t, func(*http.Request) (*http.Response, error) {
				return &http.Response{StatusCode: http.StatusOK, Body: body}, nil
			})
			got, err := callGitHubApi(`{"query":"{}"}`)
			if got != nil || err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("got %+v, %v; want %q", got, err, test.want)
			}
			if !body.closed {
				t.Error("response body was not closed")
			}
		})
	}
}
