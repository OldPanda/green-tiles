package services

type GitHubResponse struct {
	Data    GitHubData    `json:"data"`
	Errors  []GitHubError `json:"errors,omitempty"`
	Message string        `json:"message,omitempty"`
}

type GitHubData struct {
	User GitHubUser `json:"user"`
}

type GitHubUser struct {
	Name                    string                  `json:"name"`
	Login                   string                  `json:"login"`
	AvatarURL               string                  `json:"avatarUrl"`
	ContributionsCollection ContributionsCollection `json:"contributionsCollection"`
}

type ContributionsCollection struct {
	Years                []int                `json:"years,omitempty"`
	ContributionCalendar ContributionCalendar `json:"contributionCalendar,omitempty"`
}

type ContributionCalendar struct {
	Total int    `json:"total"`
	Weeks []Week `json:"weeks"`
}

type Week struct {
	Days []Day `json:"days"`
}

type Day struct {
	Level   string `json:"level"`
	Weekday int    `json:"weekday"`
}

type GitHubError struct {
	Message string `json:"message"`
}
