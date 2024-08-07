package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"play/websocket/cache"
	"play/websocket/domain"
)

func HandleWsConn(conn *websocket.Conn, UID string) {
	defer func() {
		cache.Conns.DeleteConn(UID)
		conn.Close()
	}()

	cache.Conns.SetConn(UID, conn)

	for {
		typ, message, err := conn.ReadMessage()
		if err != nil {
			handleWsError(err)
			return
		}

		switch typ {
		case websocket.BinaryMessage:
			var msg domain.MessageRequest
			if err = json.Unmarshal(message, &msg); err != nil {
				return
			}
			d := cache.Dials.Get(msg.ChID)
			// лучше тут это не делать, тк существует возможность потерять сообщения
			// на подумать: почему?
			d.Messages = append(d.Messages, msg.Msg)
			cache.Dials.Set(msg.ChID, d)

			for _, uiD := range d.UIDs {
				if uiD != UID {

					othconn := cache.Conns.GetConn(uiD)
					if othconn != nil {
						if err = othconn.WriteJSON(msg); err != nil {
							handleWsError(err)
							return
						}
					}
				}
			}

		case websocket.CloseMessage:
			return
		case websocket.PingMessage:
			if err = conn.WriteMessage(websocket.PongMessage, []byte{}); err != nil {
				handleWsError(err)
				return
			}
		}

	}
}

func messageHandler(message []byte) {
	fmt.Println(string(message))
}

func handleWsError(err error) {
	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
		log.Println("websocket session closed by client")
	default:
		log.Println("error read websocket message", err.Error())
	}
}

func write(conn *websocket.Conn, ch chan domain.MessageRequest) {
	for {
		msg := <-ch
		if err := conn.WriteJSON(msg); err != nil {
			handleWsError(err)
			return
		}
	}
}
