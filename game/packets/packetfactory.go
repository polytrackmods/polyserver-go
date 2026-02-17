package gamepackets

import (
	"encoding/binary"
	"fmt"
	"log"
)

// PacketFactory helps create packets from raw data (for receiving)
type PacketFactory struct{}

func (f *PacketFactory) FromBytes(data []byte) (HostPacket, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty packet")
	}

	packetType := HostPacketType(data[0])

	switch packetType {
	case Pong:
		return PongPacket{
			PingId: data[1],
		}, nil
	case HostCarUpdate:
		return nil, fmt.Errorf("unknown packet type: %d", packetType) // TODO: Proper car update deserializing

		carState, _, err := DecodeCarState(data[12:])
		if err != nil {
			log.Println("Error decoding car state: " + err.Error())
		}
		return HostCarUpdatePacket{
			SessionID:    binary.LittleEndian.Uint32(data[2:7]),
			ResetCounter: binary.LittleEndian.Uint32(data[7:12]),
			CarState:     carState,
		}, nil
	default:
		return nil, fmt.Errorf("unknown packet type: %d", packetType)
	}
}
