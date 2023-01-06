<template>
  <div class="hero min-h-768 bg-base-0 py-20">
    <div class="hero-content text-center">
      <div class="max-w-screen-lg">
        <p class="text-5xl font-bold">Review the contributions you have made on GitHub over the years.</p>
        <br />
        <div class="input-group flex justify-center py-20">
          <span>https://github.com/</span>
          <input required type="text" placeholder="Your GitHub Username" class="input input-bordered"
            :class="{ 'input-error': showAlert }" v-model="username" />
          <button class="btn btn-primary" @click="generate">Generate</button>
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
        <!-- *** -->
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";

const username = ref("");
let errorMsg = ref("");
let showAlert = ref(false);

function generate() {
  if (username.value === "") {
    showAlert.value = true;
    errorMsg.value = "GitHub username cannot be empty!"
    return;
  }

  showAlert.value = false;
}

function dismissAlert() {
  showAlert.value = false;
}
</script>
