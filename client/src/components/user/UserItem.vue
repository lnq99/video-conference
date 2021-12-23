<script lang="ts" setup>
import {useMainStore} from "~/stores/main";
import {useRoomStore} from '~/stores/room'
import {storeToRefs} from "pinia";

const props = defineProps<{
  id: number
}>()

const {id, isHost, conn} = storeToRefs(useMainStore())
const room = useRoomStore()
const {participants} = storeToRefs(room)

const i = participants.value.findIndex((e) => e.id == props.id)

function muteVideo() {
  if (isHost.value && participants.value[i].video) {
    const msg = {
      type: 'mute video',
      data: participants.value[i].id
    }
    conn.value.send(JSON.stringify(msg))
  }
}

function muteAudio() {
  if (isHost.value && participants.value[i].audio) {
    const msg = {
      type: 'mute audio',
      data: participants.value[i].id
    }
    conn.value.send(JSON.stringify(msg))
  }
}

// onBeforeMount(async () => {
//   if (!room.participants || room.participants.length == 0) {
//     await room.fetchParticipants()
//   }
// })
</script>

<template>
  <div :class="participants[i].hand ? 'raise-hand': ''" class="card p-3 mb-1">
    <div flex-1>{{ participants[i].name }}</div>

    <div v-if="participants[i].video" i-carbon-video mr-3 @click="muteVideo"></div>
    <div v-else i-carbon-video-off mr-3></div>
    <div v-if="participants[i].audio" i-carbon-volume-up mr-3 @click="muteAudio"></div>
    <div v-else i-carbon-volume-mute mr-3></div>
  </div>
</template>

<style scoped>
.raise-hand {
  border: lightgoldenrodyellow solid 1px;
}
</style>