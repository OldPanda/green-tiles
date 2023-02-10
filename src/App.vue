<template>
  <div class="container mx-auto font-mono flex flex-col h-screen justify-between">
    <Header />
    <Hero class="mb-auto" />
    <Contributions class="mb-auto" />
    <Footer />
  </div>
</template>

<script setup lang="ts">
import Header from './components/Header.vue';
import Hero from './components/Hero.vue';
import Contributions from './components/Contributions.vue';
import Footer from './components/Footer.vue';
</script>
