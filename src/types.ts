// export interface FetchDataResponse {
//   contributions: GitHubContributions,
//   error: boolean,
//   message: string
// }

// response body returned by backend API
export interface APIResponse {
  error: string,
  data: string
}

export interface GitHubContributions {
  login: string,
  avatarUrl: string,
  years: number[],
  calendars: GitHubCalendar[]
}

export interface GitHubCalendar {
  year: number,
  total: number,
  weeks: GitHubWeek[]
}

export interface GitHubWeek {
  days: GitHubDay[]
}

export interface GitHubDay {
  level: string,
  weekday: number,
  contributionCount: number,
  date: string
}

export interface MonthPosPair {
  month: string,
  xPos: number
}
