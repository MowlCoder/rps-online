<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { Login } from "../../wailsjs/go/main/App";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { useUserStore } from "../stores/user.store";
import { useRoomStore } from "../stores/room.store";

const router = useRouter();
const userStore = useUserStore();
const roomStore = useRoomStore();

const username = ref("");

function onSubmit() {
  if (username.value.length === 0) {
    return;
  }

  Login(username.value);
}

EventsOn("server:success_login", (payload) => {
  userStore.setUsername(username.value);
  roomStore.setRooms(payload.Rooms);
  router.push("/hub");
});
</script>

<template>
  <div class="page-wrapper">
    <h1>🪨 📄 ✂️ Online</h1>
    <form class="form" @submit.prevent="onSubmit">
      <input
        class="input"
        type="text"
        placeholder="Username"
        v-model="username"
      />
      <button class="button" type="submit">Let's Go</button>
    </form>
  </div>
</template>

<style scoped>
.page-wrapper {
  height: 100dvh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 25px;
}

.form {
  display: flex;
  align-items: center;
  flex-direction: column;
  gap: 15px;
}
</style>
