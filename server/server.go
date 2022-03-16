package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Bananenpro/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type client struct {
	id                 string
	conn               *websocket.Conn
	server             *Server
	controlClientsLock sync.RWMutex
	controlClients     []*controlClient
	messages           [][]byte
}

type controlClient struct {
	client *client
	conn   *websocket.Conn
}

type Server struct {
	clients map[string]*client

	upgrader websocket.Upgrader
}

func New() *Server {
	return &Server{
		clients: make(map[string]*client),

		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) Listen(port int) {
	r := mux.NewRouter()
	r.HandleFunc("/", s.connectClient)
	r.HandleFunc("/attach/{clientId}", s.attachControlClient)
	http.Handle("/", r)

	log.Infof("Listening on port %d...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func (s *Server) connectClient(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Failed to upgrade connection with %s: %s", r.RemoteAddr, err)
		return
	}

	client := &client{
		id:             uuid.NewString(),
		conn:           conn,
		server:         s,
		controlClients: make([]*controlClient, 0),
	}

	s.clients[client.id] = client

	go client.handleConnection()

	log.Infof("Client %s connected with id %s.", client.conn.RemoteAddr(), client.id)
}

func (s *Server) attachControlClient(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Failed to upgrade connection with %s: %s", r.RemoteAddr, err)
		return
	}

	clientId := mux.Vars(r)["clientId"]
	client := s.clients[clientId]
	if client == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("invalid client id"))
		conn.Close()
		return
	}

	ctrl := &controlClient{
		client: client,
		conn:   conn,
	}

	client.controlClientsLock.Lock()
	client.controlClients = append(client.controlClients, ctrl)
	client.controlClientsLock.Unlock()

	for _, m := range client.messages {
		conn.WriteMessage(websocket.TextMessage, m)
	}

	go ctrl.handleConnection()

	log.Tracef("Control client %s connected and attached itself to client %s.", client.conn.RemoteAddr(), client.id)
}

func (c *client) handleConnection() {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				log.Errorf("Failed to receive message from client %s: %s", c.id, err)
			}
			break
		}
		if msgType != websocket.TextMessage {
			log.Warnf("Received unsupported message type from client %s. Only text messages are supported!", c.id)
		}

		c.messages = append(c.messages, msg)

		c.controlClientsLock.RLock()
		for _, ctrl := range c.controlClients {
			err = ctrl.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Errorf("Failed to send message from client %s to control client %s: %s", c.id, ctrl.conn.RemoteAddr(), err)
			}
		}
		c.controlClientsLock.RUnlock()
	}

	c.controlClientsLock.RLock()
	for _, ctrl := range c.controlClients {
		err := ctrl.conn.Close()
		if err != nil {
			log.Errorf("Failed close connection with control client %s: %s", ctrl.conn.RemoteAddr, err)
		}
	}
	c.controlClientsLock.RUnlock()

	delete(c.server.clients, c.id)
	c.conn.Close()

	log.Infof("Client %s disconnected.", c.id)
}

func (c *controlClient) handleConnection() {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				log.Errorf("Failed to receive message from control client %s: %s", c.conn.RemoteAddr(), err)
			}
			break
		}
		if msgType != websocket.TextMessage {
			log.Warnf("Received unsupported message type from control client %s. Only text messages are supported!", c.conn.RemoteAddr())
		}

		err = c.client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Errorf("Failed to send message to client %s: %s", c.client.id, err)
		}
	}

	c.client.controlClientsLock.Lock()
	index := -1
	for i, ctrl := range c.client.controlClients {
		if ctrl == c {
			index = i
		}
	}
	c.client.controlClients[index] = c.client.controlClients[len(c.client.controlClients)-1]
	c.client.controlClients = c.client.controlClients[:len(c.client.controlClients)-1]
	c.client.controlClientsLock.Unlock()

	c.conn.Close()
	log.Tracef("Control client %s disconnected.", c.conn.RemoteAddr())
}
