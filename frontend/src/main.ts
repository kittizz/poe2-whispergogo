/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from "@/plugins"
import { useAppStore } from "./stores/app"

// Components
import App from "./App.vue"

// Composables
import { createApp } from "vue"

const app = createApp(App)

registerPlugins(app)

const appStore = useAppStore()
appStore.startGameStatusWatcher()

app.mount("#app")
