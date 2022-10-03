package ws

type Hub struct {
	ChannelSenders   map[string]*Channel
	ChannelReceivers map[string]*Channel
}

func NewHub() *Hub {
	return &Hub{
		ChannelSenders:   make(map[string]*Channel),
		ChannelReceivers: make(map[string]*Channel),
	}
}
