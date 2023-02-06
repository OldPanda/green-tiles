<template>
  <div v-if="contributions.hasData()">
    <a :href="GITHUB_PROFILE_PREFIX + contributions.getContributions().login" target="_blank">
      <div class="avatar">
        <div class="w-8 rounded-full">
          <img :src=contributions.getContributions().avatarUrl />
        </div>
        <span class="text-2xl font-semibold ml-2">
          {{ contributions.getContributions().login }}
        </span>
      </div>
    </a>
  </div>
</template>

<script setup lang="ts">
import { contributionsStore } from "@/stores/contributions";
import { GITHUB_PROFILE_PREFIX } from "@/constants";

let contributions = contributionsStore();
</script>
