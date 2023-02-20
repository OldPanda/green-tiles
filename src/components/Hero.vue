<template>
  <div class="hero min-h-768 bg-base-0 pt-20">
    <div class="hero-content text-center">
      <div class="max-w-screen-lg">
        <p class="text-3xl lg:text-5xl font-bold">{{ HERO_STATEMENT }}</p>
        <br />
        <!-- Show on desktop screen -->
        <div class="hidden lg:block">
          <div class="input-group flex justify-center py-20">
            <span>{{ GITHUB_PROFILE_PREFIX }}</span>
            <input required type="text" placeholder="Your GitHub Username" class="input input-bordered"
              :class="{ 'input-error': showAlert }" v-model="username" />
            <button class="btn btn-primary loading" v-if="loadingStore.isLoading">Generating...</button>
            <button class="btn btn-primary" @click="generate" v-else>Generate</button>
          </div>
        </div>
        <!-- Show on mobile screen -->
        <div class="lg:hidden">
          <div class="input-group flex justify-center py-20 lg:hidden">
            <input required type="text" placeholder="Your GitHub Username" class="input input-bordered"
              :class="{ 'input-error': showAlert }" v-model="username" />
            <button class="btn btn-primary loading" v-if="loadingStore.isLoading">Generating...</button>
            <button class="btn btn-primary" @click="generate" v-else>Generate</button>
          </div>
        </div>
        <br />
        <!-- Alert -->
        <div class="toast">
          <div class="alert alert-error shadow-lg" v-if="showAlert">
            <div>
              <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current flex-shrink-0 h-6 w-6" fill="none"
                viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{{ errorMsg }}</span>
            </div>
            <div class="flex-none">
              <button class="btn btn-xs btn-circle btn-ghost" @click="dismissAlert">
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5"
                  stroke="currentColor" class="w-6 h-6">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        </div>
        <!-- End Alert -->
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

import { contributionsStore } from "@/stores/contributions";
import { loadingStatusStore } from "@/stores/loading";
import { fetchGitHubContributions } from "@/services/github-data-fetcher";
import { HERO_STATEMENT, GITHUB_PROFILE_PREFIX } from "@/constants";

const username = ref("");
let errorMsg = ref("");
let showAlert = ref(false);
let loadingStore = loadingStatusStore();
let contributions = contributionsStore();

function generate() {
  if (username.value === "") {
    showErrPopUp("GitHub username cannot be empty!");
    return;
  }

  if (showAlert.value) {
    dismissErrPopUp();
  }
  fetchContributions(username.value);
}

function dismissAlert() {
  showAlert.value = false;
}

async function fetchContributions(username: string) {
  loadingStore.setTrue();
  try {
    let data = await fetchGitHubContributions(username);
    contributions.setContributions(data);
  } catch (e) {
    showErrPopUp(<string>e);
  } finally {
    loadingStore.setFalse();
  }
}

function showErrPopUp(message: string) {
  if (message.length === 0) return;
  showAlert.value = true;
  errorMsg.value = message;
}

function dismissErrPopUp() {
  showAlert.value = false;
  errorMsg.value = "";
}
</script>
