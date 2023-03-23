package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	SessionCreateRequest struct {
		ParticipantID string `json:"participant_id"`
	}

	SessionCreateResponse struct {
		SessionID string `json:"session_id"`
		Token     string `json:"token"`
	}

	SessionJoinRequest struct {
		ParticipantID string `json:"participant_id"`
	}

	SessionJoinResponse struct {
		Token string `json:"token"`
	}
)

func SessionCreateHandler(w http.ResponseWriter, r *http.Request) {
	var scr SessionCreateRequest
	err := json.NewDecoder(r.Body).Decode(&scr)
	if err != nil {
		fmt.Printf("Error decoding request: %v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var p Participant
	var s Session

	sessionLock.Lock()
	defer sessionLock.Unlock()
	for {
		id := GenerateID()
		if _, ok := sessions[id]; !ok {
			p = Participant{
				ID:    scr.ParticipantID,
				Token: uuid.New().String(),
			}

			s = Session{
				ID:           id,
				Participants: make(map[string]Participant),
				Prompts:      make([]Prompt, 0),
				OwnerID:      p.ID,
			}
			s.Participants[p.ID] = p

			sessions[id] = s
			break
		}
	}

	response := SessionCreateResponse{
		SessionID: s.ID,
		Token:     p.Token,
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func SessionJoinHandler(w http.ResponseWriter, r *http.Request) {
	var req SessionJoinRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Error decoding request: %v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	sessionLock.Lock()
	defer sessionLock.Unlock()

	s, ok := sessions[sessionID]
	if !ok {
		fmt.Printf("Attempted to join non-existent session: %v\n", sessionID)
		http.Error(w, "invalid session", http.StatusNotFound)
		return
	}

	_, ok = s.Participants[req.ParticipantID]
	if ok {
		fmt.Printf("Attempted to join session with duplicate participant ID: %v\n", req.ParticipantID)
		http.Error(w, "duplicate participant ID", http.StatusConflict)
		return
	}

	p := Participant{
		ID:    req.ParticipantID,
		Token: uuid.New().String(),
	}
	s.Participants[p.ID] = p

	response := SessionJoinResponse{
		Token: p.Token,
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}

func SessionCloseHandler(w http.ResponseWriter, r *http.Request) {
	reqToken, err := GetToken(r)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]

	sessionLock.Lock()
	defer sessionLock.Unlock()

	s, ok := sessions[sessionID]
	if !ok {
		fmt.Printf("Attempted to close non-existent session: %v\n", sessionID)
		http.Error(w, "invalid session", http.StatusNotFound)
		return
	}

	participantID := GetParticipantID(s, reqToken)
	if participantID == "" {
		fmt.Printf("Attempted to close session with invalid token: %v\n", reqToken)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	if participantID != s.OwnerID {
		fmt.Printf("Attempted to close session with non-owner token: %v\n", reqToken)
		http.Error(w, "only the session creator can close it", http.StatusUnauthorized)
		return
	}

	delete(sessions, sessionID)
	w.WriteHeader(http.StatusNoContent)
}
