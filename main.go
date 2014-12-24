package main

import (
	"github.com/go-martini/martini"
	"github.com/gorilla/websocket"
	"github.com/martini-contrib/render"
	"github.com/masayukioguni/go-lgtm-front/config"
	"github.com/masayukioguni/go-lgtm-model"
	"path"

	"log"
	"net"
	"net/http"
	"os"
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
	config       *config.Config
}

func main() {
	c, err := config.NewConfig(".env")
	if err != nil {
		log.Panic(err)
	}

	f := &Front{
		config: c,
	}

	LogPath := os.Getenv("LOG_PATH")

	file, err := os.OpenFile(LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to os.OpenFile()")
	}
	defer file.Close()

	log.SetOutput(file)

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
	s, _ := model.NewStore(f.config.MongoHost, f.config.MongoDatabase, f.config.MongoCollectionName)

	items, _ := s.All()
	names := []string{}

	for _, item := range items {
		names = append(names, path.Join(f.config.S3Url, item.Name))
	}

	r.HTML(200, "index", names)
}

func (f *Front) CommandImage(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("Name")
	c := &ImageChannel{
		Name: name,
	}
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
			Name: path.Join(f.config.S3Url, name),
		}

		time.Sleep(1000 * time.Millisecond)
		broadcastMessage(result)
	}

}
