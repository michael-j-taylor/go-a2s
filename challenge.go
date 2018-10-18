package a2s

import (
	"errors"
)

const (
	A2S_PLAYER_CHALLENGE_REPLY_HEADER = 0x41
)

var (
	ErrBadChallengeResponse = errors.New("Bad challenge response")
)

func (c *Client) GetChallenge(header byte, fullResult byte) ([]byte, bool, error) {
	if err := c.Send([]byte{
		0xff, 0xff, 0xff, 0xff,
		header,
		0xff, 0xff, 0xff, 0xff,
	}); err != nil {
		return nil, false, err
	}

	data, err := c.Receive()

	if err != nil {
		return nil, false, err
	}

	reader := NewPacketReader(data)

	switch int32(reader.ReadUint32()) {
	case -2:
		// We received an unexpected full reply

		return data, true, nil
	case -1:
		// Continue
	default:
		panic(ErrBadPacketHeader)
	}

	switch reader.ReadUint8() {
	case A2S_PLAYER_CHALLENGE_REPLY_HEADER:
		// Received a challenge number

		return data[reader.Pos() : reader.Pos()+4], false, nil
	case fullResult:
		// Received full result

		return data, true, nil
	}

	return nil, false, ErrBadChallengeResponse
}
