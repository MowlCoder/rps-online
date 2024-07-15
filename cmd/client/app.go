package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"

	"github.com/MowlCoder/rps-online/internal/network"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	connected  bool
	socketConn net.Conn
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		return
	}

	a.socketConn = conn
	a.connected = true

	go a.socketListener()
}

func (a *App) IsConnectedToServer() bool {
	return a.connected
}

func (a *App) Login(name string) {
	msg := network.NewMessage(network.CONNECT_CLIENT_EVENT, network.ConnectClientPayload{
		Username: name,
	})

	a.socketConn.Write(msg.Encode())
}

func (a *App) CreateRoom(name string) {
	msg := network.NewMessage(network.CREATE_ROOM_CLIENT_EVENT, network.CreateRoomClientPayload{
		RoomName: name,
	})

	a.socketConn.Write(msg.Encode())
}

func (a *App) JoinRoom(roomID int) {
	msg := network.NewMessage(network.JOIN_ROOM_CLIENT_EVENT, network.JoinRoomClientPayload{
		RoomID: roomID,
	})

	a.socketConn.Write(msg.Encode())
}

func (a *App) MakeChoice(choice uint8) {
	msg := network.NewMessage(network.DO_TURN_CLIENT_EVENT, network.DoTurnClientPayload{
		Choice: choice,
	})

	a.socketConn.Write(msg.Encode())
}

func (a *App) socketListener() {
	messageSizeBuffer := make([]byte, 4)

	for {
		_, err := a.socketConn.Read(messageSizeBuffer)

		if err != nil {
			if err == io.EOF {
				runtime.EventsEmit(a.ctx, "server:no_connection")
			}

			return
		}

		var messageSize uint32
		binary.Read(bytes.NewBuffer(messageSizeBuffer), binary.BigEndian, &messageSize)

		messageBuffer := make([]byte, messageSize)
		_, err = a.socketConn.Read(messageBuffer)

		if err != nil {
			if err == io.EOF {
				runtime.EventsEmit(a.ctx, "server:no_connection")
			}

			return
		}

		msg := &network.Message{}
		msg.Decode(messageBuffer)

		switch msg.EventType {
		case network.CONNECT_SERVER_EVENT:
			var payload network.ConnectServerPayload
			_ = json.Unmarshal(msg.Payload, &payload)

			runtime.EventsEmit(a.ctx, "server:success_login", payload)
		case network.ROOM_CREATED_SERVER_EVENT:
			var payload network.RoomCreatedServerPayload
			_ = json.Unmarshal(msg.Payload, &payload)

			runtime.EventsEmit(a.ctx, "server:room_created", payload)
		case network.JOIN_ROOM_SERVER_EVENT:
			var payload network.JoinRoomServerPayload
			_ = json.Unmarshal(msg.Payload, &payload)

			runtime.EventsEmit(a.ctx, "server:room_joined", payload)
		case network.JOINED_ROOM_SUCCESS_SERVER_EVENT:
			var payload network.JoinedRoomSuccessServerPayload
			_ = json.Unmarshal(msg.Payload, &payload)

			runtime.EventsEmit(a.ctx, "server:joined_room_success", payload)
		case network.MATCH_END_SERVER_EVENT:
			var payload network.MatchEndServerPayload
			_ = json.Unmarshal(msg.Payload, &payload)

			runtime.EventsEmit(a.ctx, "server:match_end", payload)
		case network.ROOM_DELETED_SERVER_EVENT:
			var payload network.RoomDeletedServerPayload
			_ = json.Unmarshal(msg.Payload, &payload)

			runtime.EventsEmit(a.ctx, "server:room_deleted", payload)
		}
	}
}
