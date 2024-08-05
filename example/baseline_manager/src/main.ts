/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'

// Components
import App from './App.vue'

// Composables
import { createApp } from 'vue'

const app = createApp(App)

import { createPinia } from 'pinia'

const pinia = createPinia()
app.use(pinia)

import 'highlight.js/styles/stackoverflow-light.css';
import hljs from 'highlight.js/lib/core';
import yaml from 'highlight.js/lib/languages/yaml';
import hljsVuePlugin from "@highlightjs/vue-plugin";

hljs.registerLanguage('yaml', yaml);
app.use(hljsVuePlugin)

registerPlugins(app)

app.mount('#app')
