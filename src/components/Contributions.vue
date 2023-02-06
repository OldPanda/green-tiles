<template>
  <div v-if="contributions.hasData()">
    <a :href="GitHubProfilePrefix + contributions.getContributions().login" target="_blank">
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
import { computed } from "vue";

import { contributionsStore } from "@/stores/contributions";
import { GitHubProfilePrefix } from "@/constants";

let contributions = contributionsStore();
</script>
