package gamepackets

import (
	"encoding/binary"
	"fmt"
)

type PingPacket struct {
	PingId uint8
}

func (p PingPacket) Type() PlayerPacketType {
	return Ping
}

func (p PingPacket) Marshal() ([]byte, error) {

	// Create packet: [type][ping id]
	packet := make([]byte, 2)
	packet[0] = byte(p.Type())
	packet[1] = p.PingId

	return packet, nil
}

type PongPacket struct {
	PingId uint8
}

func (p PongPacket) Type() HostPacketType {
	return Pong
}

func (p PongPacket) Marshal() ([]byte, error) {

	// Create packet: [type][ping id]
	packet := make([]byte, 2)
	packet[0] = byte(p.Type())
	packet[1] = p.PingId

	return packet, nil
}

type PingDataPacket struct {
	HostID      uint32       // ID of the host/sender
	PlayerPings []PlayerPing // List of player pings
}

type PlayerPing struct {
	PlayerID uint32
	Ping     uint16 // 0-65535, 65535 means unknown/default
}

func (p PingDataPacket) Type() PlayerPacketType {
	return PingData
}

func (p PingDataPacket) Marshal() ([]byte, error) {
	// Calculate total size: 1 (type) + 6 (host) + (6 * numPlayers)
	totalSize := 1 + 6 + (6 * len(p.PlayerPings))
	buf := make([]byte, totalSize)

	// Byte 0: Packet type
	buf[0] = byte(PingData)

	// Bytes 1-6: Host info
	// Host ID (4 bytes, little-endian)
	binary.LittleEndian.PutUint32(buf[1:5], p.HostID)
	// Padding (2 bytes zero)
	buf[5] = 0
	buf[6] = 0

	// Bytes 7+: Player ping entries
	offset := 7
	for _, player := range p.PlayerPings {
		if offset+6 > len(buf) {
			return nil, fmt.Errorf("buffer overflow")
		}

		// Player ID (4 bytes, little-endian)
		binary.LittleEndian.PutUint32(buf[offset:offset+4], player.PlayerID)

		// Ping value (2 bytes, little-endian)
		ping := player.Ping
		if ping == 0 {
			ping = 65535 // Default if not set (matches JS: e.ping ?? 65535)
		}
		binary.LittleEndian.PutUint16(buf[offset+4:offset+6], ping)

		offset += 6
	}

	return buf, nil
}
