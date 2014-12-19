package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/martini-contrib/render"
	"github.com/masayukioguni/go-lgtm-model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// https://github.com/patcito/martini-gorilla-websocket-chat-example/blob/master/server.go
var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(result model.Image) {
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	result.Name = S3Url + result.Name

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteJSON(result); err != nil {
			return
		}
	}
}

const (
	Dial       = "mongodb://localhost"
	DB         = "test-go-lgtm-server"
	Collection = "test_collection"
	S3Url      = "https://s3.amazonaws.com/go-lgtm/"
)

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(martini.Static("assets"))

	m.Get("/", Index)
	m.Get("/ws", WebSocket)

	m.Run()
}

func Index(r render.Render) {
	r.HTML(200, "index", "")
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("WebSocket:", r)
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshae", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	store, _ := model.NewStore(Dial, DB, Collection)

	client := ws.RemoteAddr()
	clientConn := ClientConn{ws, client}
	addClient(clientConn)

	defer store.Close()

	// Optional. Switch the session to a monotonic behavior.
	store.Session.SetMode(mgo.Monotonic, true)

	c := store.Session.DB(DB).C(Collection)

	result := model.Image{}
	err = c.Find(bson.M{}).Sort("-_id").One(&result)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("result %v\n", result)

	iter := c.Find(bson.M{"_id": bson.M{"$gt": result.ID}}).Tail(1 * time.Second)

	for {
		var lastId bson.ObjectId
		for iter.Next(&result) {
			fmt.Println(result)
			lastId = result.ID
			broadcastMessage(result)
		}
		if iter.Err() != nil {

			fmt.Println(iter.Err())
			iter.Close()
			return
		}
		if iter.Timeout() {
			continue
		}
		query := c.Find(bson.M{"_id": bson.M{"$gt": lastId}})
		iter = query.Sort("$natural").Tail(1 * time.Second)
	}
	iter.Close()

}
