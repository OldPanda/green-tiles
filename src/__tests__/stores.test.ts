import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'
import { contributionsStore } from '@/stores/contributions'
import { loadingStatusStore } from '@/stores/loading'
import { contributions } from './fixtures'

beforeEach(() => setActivePinia(createPinia()))

describe('contributions store', () => {
  it('starts empty and accepts and replaces contribution data', () => {
    const store = contributionsStore()
    expect(store.hasData()).toBe(false)
    expect(store.getContributions()).toEqual({})
    store.setContributions(contributions())
    expect(store.hasData()).toBe(true)
    expect(store.getContributions()).toEqual(contributions())
    const replacement = { ...contributions(), login: 'another-user', calendars: [] }
    store.setContributions(replacement)
    expect(store.getContributions()).toEqual(replacement)
  })

  it('treats an empty login as missing data', () => {
    const store = contributionsStore()
    store.setContributions({ ...contributions(), login: '' })
    expect(store.hasData()).toBe(false)
  })
})

it('toggles loading on and off', () => {
  const store = loadingStatusStore()
  expect(store.isLoading).toBe(false)
  store.setTrue()
  expect(store.isLoading).toBe(true)
  store.setFalse()
  expect(store.isLoading).toBe(false)
})
