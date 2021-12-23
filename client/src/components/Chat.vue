<script lang="ts" setup>
import {storeToRefs} from "pinia";
import {useMainStore} from '~/stores/main'
import {Msg, useRoomStore} from "~/stores/room";

const {id, conn} = storeToRefs(useMainStore())
const room = useRoomStore()
const {roomId, participants, messageList} = storeToRefs(room)
// const conn = new WebSocket("ws://localhost:8080/ws/" + roomId)


onBeforeMount(async () => {
  if (!room.participants || room.participants.length == 0)
    await room.fetchParticipants()
})

const data = ref({
  participants: [
    // {
    //   id: 'user1',
    //   name: 'Matteo',
    //   imageUrl: 'https://eu.ui-avatars.com/api/?background=random&name=QL'
    // },
  ], // the list of all the participant of the conversation. `name` is the user name, `id` is used to establish the author of a message, `imageUrl` is supposed to be the user avatar.
  titleImageUrl: 'https://a.slack-edge.com/66f9/img/avatars-teams/ava_0001-34.png',
  messageList: [
    // {type: 'text', author: 'me', data: {text: 'Say yes!'}},
    // {type: 'text', author: 'user1', data: {text: 'No.'}}
  ],
  newMessagesCount: 0,
  isChatOpen: false,
  alwaysScrollToBottom: true,
  messageStyling: false // enables *bold* /emph/ _underline_ and such (more info at github.com/mattezza/msgdown)
})


// conn.value.onclose = () => {
// }
// conn.value.onmessage = (event) => {
//   const message = JSON.parse(event.data)
//   if (message.author == id.value)
//     message.author = 'me'
//   // console.log(event)
//   data.value.messageList = [...data.value.messageList, message]
// }


function onMessageWasSent(message: Msg) {
  if (message.author == 'me')
    message.author = id.value
  // called when the user sends a message
  conn.value.send(JSON.stringify(message))
}

function openChat() {
  data.value.isChatOpen = true
  data.value.newMessagesCount = 0
}

function closeChat() {
  data.value.isChatOpen = false
}

const colors = {
  header: {bg: '#4e8cff', text: '#ffffff'},
  launcher: {bg: '#4e8cff'},
  messageList: {bg: '#ffffff'},
  sentMessage: {bg: '#4e8cff', text: '#ffffff'},
  receivedMessage: {bg: '#eaeaea', text: '#222222'},
  userInput: {bg: '#f4f7f9', text: '#565867'}
}
</script>

<template>
  <beautiful-chat
      :alwaysScrollToBottom="data.alwaysScrollToBottom"
      :close="closeChat"
      :colors="colors"
      :deletionConfirmation="true"
      :disableUserListToggle="false"
      :icons="data.icons"
      :isOpen="data.isChatOpen"
      :messageList="messageList"
      :messageStyling="data.messageStyling"
      :newMessagesCount="data.newMessagesCount"
      :onMessageWasSent="onMessageWasSent"
      :open="openChat"
      :participants="participants"
      :showCloseButton="true"
      :showDeletion="true"
      :showEdition="false"
      :showEmoji="true"
      :showFile="true"
      :showLauncher="true"
      :titleImageUrl="data.titleImageUrl"
  ></beautiful-chat>
</template>
