package gamepackets

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
)

// TrackIDPacket represents the initial track ID packet
type TrackIDPacket struct {
	TrackID string // 64-char hex string
}

func (p TrackIDPacket) Type() PlayerPacketType {
	return TrackID
}

func (p TrackIDPacket) Marshal() ([]byte, error) {
	if len(p.TrackID) != 64 {
		return nil, fmt.Errorf("invalid track ID length: expected 64, got %d", len(p.TrackID))
	}

	// Decode hex to bytes
	trackIDBytes, err := hex.DecodeString(p.TrackID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode track ID: %w", err)
	}

	// Create packet: [type][32 bytes track ID]
	packet := make([]byte, 33)
	packet[0] = byte(p.Type())
	copy(packet[1:], trackIDBytes)

	return packet, nil
}

// TrackChunkPacket represents a chunk of track data
type TrackChunkPacket struct {
	Data []byte // Raw ASCII data chunk
}

func (p TrackChunkPacket) Type() PlayerPacketType {
	return TrackChunk
}

func (p TrackChunkPacket) Marshal() ([]byte, error) {
	// Create packet: [type][data...]
	packet := make([]byte, 1+len(p.Data))
	packet[0] = byte(p.Type())
	copy(packet[1:], p.Data)

	return packet, nil
}

// PlayerUpdatePacket example
type PlayerUpdatePacket struct {
	PlayerID  uint32
	PositionX float32
	PositionY float32
	PositionZ float32
	Rotation  float32
	// ... other fields
}

func (p PlayerUpdatePacket) Type() PlayerPacketType {
	return PlayerUpdate
}

func (p PlayerUpdatePacket) Marshal() ([]byte, error) {
	// Example binary marshaling - adjust based on your protocol
	buf := make([]byte, 1+4+4+4+4+4) // type + fields
	buf[0] = byte(p.Type())

	binary.LittleEndian.PutUint32(buf[1:5], p.PlayerID)
	binary.LittleEndian.PutUint32(buf[5:9], math.Float32bits(p.PositionX))
	binary.LittleEndian.PutUint32(buf[9:13], math.Float32bits(p.PositionY))
	binary.LittleEndian.PutUint32(buf[13:17], math.Float32bits(p.PositionZ))
	binary.LittleEndian.PutUint32(buf[17:21], math.Float32bits(p.Rotation))

	return buf, nil
}

// PingPacket (simple packet with no data)
type PingPacket struct{}

func (p PingPacket) Type() PlayerPacketType {
	return Ping
}

func (p PingPacket) Marshal() ([]byte, error) {
	return []byte{byte(p.Type())}, nil
}

// KickPacket example
type KickPacket struct {
	Reason string
}

func (p KickPacket) Type() PlayerPacketType {
	return Kick
}

func (p KickPacket) Marshal() ([]byte, error) {
	reasonBytes := []byte(p.Reason)
	if len(reasonBytes) > 255 {
		return nil, fmt.Errorf("kick reason too long: %d bytes max 255", len(reasonBytes))
	}

	// Format: [type][reason length][reason...]
	packet := make([]byte, 1+1+len(reasonBytes))
	packet[0] = byte(p.Type())
	packet[1] = byte(len(reasonBytes))
	copy(packet[2:], reasonBytes)

	return packet, nil
}
