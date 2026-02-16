package game

import (
	"fmt"
	gamepackets "polyserver/game/packets"
	gametrack "polyserver/game/track"
	webrtc_session "polyserver/webrtc"
	"time"
)

type Player struct {
	Session                 *webrtc_session.PeerSession
	IsKicked                bool
	ID                      int
	Mods                    []string
	IsModsVanillaCompatible bool
	Nickname                string
	CountryCode             string
	ResetCounter            int
	CarStyle                string
	Record                  *Record
	Ping                    int
	PingIdCounter           int
	PingPackages            []PingPackage
	UnsentCarStates         []any // TODO: CAR STATE STUFF
}

type Record struct {
	numberOfFrames int
}

type PingPackage struct {
	pingId   int
	sentTime time.Time
}

func (player *Player) Send(packet gamepackets.PlayerPacket) error {
	data, err := packet.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal %s packet: %w", packet.Type(), err)
	}

	// Special handling for TrackChunk packets might be needed
	if packet.Type() == gamepackets.TrackChunk {
		return player.Session.ReliableDC.Send(data)
	}

	// For other packets, just send directly
	return player.Session.ReliableDC.Send(data)
}

func (player *Player) SendTrack(track gametrack.Track) error {
	// Send track ID
	trackId, _ := track.GetTrackID()
	if err := player.Send(gamepackets.TrackIDPacket{TrackID: trackId}); err != nil {
		return fmt.Errorf("failed to send track ID: %w", err)
	}

	// Split and send track data in chunks
	data := track.EncodedData
	for i := 0; i < len(data); i += 1024 {
		end := i + 1024
		if end > len(data) {
			end = len(data)
		}

		chunk := gamepackets.TrackChunkPacket{
			Data: data[i:end],
		}

		if err := player.Send(chunk); err != nil {
			return fmt.Errorf("failed to send chunk at offset %d: %w", i, err)
		}
	}

	return nil
}

func (player *Player) StartNewSession(session GameSession) {
	player.Send(gamepackets.NewSessionPacket{
		SessionID: session.SessionID,
		GameMode:  uint8(session.GameMode),
	})
}
