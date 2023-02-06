import { ref } from 'vue';
import { defineStore } from 'pinia';

export const loadingStatusStore = defineStore("loading", () => {
  const isLoading = ref(false);
  function setTrue() {
    isLoading.value = true;
  }
  function setFalse() {
    isLoading.value = false;
  }

  return { isLoading, setTrue, setFalse };
});
