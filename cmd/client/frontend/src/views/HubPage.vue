<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { CreateRoom, JoinRoom } from "../../wailsjs/go/main/App";
import { useUserStore } from "../stores/user.store";
import { useRoomStore } from "../stores/room.store";
import { useMatchStore } from "../stores/match.store";
import { EventsOn } from "../../wailsjs/runtime/runtime";

const userStore = useUserStore();
const roomsStore = useRoomStore();
const matchStore = useMatchStore();
const router = useRouter();
const roomName = ref("");

function logout() {
  router.push("/");
}

function createRoom() {
  if (roomName.value.length === 0) {
    return;
  }

  CreateRoom(roomName.value);
  roomName.value = "";
}

function joinRoom(roomId) {
  JoinRoom(roomId);
}

EventsOn("server:joined_room_success", (payload) => {
  matchStore.setRoom(payload.Room);
  router.push("/match");
});

EventsOn("server:room_created", (payload) => {
  if (payload.Room.Creator.Username === userStore.username) {
    matchStore.setRoom(payload.Room);
    router.push("/match");
  }
});
</script>

<template>
  <div class="page-wrapper">
    <header class="header">
      <p>
        Hello, <b>{{ userStore.username }}</b
        >!
      </p>
      <button class="button" @click="logout">Log out</button>
    </header>
    <div class="rooms">
      <h2>Rooms Hub</h2>
      <div class="rooms-list">
        <div
          class="rooms-list__row"
          v-for="room in roomsStore.rooms"
          :key="room.ID"
        >
          <p>
            {{ room.Name }}
            <small>Creator: {{ room.Creator.Username }}</small>
          </p>
          <div class="rooms-list__row__right">
            <small>{{ room.Opponent ? 2 : 1 }}/2</small>
            <button class="button" @click="joinRoom(room.ID)">Join</button>
          </div>
        </div>
      </div>
    </div>
    <form class="create-room" @submit.prevent="createRoom">
      <input
        class="input"
        type="text"
        placeholder="Room name"
        v-model="roomName"
      />
      <button class="button" type="submit">Create Room</button>
    </form>
  </div>
</template>

<style scoped>
.page-wrapper {
  height: 100dvh;
  max-width: 80dvw;
  margin: 0 auto;

  padding: 40px 0;

  display: flex;
  flex-direction: column;
  gap: 20px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header p {
  font-size: 24px;
}

.rooms {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  overflow: auto;
}

.rooms-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  overflow: auto;
}

.rooms-list__row {
  display: flex;
  align-items: center;
  justify-content: space-between;

  background-color: #2b2a4c;
  border-radius: 6px;
  padding: 10px;
}

.rooms-list__row p {
  font-size: 20px;
}

.rooms-list__row__right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.rooms-list__row__right small {
  font-size: 18px;
}

.create-room {
  display: flex;
  justify-content: space-between;
  gap: 20px;
}

.create-room input {
  flex: 1;
}
</style>
