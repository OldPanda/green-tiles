import { ref } from 'vue';
import { defineStore } from 'pinia';
import type { GitHubContributions } from '@/types';

export const contributionsStore = defineStore("contributions", () => {
  const contributions = ref(<GitHubContributions>{});

  function hasData(): Boolean {
    return contributions.value.login !== undefined && contributions.value.login !== "";
  }

  function setContributions(data: GitHubContributions) {
    contributions.value = data;
  }

  function getContributions(): GitHubContributions {
    return <GitHubContributions>contributions.value;
  }

  return { contributions, hasData, setContributions, getContributions };
});
