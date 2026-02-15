package signaling

import (
	"encoding/json"
	"fmt"
	"log"
	"polyserver/config"
	"polyserver/webrtc"

	"github.com/gorilla/websocket"
)

type Server struct {
	Conn          *websocket.Conn
	CurrentInvite string
	Sessions      map[string]*webrtc.PeerSession
}

func NewServer() *Server {
	return &Server{
		Sessions: make(map[string]*webrtc.PeerSession),
	}
}

func (s *Server) Connect() error {
	if s.Conn != nil {
		s.Conn.Close()
	}

	conn, _, err := websocket.DefaultDialer.Dial(config.WebsocketUrl, nil)
	if err != nil {
		return err
	}

	s.Conn = conn
	return nil
}

func (s *Server) RegenerateInvite() error {
	if err := s.Connect(); err != nil {
		return err
	}

	go s.Start()

	return s.CreateInvite()
}

func (s *Server) CreateInvite() error {
	if s.Conn == nil {
		return fmt.Errorf("not connected")
	}

	payload := map[string]interface{}{
		"version": config.PolyVersion,
		"type":    "createInvite",
	}

	data, _ := json.Marshal(payload)

	return s.Conn.WriteMessage(websocket.TextMessage, data)
}

func (s *Server) Start() {
	for {
		_, message, err := s.Conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			s.RegenerateInvite()
			return
		}

		s.route(message)
	}
}

func (s *Server) handleCreateInvite(p CreateInviteResponse) {
	s.CurrentInvite = p.InviteCode
	log.Println("Invite code:", p.InviteCode)
}

func (s *Server) handleJoinInvite(p JoinInvite) {
	log.Println("User is joining:", p.Nickname)

	session, err := webrtc.NewPeerSession(
		p.Session,
		p.Offer,
		s.send,
	)
	if err != nil {
		log.Println("failed to create session:", err)
		return
	}

	s.Sessions[p.Session] = session
	log.Println("Created session:", p.Session)

	joinPacket, _ := json.Marshal(AcceptJoinPacket{
		Type:                    "acceptJoin",
		Version:                 config.PolyVersion,
		Session:                 p.Session,
		Mods:                    config.LoadedMods,
		IsModsVanillaCompatible: config.AcceptVanillaClients,
	})

	s.send([]byte(joinPacket))
}

func (s *Server) handleICE(p IceCandidatePacket) {
	log.Println("Received ICE candidate.")
	session, ok := s.Sessions[p.Session]
	if !ok {
		log.Println("unknown session:", p.Session)
		return
	}

	err := session.AddICECandidate(p.Candidate)
	if err != nil {
		log.Println("failed to add ICE:", err)
	}
}

func (s *Server) send(data []byte) error {
	return s.Conn.WriteMessage(1, data)
}
