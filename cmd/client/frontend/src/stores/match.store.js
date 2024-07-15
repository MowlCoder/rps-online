import { defineStore } from "pinia";
import { readonly, ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime";

export const useMatchStore = defineStore("match", () => {
  const room = ref(null);

  function setRoom(value) {
    room.value = value;
  }

  function setOpponent(value) {
    room.value.Opponent = value;
  }

  function setStatus(value) {
    room.value.Status = value;
  }

  return {
    room: readonly(room),

    setRoom,
    setOpponent,
    setStatus,
  };
});
