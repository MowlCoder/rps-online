package network

import (
	"net"

	"github.com/MowlCoder/rps-online/internal/domain"
	"github.com/MowlCoder/rps-online/internal/id"
)

var userIDGenerator = id.NewGenerator()

type ConnectedClient struct {
	domain.User
	conn net.Conn
}

func NewConnectedClient(conn net.Conn) *ConnectedClient {
	user := domain.User{
		ID:       userIDGenerator.NextID(),
		Username: "",
		RoomID:   -1,
	}

	return &ConnectedClient{
		User: user,
		conn: conn,
	}
}

func (cu *ConnectedClient) SendMessage(msg *Message) {
	cu.conn.Write(msg.Encode())
}

func (cu *ConnectedClient) SendRawBytes(b []byte) {
	cu.conn.Write(b)
}

func (cu *ConnectedClient) Disconnect() {
	cu.conn.Close()
}
