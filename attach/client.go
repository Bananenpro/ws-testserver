package attach

import (
	"os"

	"github.com/Bananenpro/log"
	"github.com/Bananenpro/ws-testserver/cli"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

func Connect(url string) (*Client, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Send(message string) error {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Error("Error while sending message:", err)
		return err
	}
	log.Info("Sent:", message)
	return nil
}

func (c *Client) Listen() error {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				return err
			}
			break
		}
		if msgType != websocket.TextMessage {
			cli.PrintError("Received an invalid message type. Only text messages are supported!")
			continue
		}
		cli.PrintMessage(string(msg))
	}
	c.conn.Close()
	cli.PrintMessage("Disconnected.")
	os.Exit(0)
	return nil
}
