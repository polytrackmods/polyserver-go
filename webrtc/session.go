package webrtc

import (
	"encoding/json"

	"github.com/pion/webrtc/v4"
)

type PeerSession struct {
	SessionID string
	Peer      *webrtc.PeerConnection
	SendFunc  func([]byte) error
}

func NewPeerSession(
	sessionID string,
	offerSDP string,
	sendFunc func([]byte) error,
) (*PeerSession, error) {

	config := webrtc.Configuration{}

	peer, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	ps := &PeerSession{
		SessionID: sessionID,
		Peer:      peer,
		SendFunc:  sendFunc,
	}

	ps.setupICE()

	err = ps.handleOffer(offerSDP)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

func (ps *PeerSession) setupICE() {
	ps.Peer.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		payload := map[string]interface{}{
			"type":      "iceCandidate",
			"session":   ps.SessionID,
			"candidate": c.ToJSON(),
		}

		data, _ := json.Marshal(payload)
		ps.SendFunc(data)
	})
}

func (ps *PeerSession) handleOffer(offerSDP string) error {

	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerSDP,
	}

	err := ps.Peer.SetRemoteDescription(offer)
	if err != nil {
		return err
	}

	answer, err := ps.Peer.CreateAnswer(nil)
	if err != nil {
		return err
	}

	err = ps.Peer.SetLocalDescription(answer)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"type":    "answer",
		"session": ps.SessionID,
		"answer":  answer.SDP,
	}

	data, _ := json.Marshal(response)
	return ps.SendFunc(data)
}

func (ps *PeerSession) AddICECandidate(candidate webrtc.ICECandidateInit) error {
	return ps.Peer.AddICECandidate(candidate)
}
