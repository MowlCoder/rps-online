import { defineStore } from "pinia";
import { readonly, ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime";

export const useRoomStore = defineStore("room", () => {
  const rooms = ref([]);

  EventsOn("server:room_created", (payload) => {
    rooms.value.push(payload.Room);
  });

  EventsOn("server:room_deleted", (payload) => {
    rooms.value = rooms.value.filter((room) => room.id === payload.RoomID);
  });

  function setRooms(value) {
    rooms.value = value;
  }

  return {
    rooms: readonly(rooms),

    setRooms,
  };
});
