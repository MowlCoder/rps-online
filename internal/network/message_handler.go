package network

type MessageHandler func(connectedClient *ConnectedClient, message *Message)

type MessageHandlerRegister struct {
	handlersMap map[uint32]MessageHandler
}

func NewMessageHandlerRegister() *MessageHandlerRegister {
	return &MessageHandlerRegister{
		handlersMap: make(map[uint32]MessageHandler),
	}
}

func (r *MessageHandlerRegister) RegisterHandler(messageType uint32, handler MessageHandler) {
	r.handlersMap[messageType] = handler
}

func (r *MessageHandlerRegister) GetHandler(messageType uint32) (MessageHandler, bool) {
	handler, ok := r.handlersMap[messageType]

	return handler, ok
}
