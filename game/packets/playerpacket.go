package gamepackets

import (
	"encoding/binary"
	"fmt"
	"log"
)

type PlayerUpdatePacket struct {
	ID          uint32
	Nickname    string
	CountryCode *string
	CarStyle    *CarStyle
	NumFrames   *uint32
}

func (p PlayerUpdatePacket) Type() PlayerPacketType {
	return PlayerUpdate
}

func (p PlayerUpdatePacket) Marshal() ([]byte, error) {
	buf := make([]byte, 0)

	// Packet type
	buf = append(buf, byte(PlayerUpdate))

	// Player ID (32-bit little-endian)
	buf = binary.LittleEndian.AppendUint32(buf, p.ID)

	// Nickname (length-prefixed UTF-8)
	nicknameBytes := []byte(p.Nickname)
	if len(nicknameBytes) > 255 {
		return nil, fmt.Errorf("nickname too long: %d bytes", len(nicknameBytes))
	}
	buf = append(buf, byte(len(nicknameBytes)))
	buf = append(buf, nicknameBytes...)

	// Country code (null-terminated)
	if p.CountryCode == nil {
		buf = append(buf, 0)
	} else {
		for _, c := range *p.CountryCode {
			if c > 255 {
				return nil, fmt.Errorf("country code contains non-ASCII character")
			}
			buf = append(buf, byte(c))
		}
		buf = append(buf, 0)
	}

	// Car style - encode to binary (16 bytes) for network transmission
	if p.CarStyle != nil {
		carStyleBytes := p.CarStyle.EncodeCarStyle()
		buf = append(buf, carStyleBytes...)
	} else {
		// Send default car style if nil
		defaultStyle := DefaultCarStyle()
		buf = append(buf, defaultStyle.EncodeCarStyle()...)
	}

	// Frame count (optional)
	if p.NumFrames == nil {
		buf = append(buf, 0)
	} else {
		buf = append(buf, 1)
		log.Printf("Sending frames: %v\n", *p.NumFrames)
		tempBuf := make([]byte, 0)
		tempBuf = binary.LittleEndian.AppendUint32(tempBuf, *p.NumFrames)
		buf = append(buf, tempBuf[0])
		buf = append(buf, tempBuf[1])
		buf = append(buf, tempBuf[2])
	}

	return buf, nil
}

type RemovePlayerPacket struct {
	ID       uint32
	IsKicked bool
}

func (p RemovePlayerPacket) Type() PlayerPacketType {
	return RemovePlayer
}

func (p RemovePlayerPacket) Marshal() ([]byte, error) {
	buf := make([]byte, 0)

	buf = append(buf, byte(RemovePlayer))
	buf = binary.LittleEndian.AppendUint32(buf, p.ID)
	if p.IsKicked {
		buf = append(buf, byte(1))
	} else {
		buf = append(buf, byte(0))
	}
	return buf, nil
}

type KickPlayerPacket struct{}

func (p KickPlayerPacket) Type() PlayerPacketType {
	return Kick
}

func (p KickPlayerPacket) Marshal() ([]byte, error) {
	buf := make([]byte, 0)

	buf = append(buf, byte(Kick))
	return buf, nil
}

type HostRecordPacket struct {
	SessionID   uint32
	NumOfFrames uint32
}

func (p HostRecordPacket) Type() HostPacketType {
	return HostRecord
}
