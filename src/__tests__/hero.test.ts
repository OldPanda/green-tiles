import { createPinia, setActivePinia } from 'pinia'
import { enableAutoUnmount, flushPromises, mount } from '@vue/test-utils'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import Hero from '@/components/Hero.vue'
import { fetchGitHubContributions } from '@/services/github-data-fetcher'
import { contributionsStore } from '@/stores/contributions'
import { loadingStatusStore } from '@/stores/loading'
import type { GitHubContributions } from '@/types'
import { contributions } from './fixtures'

vi.mock('@/services/github-data-fetcher', () => ({ fetchGitHubContributions: vi.fn() }))
enableAutoUnmount(afterEach)
let pinia: ReturnType<typeof createPinia>
beforeEach(() => {
  pinia = createPinia()
  setActivePinia(pinia)
})
const render = () => mount(Hero, { global: { plugins: [pinia] } })

describe('contribution search', () => {
  it('validates an empty username and dismisses the alert', async () => {
    const wrapper = render()
    await wrapper.get('button').trigger('click')
    expect(wrapper.get('.alert').text()).toContain('GitHub username cannot be empty!')
    expect(wrapper.get('input').classes()).toContain('input-error')
    expect(fetchGitHubContributions).not.toHaveBeenCalled()
    await wrapper.get('.toast button').trigger('click')
    expect(wrapper.find('.alert').exists()).toBe(false)
    expect(wrapper.get('input').classes()).not.toContain('input-error')
  })

  it('shows loading until the response arrives and stores the result', async () => {
    let resolve!: (data: GitHubContributions) => void
    vi.mocked(fetchGitHubContributions).mockImplementation(() => new Promise(done => { resolve = done }))
    const wrapper = render()
    await wrapper.get('input').setValue('OldPanda')
    await wrapper.get('button').trigger('click')
    expect(fetchGitHubContributions).toHaveBeenCalledWith('OldPanda')
    expect(wrapper.get('button').text()).toContain('Generating')
    expect(wrapper.get('button').attributes('disabled')).toBeDefined()
    expect(loadingStatusStore().isLoading).toBe(true)
    resolve(contributions())
    await flushPromises()
    expect(contributionsStore().getContributions()).toEqual(contributions())
    expect(loadingStatusStore().isLoading).toBe(false)
    expect(wrapper.get('button').text()).toBe('Generate')
  })

  it.each([new Error('API unavailable'), 'API unavailable'])(
    'displays a rejected request and clears loading: %s', async error => {
      vi.mocked(fetchGitHubContributions).mockRejectedValue(error)
      const wrapper = render()
      await wrapper.get('input').setValue('OldPanda')
      await wrapper.get('button').trigger('click')
      await flushPromises()
      expect(wrapper.get('.alert').text()).toBe('API unavailable')
      expect(loadingStatusStore().isLoading).toBe(false)
      expect(contributionsStore().hasData()).toBe(false)
    }
  )

  it('clears an earlier alert when a valid search succeeds', async () => {
    const wrapper = render()
    await wrapper.get('button').trigger('click')
    vi.mocked(fetchGitHubContributions).mockResolvedValue(contributions())
    await wrapper.get('input').setValue('OldPanda')
    await wrapper.get('button').trigger('click')
    await flushPromises()
    expect(wrapper.find('.alert').exists()).toBe(false)
    expect(contributionsStore().hasData()).toBe(true)
  })

  it('does not show an empty error message and still resets loading', async () => {
    vi.mocked(fetchGitHubContributions).mockRejectedValue(new Error(''))
    const wrapper = render()
    await wrapper.get('input').setValue('OldPanda')
    await wrapper.get('button').trigger('click')
    await flushPromises()
    expect(wrapper.find('.alert').exists()).toBe(false)
    expect(loadingStatusStore().isLoading).toBe(false)
  })
})
