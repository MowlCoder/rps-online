import { defineStore } from "pinia";
import { readonly, ref } from "vue";

export const useUserStore = defineStore("user", () => {
  const username = ref("");

  function setUsername(value) {
    username.value = value;
  }

  return {
    username: readonly(username),

    setUsername,
  };
});
