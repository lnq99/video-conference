import {defineStore} from 'pinia'

export type Msg = {
    type: string,
    author: string,
    data: any,
}

export const useRoomStore = defineStore('room', {
    state: () => ({
        roomId: '',
        // participants: {
        //     1: {id: 1, name: 'User1', audio: false, video: false, raiseHand: false},
        //     2: {id: 2, name: 'User2', audio: false, video: false, raiseHand: false},
        //     3: {id: 3, name: 'User3', audio: false, video: false, raiseHand: false},
        // }
        participants: <any[]>[],
        messageList: <Msg[]>[],
    }),
    actions: {
        fetchParticipants() {
            // this.participants = 0
            fetch(`http://localhost:8080/participants?id=${this.roomId}`, {
                method: 'GET'
            })
                .then((res) => res.json())
                .then(data => {
                    console.log(data)
                    this.participants = []

                    for (let e of data) {
                        this.participants.push({
                            id: e.id,
                            name: e.name,
                            imageUrl: `https://eu.ui-avatars.com/api/?background=random&name=${e.name}`,
                            video: e.video,
                            audio: e.audio,
                            hand: e.hand,
                        })
                    }

                    console.log(this.participants)

                })
                .catch(console.error)
        },
    },
    // getters: {
    //     participants(): any[] {
    //         if (!this._participants) return []
    //         return this._participants.map((e) => Object({
    //             id: e.id,
    //             name: e.name,
    //             imageUrl: `https://eu.ui-avatars.com/api/?background=random&name=${e.name}`,
    //             video: e.video,
    //             audio: e.audio,
    //             hand: e.raiseHand,
    //         }))
    //     }
    // }
})
