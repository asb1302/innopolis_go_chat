package handler

import (
	"github.com/asb1302/innopolis_go_chat/pkg/authclient"
	"github.com/asb1302/innopolis_go_chat/pkg/chatdata"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

const HeaderAuthorization = "Authorization"
const HeaderUserID = "User-ID"
const HeaderUserRole = "User-Role"

var upgrader = websocket.Upgrader{
	HandshakeTimeout: time.Minute,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  &sync.Pool{},
	CheckOrigin: func(r *http.Request) bool {
		return true // Пропускаем любой запрос
	},
}

func HandleHTTPReq(resp http.ResponseWriter, req *http.Request) {
	defer func() {
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Add("Access-Control-Allow-Methods", "GET")
		resp.Header().Add("Access-Control-Allow-Methods", "OPTIONS")
	}()

	token := req.Header.Get(HeaderAuthorization)

	if token == "" {
		resp.WriteHeader(http.StatusUnauthorized)
		log.Println("Get request", req.Method, token, "error", http.StatusUnauthorized)
		return
	}

	userID, valid := authclient.ValidateToken(token)
	if !valid {
		resp.WriteHeader(http.StatusUnauthorized)
		log.Println("Get request", req.Method, token, "error", http.StatusUnauthorized)
		return
	}

	log.Println("userID", userID)

	// обновление соединения до WebSocket
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		// Upgrade сам вставляет статус код в респонс
		log.Println(err)
		return
	}
	HandleWsConn(conn, chatdata.ID(userID))
}
