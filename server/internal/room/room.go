package room

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

type J map[string]interface{}

var (
	Rooms RoomMap
	//upgrader = websocket.Upgrader{
	//	CheckOrigin: func(r *http.Request) bool {
	//		return true
	//	}}
)

type Participant struct {
	Id     string          `json:"id"`
	Name   string          `json:"name"`
	IsHost bool            `json:"-"`
	Conn   *websocket.Conn `json:"-"`
	Video  bool            `json:"video"`
	Audio  bool            `json:"audio"`
	Hand   bool            `json:"hand"`
}

type RoomMap struct {
	Mutex  sync.RWMutex
	HostId string
	Map    map[string][]Participant
}

func randString(n int) string {
	rand.Seed(time.Now().UnixNano())
	l := len(Letters)

	b := make([]byte, n)
	for i := range b {
		b[i] = Letters[rand.Intn(l)]
	}

	return string(b)
}

func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

func (r *RoomMap) Get(roomId string) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	log.Println(r.Map[roomId])

	return r.Map[roomId]
}

func (r *RoomMap) CreateRoom() string {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	roomId := randString(8)
	r.Map[roomId] = []Participant{}

	return roomId
}

func (r *RoomMap) InsertIntoRoom(roomID, id, name string, host bool, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{id, name, host, conn, false, false, false}

	log.Println("Inserting into Room with RoomID: ", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)
}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomID)
}

//type resp struct {
//	RoomID string `json:"room_id"`
//	UserID string `json:"user_id"`
//}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Create")
	roomID := Rooms.CreateRoom()
	json.NewEncoder(w).Encode(J{"room_id": roomID})
}

func JoinRoom(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Join")
	roomID, ok := r.URL.Query()["id"]
	if !ok {
		log.Println("roomID missing")
	}

	name, ok := r.URL.Query()["name"]
	if !ok {
		log.Println("username missing")
	}

	uid, ok := r.URL.Query()["uid"]
	if !ok {
		log.Println("userid missing")
	}

	isHostQuery, ok := r.URL.Query()["host"]
	isHost := false
	if ok {
		isHost, _ = strconv.ParseBool(isHostQuery[0])
	}

	//uid := randString(8)
	//res := resp{roomID[0], uid}
	//json.NewEncoder(w).Encode(res)
	Rooms.InsertIntoRoom(roomID[0], uid[0], name[0], isHost, nil)
	log.Println(Rooms.Map)

	//chat.ServeWs(w, r, roomID[0])

	//ws, err := upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	log.Fatal("Web socket upgrade error", err)
	//}

	//for {
	//	var msg broadcastMsg
	//
	//	err := ws.ReadJSON(&msg.Message)
	//	if err != nil {
	//		log.Fatal("Read err:", err)
	//	}
	//	msg.Client = ws
	//	msg.RoomID = roomID[0]
	//
	//	broatcast <- msg
	//}
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["id"]
	if !ok {
		log.Println("roomID missing")
	}

	room := Rooms.Get(roomID[0])

	json.NewEncoder(w).Encode(room)
}
