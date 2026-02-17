package signaling

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"polyserver/config"
	webrtc_session "polyserver/webrtc"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type WebRTCServer struct {
	Conn          *websocket.Conn
	ConnLock      sync.Mutex
	ICEUrls       []string
	CurrentInvite string
	SessionLock   sync.Mutex
	Sessions      map[string]*webrtc_session.PeerSession
	ClientCount   uint32

	OnOpen  func(joinPacket JoinInvite, session *webrtc_session.PeerSession)
	OnClose func(sessionId string)
}

func NewServer() *WebRTCServer {

	resp, err := http.Get(config.IceFetchUrl)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var IceServers []IceServerResponse
	err1 := json.Unmarshal(body, &IceServers)
	if err1 != nil {
		log.Fatalln("Failed: Invalid ICE server response.")
	}
	var IceUrls []string
	var numOfUrls = 0
	for _, urls := range IceServers {
		IceUrls = append(IceUrls, urls.Urls)
		numOfUrls++
	}
	log.Println("Got " + strconv.Itoa(numOfUrls) + " ICE URLs")
	return &WebRTCServer{
		Sessions:    make(map[string]*webrtc_session.PeerSession),
		ClientCount: 1,
		ICEUrls:     IceUrls,
	}
}

func (s *WebRTCServer) Connect() error {
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

func (s *WebRTCServer) RegenerateInvite() error {
	if err := s.Connect(); err != nil {
		return err
	}

	go s.Start()

	return s.CreateInvite()
}

func (s *WebRTCServer) CreateInvite() error {
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

func (s *WebRTCServer) Start() {
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

func (s *WebRTCServer) handleCreateInvite(p CreateInviteResponse) {
	s.CurrentInvite = p.InviteCode
	log.Println("Invite code:", p.InviteCode)
}

func (s *WebRTCServer) onConnectionClosed(sessionId string) {
	s.SessionLock.Lock()
	defer s.SessionLock.Unlock()
	for k := range s.Sessions {
		if k == sessionId {
			log.Printf("Removing %v from Sessions...\n", sessionId)
			s.OnClose(sessionId)
			delete(s.Sessions, k)
			break
		}
	}
}

func (s *WebRTCServer) handleJoinInvite(p JoinInvite) {
	log.Println("User is joining:", p.Nickname)

	session, answer, err := webrtc_session.NewPeerSession(
		p.Session,
		p.Offer,
		s.OnIceCandidateServer,
		s.ICEUrls,
		s.onConnectionClosed,
	)
	if err != nil {
		log.Println("failed to create session:", err)
		return
	}

	s.SessionLock.Lock()
	s.Sessions[p.Session] = session
	s.SessionLock.Unlock()

	session.ReliableDC.OnOpen(func() {
		s.OnOpen(p, session)
	})

	log.Println("Created session:", p.Session)
	joinPacket, _ := json.Marshal(AcceptJoinPacket{
		Type:                    "acceptJoin",
		Version:                 config.PolyVersion,
		Session:                 p.Session,
		Mods:                    config.LoadedMods,
		IsModsVanillaCompatible: config.AcceptVanillaClients,
		CliendId:                s.ClientCount,
		Answer:                  answer,
	})
	s.ClientCount++
	log.Println("Answering...")

	s.send([]byte(joinPacket))
}

func (s *WebRTCServer) handleICE(p IceCandidateResponse) {
	s.SessionLock.Lock()
	session, ok := s.Sessions[p.Session]
	s.SessionLock.Unlock()

	if !ok {
		log.Println("unknown session:", p.Session)
		return
	}

	err := session.AddICECandidate(p.Candidate)
	if err != nil {
		log.Println("failed to add ICE:", err)
	}
}

func (s *WebRTCServer) OnIceCandidateServer(candidate []byte, session string) error {
	var iceCandidate IceCandidate
	err := json.Unmarshal(candidate, &iceCandidate)
	if err != nil {
		return err
	}
	icePacket, err := json.Marshal(IceCandidatePacket{
		Type:      "iceCandidate",
		Candidate: iceCandidate,
		Version:   config.PolyVersion,
		Session:   session,
	})
	if err != nil {
		return err
	}
	return s.send(icePacket)
}

func (s *WebRTCServer) send(data []byte) error {
	s.ConnLock.Lock()
	defer s.ConnLock.Unlock()
	return s.Conn.WriteMessage(1, data)
}
