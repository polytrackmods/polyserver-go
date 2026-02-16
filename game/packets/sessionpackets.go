package gamepackets

import (
	"encoding/binary"
)

// TrackIDPacket represents the initial track ID packet
type EndSessionPacket struct{}

func (p EndSessionPacket) Type() PlayerPacketType {
	return EndSession
}

func (p EndSessionPacket) Marshal() ([]byte, error) {

	// Create packet: [type]
	packet := make([]byte, 1)
	packet[0] = byte(p.Type())

	return packet, nil
}

type NewSessionPacket struct {
	SessionID  uint32 // The session identifier (un in JS)
	GameMode   uint8  // The track identifier (dn in JS)
	MaxPlayers uint8  // Maximum players allowed (fn in JS)
}

func (p NewSessionPacket) Type() PlayerPacketType {
	return NewSession
}

func (p NewSessionPacket) Marshal() ([]byte, error) {
	// Create 7-byte packet
	// Format: [type][sessionID (4 bytes little-endian)][trackID][maxPlayers]
	packet := make([]byte, 7)
	packet[0] = byte(p.Type())

	// Write session ID as little-endian uint32 (bytes 1-4)
	binary.LittleEndian.PutUint32(packet[1:5], p.SessionID)

	// Write track ID (byte 5)
	packet[5] = p.GameMode

	// Write max players (byte 6)
	packet[6] = p.MaxPlayers

	return packet, nil
}
