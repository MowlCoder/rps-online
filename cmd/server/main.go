package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"

	"github.com/MowlCoder/rps-online/internal/handlers"
	"github.com/MowlCoder/rps-online/internal/network"
	"github.com/MowlCoder/rps-online/internal/storage"
)

var messageHandlerRegister = network.NewMessageHandlerRegister()

func main() {
	listener, err := net.Listen("tcp", ":9090")

	if err != nil {
		panic(err)
	}

	defer listener.Close()

	h := handlers.NewHandler(network.NewClientManager(), storage.NewRoomStorage())

	messageHandlerRegister.RegisterHandler(network.CONNECT_CLIENT_EVENT, h.HandleConnectClientMessage)
	messageHandlerRegister.RegisterHandler(network.CREATE_ROOM_CLIENT_EVENT, h.HandleCreateRoomMessage)
	messageHandlerRegister.RegisterHandler(network.JOIN_ROOM_CLIENT_EVENT, h.HandleJoinRoomMessage)
	messageHandlerRegister.RegisterHandler(network.DO_TURN_CLIENT_EVENT, h.HandleDoTurnMessage)

	fmt.Println("Server is listenin on port 9090")

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Failed to accept connection", err)
			return
		}

		go handleUserConnection(conn, h)
	}
}

func handleUserConnection(c net.Conn, handler *handlers.Handler) {
	connectedClient := network.NewConnectedClient(c)

	defer func() {
		handler.HandleUserDisconnect(connectedClient)
	}()

	handler.HandleUserConnect(connectedClient)

	messageSizeBuffer := make([]byte, 4)

	for {
		_, err := c.Read(messageSizeBuffer)

		if err != nil {
			return
		}

		var messageSize uint32
		binary.Read(bytes.NewBuffer(messageSizeBuffer), binary.BigEndian, &messageSize)

		messageBuffer := make([]byte, messageSize)
		_, err = c.Read(messageBuffer)

		if err != nil {
			return
		}

		message := &network.Message{}
		message.Decode(messageBuffer)

		handler, ok := messageHandlerRegister.GetHandler(message.EventType)
		if !ok {
			fmt.Printf("Handler for message with type %d hasn't registered yet\n", message.EventType)
			continue
		}

		handler(connectedClient, message)
	}
}
