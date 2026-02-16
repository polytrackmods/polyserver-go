package gamepackets

import "fmt"

type HostPacketType uint8

const (
	HostCarReset HostPacketType = iota
	HostCarUpdate
	HostRecord
	Pong
	HostModCustomMessage
)

func (pt HostPacketType) String() string {
	switch pt {
	case HostCarReset:
		return "HostCarReset"
	case HostCarUpdate:
		return "HostCarUpdate"
	case HostRecord:
		return "HostRecord"
	case Pong:
		return "Pong"
	case HostModCustomMessage:
		return "HostModCustomMessage"
	default:
		return fmt.Sprintf("Unknown(%d)", pt)
	}
}

type PlayerPacketType uint8

const (
	PlayerUpdate PlayerPacketType = iota
	RemovePlayer
	PlayerCarReset
	PlayerCarUpdate
	Kick
	TrackID
	TrackChunk
	EndSession
	NewSession
	Ping
	PingData
	PlayerModCustomMessage
)

func (pt PlayerPacketType) String() string {
	switch pt {
	case PlayerUpdate:
		return "PlayerUpdate"
	case RemovePlayer:
		return "RemovePlayer"
	case PlayerCarReset:
		return "PlayerCarReset"
	case PlayerCarUpdate:
		return "PlayerCarUpdate"
	case Kick:
		return "Kick"
	case TrackID:
		return "TrackID"
	case TrackChunk:
		return "TrackChunk"
	case EndSession:
		return "EndSession"
	case NewSession:
		return "NewSession"
	case Ping:
		return "Ping"
	case PingData:
		return "PingData"
	case PlayerModCustomMessage:
		return "PlayerModCustomMessage"
	default:
		return fmt.Sprintf("Unknown(%d)", pt)
	}
}

// Packet is the interface that all packets must implement
type PlayerPacket interface {
	Type() PlayerPacketType
	Marshal() ([]byte, error)
}

type HostPacket interface {
	Type() PlayerPacketType
	Marshal() ([]byte, error)
}

// Sender handles sending packets through WebRTC
type Sender interface {
	Send(packet PlayerPacket) error
}
