<script lang="ts" setup>
import {storeToRefs} from "pinia";
import {useMainStore} from '~/stores/main'
import {useRoomStore} from "~/stores/room";
import {Client, LocalStream, RemoteStream} from 'ion-sdk-js'
import {IonSFUJSONRPCSignal} from 'ion-sdk-js/lib/signal/json-rpc-impl'

const router = useRouter()
const {id, name, isHost, op, conn, forceMuteVideo, forceMuteAudio} = storeToRefs(useMainStore())
const room = useRoomStore()
const {roomId} = storeToRefs(room)

if (!id.value || !roomId.value) {
  router.push('/')
}

let localVideoStream: any
let localScreenStream: any
let streams: Object = {};


// console.log(router.currentRoute.value.params)
// roomId.value = <string>router.currentRoute.value.params.id

const config = {
  iceServers: [
    {urls: "stun:stun.l.google.com:19302"},
  ],
}
const subVideo = ref<HTMLDivElement>()

let client: Client

const signal = new IonSFUJSONRPCSignal("ws://localhost:7000/ws")
client = new Client(signal, config)
console.log(roomId, id)
signal.onopen = () => client.join(roomId.value, id.value)

client.ontrack = (track: MediaStreamTrack, stream: RemoteStream) => {
  // console.log("got track: ", track.id, "for stream: ", stream.id)
  // console.log(track)
  // console.log(stream)
  const videoContainer = subVideo.value
  const videoEl = document.createElement('video')
  if (track.kind === 'video') {
    track.onunmute = () => {
      videoEl.srcObject = stream
      videoEl.autoplay = true
      videoEl.controls = true
      videoEl.muted = false
      videoEl.classList.add('card')
      videoContainer.appendChild(videoEl)

      streams[stream.id] = stream

      pubStream(stream, 'big-video')

      // when the publisher leave
      stream.onremovetrack = () => {
        videoContainer.removeChild(videoEl)
        delete streams[stream.id]
      }
    }
  }
}

function pubStream(stream, id = 'local-video') {
  const videoEl = <HTMLVideoElement>document.getElementById(id)
  videoEl.autoplay = true
  videoEl.controls = false
  videoEl.muted = true
  videoEl.srcObject = stream
}

async function startPublish(isCam: boolean) {
  if (isCam) {
    await LocalStream.getUserMedia({
      resolution: "vga",
      video: true,
      audio: true,
      codec: "vp8"
    }).then((stream) => {
      localVideoStream = stream
      pubStream(stream)
      client.publish(stream)
      // console.log(client)
    }).catch(console.error)
  } else {
    await LocalStream.getDisplayMedia({
      resolution: "vga",
      video: true,
      audio: false,
      codec: "vp8"
    }).then((stream) => {
      localScreenStream = stream
      client.publish(stream)
      op.value.screen = !op.value.screen
    }).catch(console.error)
  }
}

watch(forceMuteVideo, (cur, prev) => {
  if (cur && op.value.video) {
    onToggleVideo()
    forceMuteVideo.value = false
  }
})

watch(forceMuteAudio, (cur, prev) => {
  if (cur && op.value.audio) {
    onToggleMic()
    forceMuteAudio.value = false
  }
})

function onToggleMic() {
  if (!localVideoStream) return
  op.value.audio = !op.value.audio
  if (op.value.audio)
    localVideoStream.unmute('audio')
  else
    localVideoStream.mute('audio')

  const msg = {
    type: 'audio',
    author: id.value,
    data: op.value.audio
  }
  conn.value.send(JSON.stringify(msg))
}

function onToggleVideo() {
  if (!localVideoStream) startPublish(true)
  op.value.video = !op.value.video
  if (op.value.video)
    localVideoStream.unmute('video')
  else
    localVideoStream.mute('video')

  const msg = {
    type: 'video',
    author: id.value,
    data: op.value.video
  }
  conn.value.send(JSON.stringify(msg))
}

function onToggleScreen() {
  if (!op.value.screen)
    startPublish(false)
  else if (localScreenStream) {
    localScreenStream.unpublish()
    localScreenStream = null
  }
}

function onPowerOff() {
  if (confirm('are you sure to exit?')) {
    client.close()
    router.push('/')
  }
}

function raiseHand() {
  op.value.hand = !op.value.hand
  const msg = {
    type: 'hand',
    author: id.value,
    data: op.value.hand
  }
  conn.value.send(JSON.stringify(msg))
}

function copyRoomId() {
  navigator.clipboard.writeText(roomId.value);
}

// onBeforeMount(async () => {
//   if (!room.participants || room.participants.length == 0)
//     await room.fetchParticipants()
// })

onMounted(() => {
  if (!localVideoStream) {
    startPublish(true).then(() => {
      if (!op.value.video) localVideoStream.mute('video')
      if (!op.value.audio) localVideoStream.mute('audio')
      conn.value.send(JSON.stringify({type: 'up'}))
    })
  }
})
</script>

<template>
  <div id="main-layout">

    <div id="main">
      <div id="list-video">
        <video id="local-video" autoplay class="card video" muted></video>
        <h3> Remote Video </h3>
        <div ref="subVideo"></div>
      </div>

      <div id="main-video" class="m-2 mt-0">
        <video id="big-video" autoplay class="card" muted></video>
      </div>

    </div>
    <div id="right-panel">
      <p class="border" @click="copyRoomId">
        Room ID:<br> {{ roomId }}
      </p>
      <br>
      <div i-carbon-pedestrian inline-block text-4xl/>
      <p>
        Hi, {{ name }}
      </p>
      <br/>
      <UserList/>
    </div>
  </div>
  <Chat></Chat>
  <div>
    <button
        class="btn btn-ctrl"
        @click=onToggleMic>
      <template v-if="op.audio">
        <i i-carbon-volume-mute>ico</i> Mute
      </template>
      <template v-else>
        <i i-carbon-volume-up>ico</i> Unmute
      </template>
    </button>
    <button
        class="btn btn-ctrl"
        @click=onToggleVideo>
      <template v-if="op.video">
        <i i-carbon-video-off>ico</i> Stop Video
      </template>
      <template v-else>
        <i i-carbon-video>ico</i> Start Video
      </template>
    </button>
    <button
        class="btn btn-ctrl"
        @click=onToggleScreen>
      <template v-if="op.screen">
        <i i-carbon-screen-off>ico</i> Stop Screen
      </template>
      <template v-else>
        <i i-carbon-screen>ico</i> Start Screen
      </template>
    </button>
    <button
        class="btn btn-ctrl"
        @click=raiseHand>
      <i i-carbon-voice-activate>ico</i>
      <template v-if="op.hand">
        Lower Hand
      </template>
      <template v-else>
        Raise Hand
      </template>
    </button>
    <button
        class="btn btn-ctrl bg-red-800 hover:bg-red-900"
        @click=onPowerOff>
      <i i-carbon-power>ico</i> Disconnect
    </button>

  </div>
</template>


<style scoped>
#main-layout {
  display: flex;
  flex: 1;
  width: 100%;
}

.btn-ctrl {
  height: 40px;
  width: 120px;
  border-radius: 20px;
  padding: 0;
  margin: 0 5px;
}

.btn-ctrl {
  font-size: small;
}

.btn-ctrl > i {
  font-size: medium;
}

.btn-ctrl:hover {
  /*background-color: steelblue;*/
}


#main {
  flex: 1;
  display: flex;
}

#main-video {
  flex: 1;
}

#list-video {
  width: 250px;
  display: flex;
  flex-direction: column;
  height: 90vh;
  overflow-x: hidden;
  overflow-y: scroll;
}

video {
  margin-right: 8px;
  margin-bottom: 8px;
  width: 100%;
}

#big-video {
  height: 100%;
  width: 100%;
}

#right-panel {
  width: 240px;
}
</style>
