<script setup>
import { ref } from "vue";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { MakeChoice } from "../../wailsjs/go/main/App";
import { useMatchStore } from "../stores/match.store";
import { useUserStore } from "../stores/user.store";
import { useRouter } from "vue-router";

const actionTextMap = {
  0: "Waiting for opponent...",
  1: "Make your turn!",
  2: "Duel ended",
};

const choiceMap = {
  0: "❔",
  1: "🪨",
  2: "📄",
  3: "✂️",
};

const matchResultMap = {
  0: "You won!",
  1: "You lost!",
  2: "Draw!",
};

const router = useRouter();

const alreadyMadeChoice = ref(false);
const matchResult = ref(null);
const creatorChoice = ref("❔");
const opponentChoice = ref("❔");

const userStore = useUserStore();
const matchStore = useMatchStore();

function makeChoice(choice) {
  if (matchStore.room.Status !== 1 || alreadyMadeChoice.value) {
    return;
  }

  if (userStore.username === matchStore.room.Creator.Username) {
    creatorChoice.value = choiceMap[choice];
  } else {
    opponentChoice.value = choiceMap[choice];
  }

  alreadyMadeChoice.value = true;

  MakeChoice(choice);
}

function goToHub() {
  router.push("/hub");
}

EventsOn("server:room_joined", (payload) => {
  matchStore.setOpponent(payload.JoinedUser);
  matchStore.setStatus(1);
});

EventsOn("server:match_end", (payload) => {
  creatorChoice.value = choiceMap[payload.CreatorChoice];
  opponentChoice.value = choiceMap[payload.OpponentChoice];
  matchResult.value = payload.MatchResult;
  matchStore.setStatus(2);
});
</script>

<template>
  <div class="page-wrapper">
    <div class="action-text">
      <p>{{ actionTextMap[matchStore.room.Status] }}</p>
      <p v-if="matchStore.room.Status === 2">
        {{ matchResultMap[matchResult] }}
      </p>
    </div>

    <div class="match-zone">
      <div class="left">
        <p
          :class="{
            own: matchStore.room.Creator.Username === userStore.username,
          }"
        >
          {{ matchStore.room.Creator.Username }}
        </p>
        <div class="selected">{{ creatorChoice }}</div>
      </div>
      <div class="center">
        <p>vs</p>
        <small>bo1</small>
      </div>
      <div class="right">
        <p
          v-if="matchStore.room.Opponent"
          :class="{
            own: matchStore.room.Opponent.Username === userStore.username,
          }"
        >
          {{ matchStore.room.Opponent.Username }}
        </p>
        <p v-else>❔❔❔</p>
        <div class="selected">{{ opponentChoice }}</div>
      </div>
    </div>
    <div class="control-zone">
      <div
        v-if="matchStore.room.Status !== 2"
        class="control-zone__panel"
        :class="{ disabled: matchStore.room.Status !== 1 || alreadyMadeChoice }"
      >
        <p @click="makeChoice(1)">🪨</p>
        <p @click="makeChoice(2)">📄</p>
        <p @click="makeChoice(3)">✂️</p>
      </div>
      <div v-else>
        <button class="button" @click="goToHub">Go to Hub</button>
      </div>
    </div>
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

.action-text {
  font-size: 30px;
  text-align: center;
}

.match-zone {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex: 1;
}

.match-zone .left,
.match-zone .right {
  width: 40%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
}

.match-zone .left > p,
.match-zone .right > p {
  font-size: 40px;
  font-weight: 600;
}

.match-zone .left > p.own,
.match-zone .right > p.own {
  color: red;
}

.match-zone .center {
  text-align: center;
}

.match-zone .center p {
  font-size: 60px;
  text-transform: uppercase;
  font-weight: 700;
}

.match-zone .center small {
  font-size: 30px;
}

.match-zone .selected {
  font-size: 100px;
}

.control-zone {
  display: flex;
  justify-content: center;
}

.control-zone__panel {
  display: flex;
  gap: 20px;
}

.control-zone__panel.disabled p {
  opacity: 0.5;
  cursor: not-allowed;
}

.control-zone__panel p {
  display: flex;
  align-items: center;
  justify-content: center;

  font-size: 50px;
  padding: 30px;
  border: 2px solid #433d8b;
  border-radius: 6px;

  cursor: pointer;
  transition: 0.2s;
}

.control-zone__panel:not(.disabled) p:hover {
  background-color: #433d8b;
}
</style>
