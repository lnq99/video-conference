<script lang="ts" setup>
import {storeToRefs} from 'pinia'
import {useMainStore} from '~/stores/main'
import {useRoomStore} from "~/stores/room";

const main = useMainStore()
const room = useRoomStore()
const {id, name, isHost, conn} = storeToRefs(main)
const {roomId} = storeToRefs(room)

const router = useRouter()


function createRoom() {
  if (name.value) {
    fetch(`http://localhost:8080/create?name=${name.value}`, {
      method: 'GET'
    })
        .then((res) => res.json())
        .then(data => {
          console.log(data)
          roomId.value = data.room_id
          isHost.value = true
          joinRoom()
        })
        .catch(console.error)
  }
}

function joinRoom() {
  if (name.value && roomId.value) {
    main.createWsConn()
    room.fetchParticipants()
    router.push(`/room/${encodeURIComponent(roomId.value)}`)
  }
}
</script>

<template>
  <div p="y-10">
    <div i-carbon-video inline-block text-5xl/>
    <p>
      <a href="" rel="noreferrer" target="_blank" text-xl>
        Video Conference
      </a>
    </p>
    <p>
      <em op75>Реализация обучающей среды с возможностью проведения видеоконференций</em>
    </p>

    <div py-6/>

    <input
        id="name"
        v-model="name"
        autocomplete="false"
        bg="transparent"
        border="~ rounded gray-200 dark:gray-700"
        outline="none active:none"
        p="x-4 y-2"
        placeholder="What's your name?"
        text="center"
        type="text"
        w="250px"
    >

    <div py-2/>

    <input
        id="roomId"
        v-model="roomId"
        autocomplete="false"
        bg="transparent"
        border="~ rounded gray-200 dark:gray-700"
        outline="none active:none"
        p="x-4 y-2"
        placeholder="Room ID to join?"
        text="center"
        type="text"
        w="250px"
        @keydown.enter="joinRoom"
    >

    <div py-2/>

    <div>
      <button
          :disabled="!name"
          class="m-2 btn bg-sky-700 hover:bg-sky-800"
          @click="createRoom"
      >
        Create Room
      </button>
      <button
          :disabled="!name || !roomId"
          class="m-2 btn btn-secondary"
          @click="joinRoom"
      >
        Join Room
      </button>
    </div>
  </div>
</template>
