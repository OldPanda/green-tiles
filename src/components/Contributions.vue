<template>
  <div class="lg:grid justify-center mx-2" v-if="contributions.hasData()">
    <!-- header -->
    <a class="grid justify-items-start mb-4" :href="GITHUB_PROFILE_PREFIX + contributions.getContributions().login"
      target="_blank">
      <div class="avatar">
        <div class="w-8 rounded-full">
          <img :src=contributions.getContributions().avatarUrl />
        </div>
        <span class="text-2xl font-semibold ml-2">
          {{ contributions.getContributions().login }}
        </span>
      </div>
    </a>
    <!-- end header -->
    <div class="grid justify-items-start" v-for="calendar in contributions.getContributions().calendars">
      <p>{{ calendar.total }} contributions in {{ calendar.year }}</p>
      <div class="container" style="overflow-x: auto;">
        <svg width="728" height="112" class="mb-10">
          <g transform="translate(15, 20)">
            <g v-for="(week, idx) in calendar.weeks" :transform="`translate(` + idx * 14 + `, 0)`">
              <rect width="10" height="10" v-for="day in week.days" :x="25 - idx" :y="day.weekday * 13" rx="2" ry="2"
                :class="day.level">
                <title>{{ tooltipStatement(day) }}</title>
              </rect>
            </g>

            <text v-for=" pair in monthPositions(calendar)" :x="pair.xPos" y="-7" class="fill-base-content">
              {{ pair.month }}
            </text>

            <text dx="-15" dy="8" style="display: none;" class="fill-base-content">Sun</text>
            <text dx="-15" dy="22" class="fill-base-content">Mon</text>
            <text dx="-15" dy="32" style="display: none;" class="fill-base-content">Tue</text>
            <text dx="-15" dy="48" class="fill-base-content">Wed</text>
            <text dx="-15" dy="57" style="display: none;" class="fill-base-content">Thu</text>
            <text dx="-15" dy="73" class="fill-base-content">Fri</text>
            <text dx="-15" dy="81" style="display: none;" class="fill-base-content">Sat</text>
          </g>
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { contributionsStore } from "@/stores/contributions";
import { GITHUB_PROFILE_PREFIX, MONTHS } from "@/constants";
import type { GitHubCalendar, GitHubDay, MonthPosPair } from "@/types";

let contributions = contributionsStore();

function monthPositions(calendar: GitHubCalendar): MonthPosPair[] {
  let pairs: MonthPosPair[] = [];
  let baseXPos = 25;
  pairs.push({
    month: "Jan",
    xPos: baseXPos
  });
  let step = 13;
  let weekCount = 0;
  let visited = new Set<number>();
  for (let week of calendar.weeks) {
    let month = +week.days[0].date.split("-")[1];
    if (visited.has(month)) {
      weekCount++;
      continue;
    }
    pairs.push({
      month: MONTHS[month - 1],
      xPos: baseXPos + step * weekCount
    });
    visited.add(month);
    weekCount++;
  }
  return pairs;
}

function tooltipStatement(day: GitHubDay): string {
  let count = day.contributionCount;
  if (count == 0) {
    return `No contributions on ${day.date}`;
  } else if (count == 1) {
    return `1 contribution on ${day.date}`;
  } else {
    return `${count} contributions on ${day.date}`;
  }
}
</script>

<style scoped>
.NONE {
  fill: var(--none);
}

.FIRST_QUARTILE {
  fill: var(--first);
}

.SECOND_QUARTILE {
  fill: var(--second);
}

.THIRD_QUARTILE {
  fill: var(--third);
}

.FOURTH_QUARTILE {
  fill: var(--fourth);
}
</style>
