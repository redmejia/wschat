package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

var wsChan = make(chan WsPayload)

var clients = make(map[WebSocketConnection]string)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(), // only for develoment
)

// Home page
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WsJsonResp json response format
type WsJsonResp struct {
	Action      string `json:"action"`
	Message     string `json:"message"`
	MessageType string `json:"message_type"`
}

// WebSocketConnection connection
type WebSocketConnection struct {
	*websocket.Conn
}

// WsPalload  WebSocket payload
type WsPayload struct {
	Action   string              `json:"action"`
	UserName string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// WsEndPoint updgrade connection
func WsEndPoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("client connected to encpoint ")

	var response WsJsonResp
	response.Message = `<em><small>connected to serv</small></em>`

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}

	go ListenForWs(&conn)
}
func ListenForWs(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()

	var payload WsPayload

	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			//
		} else {
			payload.Conn = *conn
			wsChan <- payload
		}
	}
}
func ListenToWsChannel() {
	var response WsJsonResp
	for {
		e := <-wsChan

		response.Action = "Got Here"
		response.Message = fmt.Sprintf("Some message %s", e.Action)
		bradcastToAll(response)
	}
}
func bradcastToAll(response WsJsonResp) {
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			log.Println("Websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Panicln(err)
		return err
	}

	err = view.Execute(w, data, nil)

	if err != nil {
		log.Panicln(err)
		return err
	}

	return nil
}
