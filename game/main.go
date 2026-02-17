package game

import (
	"fmt"
	"log"

	gamepackets "polyserver/game/packets"
	"polyserver/signaling"
	webrtc_session "polyserver/webrtc"
	"sync"
	"time"
)

type GameServer struct {
	SignalingServer *signaling.WebRTCServer
	Players         []*Player
	playersLock     sync.Mutex
	Factory         gamepackets.PacketFactory
	GameSession     *GameSession
}

type GameMode uint8

const (
	Casual GameMode = iota
	Competitive
)

func (gm GameMode) String() string {
	switch gm {
	case Casual:
		return "Casual"
	case Competitive:
		return "Competitive"
	default:
		return fmt.Sprintf("Unknown(%d)", gm)
	}
}

func NewServer(signalingServer *signaling.WebRTCServer) *GameServer {
	server := &GameServer{
		SignalingServer: signalingServer,
		Players:         make([]*Player, 0),
		Factory:         gamepackets.PacketFactory{},
		GameSession:     nil,
	}
	signalingServer.OnOpen = server.onPlayerJoin
	signalingServer.OnClose = server.onPlayerDisconnect

	schedule(server.sendPings, time.Second)

	return server
}

func (server *GameServer) onPlayerJoin(p signaling.JoinInvite, session *webrtc_session.PeerSession) {
	log.Println("Creating player " + p.Nickname)
	carStyle, err := gamepackets.FromBase64String(p.CarStyle)
	if err != nil {
		carStyle = gamepackets.DefaultCarStyle()
		log.Println("Failed fromBase64String: " + err.Error())
	}
	log.Printf("Car Style: %v\n", carStyle)
	newPlayer := NewPlayer(&Player{
		Server:                  server,
		Session:                 session,
		IsKicked:                false,
		ID:                      server.SignalingServer.ClientCount - 1,
		Mods:                    p.Mods,
		IsModsVanillaCompatible: p.IsModsVanillaCompatible,
		Nickname:                p.Nickname,
		CountryCode:             p.CountryCode,
		ResetCounter:            0,
		CarStyle:                carStyle,
		NumberOfFrames:          nil,
		Ping:                    0,
		PingIdCounter:           0,
		PingPackages:            make([]PingPackage, 0),
		UnsentCarStates:         make([]gamepackets.CarState, 0),
	})

	newPlayer.Send(gamepackets.EndSessionPacket{})
	newPlayer.SendTrack()
	newPlayer.StartNewSession()
	for _, player := range server.Players {
		log.Printf("Sending player %s to %s", player.Nickname, newPlayer.Nickname)
		newPlayer.SendPlayerUpdate(player)
	}
	server.propagateUpdate(newPlayer)
	server.Players = append(server.Players, newPlayer)
}

func (server *GameServer) onPlayerDisconnect(sessionId string) {

	for index, player := range server.Players {
		if player.Session.SessionID == sessionId {
			log.Println("Removing player " + player.Nickname)
			server.playersLock.Lock()
			server.Players = append(server.Players[:index], server.Players[index+1:]...)
			server.playersLock.Unlock()
		}
	}
}

func schedule(f func(), interval time.Duration) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		// Loop indefinitely, running the function every time a tick is received
		for range ticker.C {
			f()
		}
	}()
	return ticker
}

func (server *GameServer) sendPings() {
	for _, player := range server.Players {
		// log.Printf("Sending ping(%v): %v\n", player.Ping, player.PingIdCounter+1)
		player.SendPing()
	}
	server.sendPingDatas()
}

func (server *GameServer) sendPingDatas() {
	pings := server.getPlayerPings()

	for _, player := range server.Players {
		player.SendUnreliable(gamepackets.PingDataPacket{
			HostID:      0,
			PlayerPings: pings,
		})
	}
}

func (server *GameServer) getPlayerPings() []gamepackets.PlayerPing {
	pings := make([]gamepackets.PlayerPing, len(server.Players))
	for _, player := range server.Players {
		pings = append(pings, gamepackets.PlayerPing{
			PlayerID: player.ID,
			Ping:     uint16(player.Ping),
		})
	}
	return pings
}

func (server *GameServer) propagateUpdate(p *Player) {
	for _, player := range server.Players {
		log.Printf("Sending player %s to %s", p.Nickname, player.Nickname)
		player.SendPlayerUpdate(p)
	}
}
