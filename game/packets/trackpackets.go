package gamepackets

import (
	"encoding/hex"
	"fmt"
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
