package gamepackets

import (
	"fmt"
)

// PacketFactory helps create packets from raw data (for receiving)
type PacketFactory struct{}

func (f *PacketFactory) FromBytes(data []byte) (HostPacket, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty packet")
	}

	packetType := HostPacketType(data[0])

	switch packetType {

	default:
		return nil, fmt.Errorf("unknown packet type: %d", packetType)
	}
}
