import { createApp } from 'vue';
import { createPinia } from 'pinia';
import '@fontsource/inter/400.css';
import '@fontsource/inter/500.css';
import '@fontsource/inter/600.css';
import '@fontsource/inter/700.css';
import '@fontsource/inter/800.css';
import '@fontsource/jetbrains-mono/400.css';
import '@fontsource/jetbrains-mono/500.css';
import '@fontsource/jetbrains-mono/600.css';
import '@fontsource/jetbrains-mono/700.css';
import App from './App.vue';
import './assets/main.css';
import './assets/themes.css';
import { useLanguageStore } from './stores/languageStore';

// Restaurer immédiatement le thème de couleur et le background personnalisés
if (typeof window !== 'undefined') {
  const savedColorTheme = localStorage.getItem('autoschool_color_theme');
  if (savedColorTheme && savedColorTheme !== 'default') {
    document.documentElement.setAttribute('data-color-theme', savedColorTheme);
  }
  const savedBgStyle = localStorage.getItem('autoschool_bg_style');
  if (savedBgStyle && savedBgStyle !== 'default') {
    document.documentElement.setAttribute('data-bg-style', savedBgStyle);
  }
}

const app = createApp(App);
const pinia = createPinia();

app.use(pinia);

// Rendre la fonction $t universelle dans tous les templates de l'application
app.config.globalProperties.$t = (text: string) => {
  const languageStore = useLanguageStore();
  return languageStore.t(text);
};

app.mount('#app');
