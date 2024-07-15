import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import "./style.css";
import { router } from "./router/router";

createApp(App).use(createPinia()).use(router).mount("#app");
