package network

import (
	"sync"
)

type ClientManager struct {
	mu      sync.RWMutex
	connMap map[int]*ConnectedClient
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		connMap: make(map[int]*ConnectedClient),
	}
}

func (p *ClientManager) BroadcastMessage(message *Message) {
	encodedMessage := message.Encode()

	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, user := range p.connMap {
		user.SendRawBytes(encodedMessage)
	}
}

func (p *ClientManager) GetUser(id int) (*ConnectedClient, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	user, ok := p.connMap[id]
	return user, ok
}

func (p *ClientManager) PutUser(id int, user *ConnectedClient) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.connMap[id] = user
}

func (p *ClientManager) RemoveUser(id int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.connMap, id)
}
