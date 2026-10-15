import type { GitHubContributions, GitHubDay } from '@/types'

export function day(date: string, contributionCount: number, weekday = 0, level = 'NONE'): GitHubDay {
  return { date, contributionCount, weekday, level }
}

export function contributions(): GitHubContributions {
  return {
    login: 'OldPanda',
    avatarUrl: 'https://example.com/avatar.png',
    years: [2025, 2024],
    calendars: [
      {
        year: 2025,
        total: 10,
        weeks: [
          { days: [day('2025-01-01', 0, 3), day('2025-01-02', 1, 4, 'FIRST_QUARTILE')] },
          { days: [day('2025-01-05', 2, 0, 'SECOND_QUARTILE')] },
          { days: [day('2025-02-02', 3, 0, 'THIRD_QUARTILE')] },
          { days: [] },
          { days: [day('2025-03-02', 4, 0, 'FOURTH_QUARTILE')] }
        ]
      },
      { year: 2024, total: 0, weeks: [] }
    ]
  }
}
