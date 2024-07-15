package network

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"unsafe"

	"github.com/MowlCoder/rps-online/internal/domain"
	"github.com/MowlCoder/rps-online/internal/id"
)

const (
	CONNECT_SERVER_EVENT = iota
	CONNECT_CLIENT_EVENT
	CREATE_ROOM_CLIENT_EVENT
	ROOM_CREATED_SERVER_EVENT
	ROOM_DELETED_SERVER_EVENT
	JOIN_ROOM_CLIENT_EVENT
	JOIN_ROOM_SERVER_EVENT
	JOINED_ROOM_SUCCESS_SERVER_EVENT
	DO_TURN_CLIENT_EVENT
	MATCH_END_SERVER_EVENT
)

type ConnectClientPayload struct {
	Username string
}

type ConnectServerPayload struct {
	Rooms []domain.Room
}

type CreateRoomClientPayload struct {
	RoomName string
}

type RoomCreatedServerPayload struct {
	Room domain.Room
}

type RoomDeletedServerPayload struct {
	RoomID int
}

type JoinRoomClientPayload struct {
	RoomID int
}

type JoinRoomServerPayload struct {
	JoinedUser domain.User
}

type JoinedRoomSuccessServerPayload struct {
	Room domain.Room
}

type DoTurnClientPayload struct {
	Choice uint8
}

type MatchEndServerPayload struct {
	MatchResult    uint8
	CreatorChoice  uint8
	OpponentChoice uint8
}

var messageIDGenerator = id.NewGenerator()

type Message struct {
	ID        int32
	EventType uint32
	Payload   []byte
}

func NewMessage(eventType uint32, payload any) *Message {
	payloadBytes, _ := json.Marshal(payload)

	return &Message{
		ID:        int32(messageIDGenerator.NextID()),
		EventType: eventType,
		Payload:   payloadBytes,
	}
}

func (m *Message) Encode() []byte {
	bin_buf := new(bytes.Buffer)

	binary.Write(bin_buf, binary.BigEndian, m.ID)
	binary.Write(bin_buf, binary.BigEndian, m.EventType)
	binary.Write(bin_buf, binary.BigEndian, m.Payload)

	return bin_buf.Bytes()
}

func (m *Message) Decode(b []byte) {
	buffer := bytes.NewBuffer(b)

	binary.Read(buffer, binary.BigEndian, &m.ID)
	binary.Read(buffer, binary.BigEndian, &m.EventType)
	m.Payload = b[unsafe.Sizeof(m.ID)+unsafe.Sizeof(m.EventType):]
}
