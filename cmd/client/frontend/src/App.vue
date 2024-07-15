<script setup>
import { RouterView } from "vue-router";
import { IsConnectedToServer } from "../wailsjs/go/main/App";
import { onMounted, ref } from "vue";
import { EventsOn } from "../wailsjs/runtime/runtime";

const connected = ref(true);

onMounted(async () => {
  connected.value = await IsConnectedToServer();
});

EventsOn("server:no_connection", () => {
  connected.value = false;
});
</script>

<template>
  <main>
    <div class="connection-warning" v-if="!connected">
      No connection to server...
    </div>
    <RouterView />
  </main>
</template>

<style scoped>
.connection-warning {
  width: 100%;
  text-align: center;

  position: fixed;
  top: 0;
  left: 0;
  padding: 30px;
  font-size: 30px;
}
</style>
