package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/martini-contrib/render"
	"github.com/masayukioguni/go-lgtm-front/config"
	"github.com/masayukioguni/go-lgtm-model"
	"path"

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

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteJSON(result); err != nil {
			return
		}
	}
}

type ImageChannel struct {
	Name string
}

type Front struct {
	ImageChannel chan *ImageChannel
	m            *martini.ClassicMartini
	S3Url        string
}

func main() {
	c, err := config.NewConfig(".env")
	if err != nil {
		log.Panic(err)
	}

	f := &Front{
		S3Url: c.S3Url,
	}

	f.ImageChannel = make(chan *ImageChannel)

	f.m = martini.Classic()
	f.m.Use(render.Renderer())
	f.m.Use(martini.Static("assets"))

	f.m.Get("/", f.Index)
	f.m.Post("/command/image", f.CommandImage)
	f.m.Get("/ws", f.WebSocket)

	f.m.Run()
}

func (f *Front) Index(r render.Render) {
	r.HTML(200, "index", "")
}

func (f *Front) CommandImage(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("Name")
	c := &ImageChannel{
		Name: name,
	}
	fmt.Printf("f.ImageChannel <- c\n")
	f.ImageChannel <- c
}

func (f *Front) WebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshae", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	client := ws.RemoteAddr()
	clientConn := ClientConn{ws, client}
	addClient(clientConn)

	for {
		c := <-f.ImageChannel
		name := c.Name
		result := model.Image{
			Name: path.Join(f.S3Url, name),
		}

		time.Sleep(1000 * time.Millisecond)
		broadcastMessage(result)
	}

}
