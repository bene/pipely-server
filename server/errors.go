package server

import "errors"

var (
	ErrorInvalidChannelId       = errors.New("channel id must be twelve characters long")
	ErrorChannelDoesNotExist    = errors.New("channel with given id does not exist")
	ErrorInvalidChannelPassword = errors.New("invalid channel password")
	ErrorInvalidClientId        = errors.New("client id must not be shorter than three characters")
	ErrorClientIdAlreadyUsed    = errors.New("client id already used in channel")
	ErrorInvalidEventType       = errors.New("invalid event type, must not be shorter than one character")
	ErrorInvalidOriginId        = errors.New("origin id must not be shorter than three characters")
	ErrorStreamingUnsupported   = errors.New("client does not support streaming")
)
