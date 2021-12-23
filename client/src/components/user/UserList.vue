<script lang="ts" setup>
import {useRoomStore} from '~/stores/room'
import {storeToRefs} from "pinia";

const room = useRoomStore()
const {participants} = storeToRefs(room)

// onBeforeMount(async () => {
//   if (!room.participants || room.participants.length == 0)
//     await room.fetchParticipants()
// })
</script>

<template>
  <div id="list">
    <UserItem v-for="user of participants" :id="user.id"/>
  </div>
</template>

<style scoped>
#list {
  display: flex;
  flex-direction: column;
}
</style>