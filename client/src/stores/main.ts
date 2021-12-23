import {defineStore} from 'pinia'
import {nanoid} from "nanoid";
import {useRoomStore} from "~/stores/room";

export const useMainStore = defineStore('main', {
    state: () => ({
        id: nanoid(10),
        name: 'User',
        isHost: false,
        conn: <WebSocket>{},
        counter: 0,
        forceMuteVideo: false,
        forceMuteAudio: false,
        op: {
            video: false,
            audio: false,
            screen: false,
            hand: false,
        },
        selectTrack: 0,
    }),
    getters: {
        // getters receive the state as first parameter
        doubleCount: (state) => state.counter * 2,
        // use getters in other getters
        doubleCountPlusOne(): number {
            return this.doubleCount * 2 + 1
        },
    },
    actions: {
        reset() {
            this.counter = 0
        },
        createWsConn() {
            const room = useRoomStore()
            this.conn = new WebSocket(
                `ws://localhost:8080/join?id=${room.roomId}&name=${this.name}&uid=${this.id}&host=${this.isHost}`)
            this.conn.onclose = () => {
            }
            this.conn.onmessage = (event) => {
                // console.log(event)
                const msg = JSON.parse(event.data)
                console.log(msg)
                if (msg.type == 'text') {
                    if (msg.author == this.id)
                        msg.author = 'me'
                    room.messageList = [...room.messageList, msg]
                } else if (msg.type == 'file') {
                    if (msg.author == this.id)
                        msg.author = 'me'
                    room.messageList = [...room.messageList, msg]
                } else if (msg.type == 'up') {
                    room.fetchParticipants()
                } else if (msg.type == 'mute video' && this.id == msg.data) {
                    if (this.op.video)
                        this.forceMuteVideo = true
                } else if (msg.type == 'mute audio' && this.id == msg.data) {
                    if (this.op.audio)
                        this.forceMuteAudio = true
                } else {
                    const i = room.participants.findIndex((e) => e.id == msg.author)
                    if (msg.type == 'hand') {
                        room.participants[i].hand = msg.data
                    } else if (msg.type == 'video') {
                        room.participants[i].video = msg.data
                    } else if (msg.type == 'audio') {
                        room.participants[i].audio = msg.data
                    }
                    // console.log(room.participants.values)
                    this.conn.send('{"type":"up"}')
                }
            }
        }
    },
})
