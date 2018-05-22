package server

type Event struct {
	ChannelId string      `json:"channel_id"`
	Type      string      `json:"type"`
	Origin    string      `json:"origin"`
	Payload   interface{} `json:"payload"`
}

type Client struct {
	ClientId  string
	ChannelId string
	Channel   chan Event
}
