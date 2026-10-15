import * as Vue from 'vue'
import { flushPromises } from '@vue/test-utils'
import { expect, it, vi } from 'vitest'
import { inject } from '@vercel/analytics'
import { contributions } from './fixtures'

vi.mock('@vercel/analytics', () => ({ inject: vi.fn() }))
vi.mock('vue', async importOriginal => {
  const actual = await importOriginal<typeof Vue>()
  return { ...actual, createApp: vi.fn(actual.createApp) }
})

it('boots the real app and displays a successful contribution search', async () => {
  document.body.innerHTML = '<div id="app"></div>'
  const fetchMock = vi.fn().mockResolvedValue(
    new Response(JSON.stringify({ data: JSON.stringify(contributions()) }))
  )
  vi.stubGlobal('fetch', fetchMock)
  await import('@/main')
  const app = vi.mocked(Vue.createApp).mock.results[0].value as Vue.App
  try {
    expect(document.querySelector('.navbar')?.textContent).toContain('Green Tiles')
    expect(document.querySelector('.navbar a[target="_blank"]')?.getAttribute('href'))
      .toBe('https://github.com/OldPanda/green-tiles')
    expect(document.querySelector('footer')?.textContent).toContain('Made with')
    expect(document.querySelector('footer a')?.getAttribute('href')).toBe('https://old-panda.com')
    expect(inject).toHaveBeenCalledOnce()
    const input = document.querySelector('input')!
    input.value = 'OldPanda'
    input.dispatchEvent(new Event('input', { bubbles: true }))
    document.querySelector<HTMLButtonElement>('.join button')!.click()
    await flushPromises()
    expect(fetchMock).toHaveBeenCalledOnce()
    expect(document.querySelector('a[href="https://github.com/OldPanda"]')).not.toBeNull()
    expect(document.querySelectorAll('rect')).toHaveLength(5)
  } finally {
    app.unmount()
    document.body.innerHTML = ''
  }
})
