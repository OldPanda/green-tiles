import { beforeEach, describe, expect, it, vi } from 'vitest'
import { fetchGitHubContributions } from '@/services/github-data-fetcher'
import { contributions } from './fixtures'

const fetchMock = vi.fn<typeof fetch>()

beforeEach(() => {
  vi.stubGlobal('fetch', fetchMock)
})

describe('fetchGitHubContributions', () => {
  it('rejects an empty username before contacting the API', async () => {
    await expect(fetchGitHubContributions('')).rejects.toThrow('Username is empty.')
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('decodes the API envelope and encodes the username as one query parameter', async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ data: JSON.stringify(contributions()) })))
    await expect(fetchGitHubContributions('user&admin=true #')).resolves.toEqual(contributions())
    expect(fetchMock).toHaveBeenCalledWith(
      'https://api.old-panda.com/github/contributions?username=user%26admin%3Dtrue%20%23'
    )
  })

  it('returns API errors to the caller', async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ error: 'User not found' }), { status: 400 }))
    await expect(fetchGitHubContributions('missing')).rejects.toThrow('User not found')
  })

  it('propagates network failures', async () => {
    fetchMock.mockRejectedValue(new Error('Network unavailable'))
    await expect(fetchGitHubContributions('OldPanda')).rejects.toThrow('Network unavailable')
  })

  it.each(['not JSON', JSON.stringify({ data: 'not JSON' })])('rejects malformed response %s', async body => {
    fetchMock.mockResolvedValue(new Response(body))
    await expect(fetchGitHubContributions('OldPanda')).rejects.toBeInstanceOf(SyntaxError)
  })
})
