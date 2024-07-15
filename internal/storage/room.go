package storage

import (
	"sync"

	"github.com/MowlCoder/rps-online/internal/domain"
	"github.com/MowlCoder/rps-online/internal/id"
)

type IDGenerator interface {
	NextID() int
}

type RoomStorage struct {
	mu          sync.RWMutex
	idGenerator IDGenerator
	storage     map[int]domain.Room
}

func NewRoomStorage() *RoomStorage {
	return &RoomStorage{
		idGenerator: id.NewGenerator(),
		storage:     make(map[int]domain.Room),
	}
}

func (s *RoomStorage) GetByID(id int) (*domain.Room, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	room, ok := s.storage[id]
	return &room, ok
}

func (s *RoomStorage) GetAll() []domain.Room {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rooms := make([]domain.Room, 0)

	for _, room := range s.storage {
		rooms = append(rooms, room)
	}

	return rooms
}

func (s *RoomStorage) Put(id int, room *domain.Room) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[id] = *room
}

func (s *RoomStorage) AddNewRoom(name string, creator domain.User) *domain.Room {
	room := domain.Room{
		ID:       s.idGenerator.NextID(),
		Name:     name,
		Creator:  creator,
		Opponent: nil,
		Status:   domain.ROOM_WAITING_PLAYERS,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[room.ID] = room

	return &room
}

func (s *RoomStorage) DeleteByID(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.storage, id)
}
