package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

const (
	github_api                        = "https://api.github.com/graphql"
	content_type_key                  = "Content-Type"
	content_type_value                = "application/json"
	authorization_key                 = "authorization"
	authorization_value_template      = "Bearer %s"
	date_iso_format                   = "2006-01-02T15:04:05Z"
	contribution_years_request_body   = `{"query": "{ user(login: \"%s\") { name login avatarUrl contributionsCollection { years: contributionYears } } }"}`
	contribution_details_request_body = `{"query": "{ user(login: \"%s\") { contributionsCollection(from: \"%s\", to: \"%s\") { contributionCalendar { total: totalContributions weeks { days: contributionDays { level: contributionLevel weekday contributionCount date } } } } } }"}`
)

var github_oauth_token string
var client http.Client

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	github_oauth_token = os.Getenv("GITHUB_OAUTH_TOKEN")
	client = http.Client{
		Timeout: 5 * time.Second,
	}
}

type ContributionResult struct {
	Login                 string     `json:"login"`
	AvatarURL             string     `json:"avatarUrl"`
	Years                 []int      `json:"years"`
	ContributionCalendars []Calendar `json:"calendars"`
}

type Calendar struct {
	Year  int    `json:"year"`
	Total int    `json:"total"`
	Weeks []Week `json:"weeks"`
}

type ContributionInAYear struct{}

func GetAllContributions(username string) (*ContributionResult, error) {
	if github_oauth_token == "" {
		return nil, errors.New("GitHub OAuth token is missing")
	}

	// 1. Get contribution years
	requestBody := fmt.Sprintf(contribution_years_request_body, username)
	resp, err := callGitHubApi(requestBody)
	if err != nil {
		return nil, err
	}

	if resp.Errors != nil {
		return nil, errors.New(resp.Errors[0].Message)
	}

	response := &ContributionResult{
		Login:                 resp.Data.User.Login,
		AvatarURL:             resp.Data.User.AvatarURL,
		Years:                 resp.Data.User.ContributionsCollection.Years,
		ContributionCalendars: []Calendar{},
	}

	// 2. Get contributions in each year
	response.ContributionCalendars = make([]Calendar, len(resp.Data.User.ContributionsCollection.Years))
	type CalendarData struct {
		calendar Calendar
		idx      int
		err      error
	}
	calendarChan := make(chan CalendarData, len(resp.Data.User.ContributionsCollection.Years))
	var wg sync.WaitGroup
	for idx, year := range resp.Data.User.ContributionsCollection.Years {
		wg.Add(1)
		go func(idx int, year int) {
			defer wg.Done()
			r, err := GetContributionsDetailsInAYear(username, year)
			if err != nil {
				calendarChan <- CalendarData{
					calendar: Calendar{},
					idx:      idx,
					err:      err,
				}
				return
			}
			calendar := Calendar{
				Year:  year,
				Total: r.Total,
				Weeks: r.Weeks,
			}
			calendarChan <- CalendarData{
				calendar: calendar,
				idx:      idx,
				err:      nil,
			}
		}(idx, year)
	}
	wg.Wait()
	close(calendarChan)

	for data := range calendarChan {
		if data.err != nil {
			return nil, data.err
		}
		response.ContributionCalendars[data.idx] = data.calendar
	}

	return response, nil
}

func GetContributionsDetailsInAYear(username string, year int) (*ContributionCalendar, error) {
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).Format(date_iso_format)
	endDate := time.Date(year, 12, 31, 0, 0, 0, 0, time.UTC).Format(date_iso_format)
	requestBody := fmt.Sprintf(contribution_details_request_body, username, startDate, endDate)

	resp, err := callGitHubApi(requestBody)
	if err != nil {
		return nil, err
	}

	if resp.Errors != nil {
		return nil, errors.New(resp.Errors[0].Message)
	}

	return &resp.Data.User.ContributionsCollection.ContributionCalendar, nil
}

func callGitHubApi(requestBody string) (*GitHubResponse, error) {
	req, err := http.NewRequest(http.MethodPost, github_api, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set(content_type_key, content_type_value)
	req.Header.Set(authorization_key, fmt.Sprintf(authorization_value_template, github_oauth_token))

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	respJson := &GitHubResponse{}
	rawData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rawData, &respJson); err != nil {
		return nil, err
	}
	if respJson.Message != "" {
		// got error message
		return nil, fmt.Errorf("got message: %s", respJson.Message)
	}

	return respJson, nil
}
