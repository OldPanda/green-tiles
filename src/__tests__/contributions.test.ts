import { createPinia, setActivePinia } from 'pinia'
import { enableAutoUnmount, mount } from '@vue/test-utils'
import { afterEach, beforeEach, expect, it } from 'vitest'
import Contributions from '@/components/Contributions.vue'
import { contributionsStore } from '@/stores/contributions'
import { contributions } from './fixtures'

enableAutoUnmount(afterEach)
let pinia: ReturnType<typeof createPinia>
beforeEach(() => {
  pinia = createPinia()
  setActivePinia(pinia)
})
const render = () => mount(Contributions, { global: { plugins: [pinia] } })

it('stays hidden before data is loaded and reacts when it arrives', async () => {
  const wrapper = render()
  expect(wrapper.find('svg').exists()).toBe(false)
  contributionsStore().setContributions(contributions())
  await wrapper.vm.$nextTick()
  expect(wrapper.findAll('svg')).toHaveLength(2)
  expect(wrapper.get('a').attributes('href')).toBe('https://github.com/OldPanda')
  expect(wrapper.get('img').attributes('src')).toBe(contributions().avatarUrl)
  expect(wrapper.get('a').text()).toContain('OldPanda')
  expect(wrapper.findAll('p').map(p => p.text())).toEqual(['10 contributions in 2025', '0 contributions in 2024'])
})

it('renders contribution levels, day positions and singular/plural tooltips', () => {
  contributionsStore().setContributions(contributions())
  const wrapper = render()
  const tiles = wrapper.findAll('rect')
  expect(tiles.map(tile => tile.attributes('class'))).toEqual([
    'NONE', 'FIRST_QUARTILE', 'SECOND_QUARTILE', 'THIRD_QUARTILE', 'FOURTH_QUARTILE'
  ])
  expect(tiles[0].attributes('y')).toBe('39')
  expect(tiles[1].attributes('y')).toBe('52')
  expect(tiles[2].attributes('x')).toBe('24')
  expect(wrapper.findAll('title').map(title => title.text())).toEqual([
    'No contributions on 2025-01-01', '1 contribution on 2025-01-02',
    '2 contributions on 2025-01-05', '3 contributions on 2025-02-02', '4 contributions on 2025-03-02'
  ])
})

it('labels each month once and preserves positions across empty weeks', () => {
  contributionsStore().setContributions(contributions())
  const wrapper = render()
  const months = wrapper.findAll('svg')[0].findAll('text').filter(text => text.attributes('y') === '-7')
  expect(months.map(month => [month.text(), month.attributes('x')])).toEqual([
    ['Jan', '25'], ['Feb', '51'], ['Mar', '77']
  ])
})
