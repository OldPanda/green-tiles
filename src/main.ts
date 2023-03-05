import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from '@/App.vue';
import { inject } from '@vercel/analytics';

import '@/style.css';

const app = createApp(App);

app.use(createPinia());

app.mount('#app');

inject();
