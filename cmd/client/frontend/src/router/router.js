import { createRouter, createWebHashHistory } from "vue-router";
import LoginPage from "../views/LoginPage.vue";
import HubPage from "../views/HubPage.vue";
import MatchPage from "../views/MatchPage.vue";

const routes = [
  {
    path: "/hub",
    component: HubPage,
  },
  {
    path: "/match",
    component: MatchPage,
  },
  {
    path: "/",
    component: LoginPage,
  },
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
});
