package main

import (
	"github.com/fukata/golang-stats-api-handler"
	"github.com/gorilla/websocket"
	"github.com/masayukioguni/go-lgtm-front/config"
	"github.com/masayukioguni/go-lgtm-front/templates"
	"github.com/masayukioguni/go-lgtm-model"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web/middleware"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"sync"
	"text/template"
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

	goji.Use(middleware.Recoverer)
	goji.Use(middleware.NoCache)

	goji.Get("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	goji.Get("/", f.Index)
	goji.Get("/stats", stats_api.Handler)

	goji.Post("/command/image", f.CommandImage)
	goji.Get("/ws", f.WebSocket)

	goji.Serve()
}

type IndexView struct {
	WebSocketUrl string
	Names        []string
}

func (f *Front) Index(w http.ResponseWriter, r *http.Request) {
	s, _ := model.NewStore(f.config.MongoHost, f.config.MongoDatabase, f.config.MongoCollectionName)

	items, _ := s.All()
	names := []string{}

	for _, item := range items {
		names = append(names, path.Join(f.config.S3Url, item.Name))
	}

	indexView := &IndexView{
		WebSocketUrl: f.config.WebSocketUrl,
		Names:        names,
	}

	t := template.New("index")
	tmpl, _ := templates.Asset("assets/templates/index.tmpl")
	t, _ = t.Parse(string(tmpl))

	//t := template.Must(template.ParseFile("templates/index.tmpl"))
	err := t.Execute(w, indexView)
	if err != nil {
		log.Printf("template execution: %s\n", err)
	}

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
