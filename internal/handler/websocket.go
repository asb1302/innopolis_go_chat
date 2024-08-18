package handler

import (
	"encoding/json"
	"github.com/asb1302/innopolis_go_chat/internal/service"
	"github.com/asb1302/innopolis_go_chat/internal/service/pools"
	"github.com/asb1302/innopolis_go_chat/pkg/chatdata"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait  = 2 * time.Second
	pongWait   = 20 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

func HandleWsConn(conn *websocket.Conn, UID chatdata.ID) {

	defer func() {
		// closing the user channel and ending write goroutine
		if pools.Users.Delete(UID) {
			conn.Close()
		}
	}()

	ch := pools.Users.New(UID)

	// write to conn from channel
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			if pools.Users.Delete(UID) {
				conn.Close()
			}
		}()
		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					log.Println("channel for " + UID + " closed")
					return
				}
				log.Println("SEND to ", UID, msg)

				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteJSON(msg); err != nil {
					handleWsError(err, UID)
					return
				}
			case <-ticker.C:
				//log.Println("SEND Ping to", UID)
				conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					handleWsError(err, UID)
					return
				}
			}
		}
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		//log.Println("GOT Pong from", UID)
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	// read from conn
	for {
		typ, message, err := conn.ReadMessage()
		if err != nil {
			handleWsError(err, UID)
			return
		}
		log.Println("GOT from", UID, typ, string(message))
		switch typ {
		case websocket.TextMessage, websocket.BinaryMessage:
			var req chatdata.Request
			if err = json.Unmarshal(message, &req); err != nil {
				sendErrorResp(UID, err)
				continue
			}
			switch req.Type {
			case chatdata.ReqTypeNewChat:
				var newChatReq chatdata.NewChatRequest
				if err = json.Unmarshal(req.Data, &newChatReq); err != nil {
					sendErrorResp(UID, err)
					continue
				}
				chatid := service.NewChat(append(newChatReq.UserIDs, UID))
				sendResp(UID, chatdata.DeliveryTypeNewChat, chatid)

			case chatdata.ReqTypeNewMsg:
				var msg chatdata.MessageChatRequest
				if err = json.Unmarshal(req.Data, &msg); err != nil {
					sendErrorResp(UID, err)
					continue
				}

				switch msg.Type {
				case chatdata.MsgTypeAdd:
					if err := service.NewMessage(msg, UID); err != nil {
						sendErrorResp(UID, err)
						continue
					}
				}
			}
		case websocket.CloseMessage:
			return
		}
	}
}

func handleWsError(err error, uid chatdata.ID) {
	switch {
	case websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway):
		log.Println("websocket session closed by client", uid)
	default:
		log.Println("error websocket message", err.Error(), "for", uid)
	}
}

func sendErrorResp(UID chatdata.ID, err error) {
	sendResp(UID, chatdata.DeliveryTypeError, chatdata.ErrorResponse{Error: err.Error()})
}

func sendResp(UID chatdata.ID, typ chatdata.DeliveryType, data interface{}) {

	var resp chatdata.Delivery
	resp.Type = typ
	resp.Data = data
	//log.Println("SEND to channel", UID, resp)
	pools.Users.Send(UID, resp)
}
