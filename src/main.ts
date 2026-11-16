import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';
import './assets/main.css';
import { useLanguageStore } from './stores/languageStore';

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);

// Rendre la fonction $t universelle dans tous les templates de l'application
app.config.globalProperties.$t = (text: string) => {
  const languageStore = useLanguageStore();
  return languageStore.t(text);
};

app.mount('#app');
