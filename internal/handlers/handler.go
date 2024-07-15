package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/MowlCoder/rps-online/internal/domain"
	"github.com/MowlCoder/rps-online/internal/network"
	"github.com/MowlCoder/rps-online/internal/storage"
)

type Handler struct {
	networkClientManager *network.ClientManager
	roomStorage          *storage.RoomStorage
}

func NewHandler(networkClientManager *network.ClientManager, roomStorage *storage.RoomStorage) *Handler {
	return &Handler{
		networkClientManager: networkClientManager,
		roomStorage:          roomStorage,
	}
}

func (h *Handler) HandleUserConnect(connectedClient *network.ConnectedClient) {
	fmt.Printf("Client with id %d connected\n", connectedClient.ID)
	h.networkClientManager.PutUser(connectedClient.ID, connectedClient)
}

func (h *Handler) HandleUserDisconnect(connectedClient *network.ConnectedClient) {
	fmt.Printf("Client with id %d disconnected\n", connectedClient.ID)
	defer connectedClient.Disconnect()
	h.networkClientManager.RemoveUser(connectedClient.ID)

	if connectedClient.RoomID > 0 {
		room, exist := h.roomStorage.GetByID(connectedClient.RoomID)

		if !exist {
			return
		}

		if room.Opponent != nil && room.Opponent.ID == connectedClient.ID {
			h.endRoomMatch(*room, domain.MATCH_WON, domain.MATCH_LOST)
			return
		}

		if room.Opponent != nil && room.Creator.ID == connectedClient.ID {
			h.endRoomMatch(*room, domain.MATCH_LOST, domain.MATCH_WON)
			return
		}

		h.roomStorage.DeleteByID(connectedClient.RoomID)
		roomDeletedMessage := network.NewMessage(network.ROOM_DELETED_SERVER_EVENT, network.RoomDeletedServerPayload{
			RoomID: connectedClient.RoomID,
		})
		h.networkClientManager.BroadcastMessage(roomDeletedMessage)
	}
}

func (h *Handler) HandleConnectClientMessage(connectedClient *network.ConnectedClient, message *network.Message) {
	var payload network.ConnectClientPayload
	json.Unmarshal(message.Payload, &payload)

	connectedClient.Username = payload.Username

	rooms := h.roomStorage.GetAll()

	msg := network.NewMessage(network.CONNECT_SERVER_EVENT, network.ConnectServerPayload{
		Rooms: rooms,
	})

	connectedClient.SendMessage(msg)
}

func (h *Handler) HandleCreateRoomMessage(connectedClient *network.ConnectedClient, message *network.Message) {
	var payload network.CreateRoomClientPayload
	json.Unmarshal(message.Payload, &payload)

	room := h.roomStorage.AddNewRoom(payload.RoomName, connectedClient.User)
	connectedClient.RoomID = room.ID

	msg := network.NewMessage(network.ROOM_CREATED_SERVER_EVENT, network.RoomCreatedServerPayload{
		Room: *room,
	})

	h.networkClientManager.BroadcastMessage(msg)
}

func (h *Handler) HandleJoinRoomMessage(connectedClient *network.ConnectedClient, message *network.Message) {
	var payload network.JoinRoomClientPayload
	json.Unmarshal(message.Payload, &payload)

	room, ok := h.roomStorage.GetByID(payload.RoomID)
	if !ok {
		return
	}

	connectedClient.RoomID = room.ID

	room.Opponent = &connectedClient.User
	room.Status = domain.ROOM_PLAYING
	room.Match = &domain.Match{
		Type:                  domain.MATCH_BO1,
		CreatorCount:          0,
		OpponentCount:         0,
		CurrentCreatorChoice:  domain.CHOICE_NONE,
		CurrentOpponentChoice: domain.CHOICE_NONE,
	}

	h.roomStorage.Put(payload.RoomID, room)

	msg := network.NewMessage(network.JOINED_ROOM_SUCCESS_SERVER_EVENT, network.JoinedRoomSuccessServerPayload{
		Room: *room,
	})
	connectedClient.SendMessage(msg)

	msg = network.NewMessage(network.JOIN_ROOM_SERVER_EVENT, network.JoinRoomServerPayload{
		JoinedUser: connectedClient.User,
	})

	creatorClient, exist := h.networkClientManager.GetUser(room.Creator.ID)
	if !exist {
		return
	}

	creatorClient.SendMessage(msg)
}

func (h *Handler) HandleDoTurnMessage(connectedClient *network.ConnectedClient, message *network.Message) {
	var payload network.DoTurnClientPayload
	json.Unmarshal(message.Payload, &payload)

	room, ok := h.roomStorage.GetByID(connectedClient.RoomID)
	if !ok {
		return
	}

	if connectedClient.ID == room.Creator.ID {
		if room.Match.CurrentCreatorChoice != domain.CHOICE_NONE {
			return
		}

		room.Match.CurrentCreatorChoice = payload.Choice
	} else {
		if room.Match.CurrentOpponentChoice != domain.CHOICE_NONE {
			return
		}

		room.Match.CurrentOpponentChoice = payload.Choice
	}

	if room.Match.BothPlayerHaveChosen() {
		choiceResult := domain.ComputeChoiceResult(room.Match.CurrentCreatorChoice, room.Match.CurrentOpponentChoice)

		switch choiceResult {
		case domain.CHOICE_RESULT_DRAW:
			room.Match.CreatorCount += 1
			room.Match.OpponentCount += 1
		case domain.CHOICE_RESULT_LEFT:
			room.Match.CreatorCount += 1
		case domain.CHOICE_RESULT_RIGHT:
			room.Match.OpponentCount += 1
		}

		room.Match.ResetChoices()
	}

	if room.Match.RoundsIsOver() {
		room.Status = domain.ROOM_MATCH_END
		var creatorMatchResult uint8
		var opponentMatchResult uint8

		if room.Match.CreatorCount == room.Match.OpponentCount {
			creatorMatchResult = domain.MATCH_DRAW
			opponentMatchResult = domain.MATCH_DRAW
		} else if room.Match.CreatorCount > room.Match.OpponentCount {
			creatorMatchResult = domain.MATCH_WON
			opponentMatchResult = domain.MATCH_LOST
		} else {
			creatorMatchResult = domain.MATCH_LOST
			opponentMatchResult = domain.MATCH_WON
		}

		h.endRoomMatch(*room, creatorMatchResult, opponentMatchResult)
		return
	}

	h.roomStorage.Put(room.ID, room)
}

func (h *Handler) endRoomMatch(room domain.Room, creatorMatchResult, opponentMatchResult uint8) {
	h.sendMatchResult(room.Creator.ID, creatorMatchResult, *room.Match)
	h.sendMatchResult(room.Opponent.ID, opponentMatchResult, *room.Match)

	h.roomStorage.DeleteByID(room.ID)
	roomDeletedMessage := network.NewMessage(network.ROOM_DELETED_SERVER_EVENT, network.RoomDeletedServerPayload{
		RoomID: room.ID,
	})

	h.networkClientManager.BroadcastMessage(roomDeletedMessage)
}

func (h *Handler) sendMatchResult(userID int, matchResult uint8, match domain.Match) {
	msgForCreator := network.NewMessage(network.MATCH_END_SERVER_EVENT, network.MatchEndServerPayload{
		MatchResult:    matchResult,
		CreatorChoice:  match.LastCreatorChoice,
		OpponentChoice: match.LastOpponentChoice,
	})

	client, exist := h.networkClientManager.GetUser(userID)
	if !exist {
		return
	}

	client.RoomID = 0

	client.SendMessage(msgForCreator)
	h.networkClientManager.PutUser(client.ID, client)
}
