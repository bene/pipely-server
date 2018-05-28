package server

type Event struct {
	ChannelId string      `json:"channel_id"`
	Type      string      `json:"type"`
	OriginId  string      `json:"origin_id"`
	Payload   interface{} `json:"payload"`
}

type Channel struct {
	ChannelId string
	Password  string
	Clients   []Client
}

type Client struct {
	ClientId  string
	ChannelId string
	Channel   chan Event
}
