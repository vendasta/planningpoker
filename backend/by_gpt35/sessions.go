package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Participant struct {
	ID    string
	Token string
}

type CreateSessionRequest struct {
	ParticipantID string `json:"participant_id"`
}

type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
	Token     string `json:"token"`
}

type JoinSessionRequest struct {
	ParticipantID string `json:"participant_id"`
}

type JoinSessionResponse struct {
	Token string `json:"token"`
}

func closeSession(w http.ResponseWriter, r *http.Request) {
	// Parse session ID from URL path
	sessionID := mux.Vars(r)["session_id"]

	err := store.deleteSession(sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}

	// Return response
	res := struct{}{}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func joinSession(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}

	// Parse session ID from URL path
	sessionID := mux.Vars(r)["session_id"]

	// Read session ID from cookie
	session, err := getSession(sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req JoinSessionRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate session ID for participant
	token := generateRandomString(TokenLength)
	participants := session.Participants
	participants[req.ParticipantID] = Participant{
		ID:    req.ParticipantID,
		Token: token,
	}

	// update the participants in the session
	session.Participants = participants

	err = setSession(session)
	if err != nil {
		http.Error(w, "unable to create new session", http.StatusInternalServerError)
		return
	}

	// Create response
	res := JoinSessionResponse{
		Token: token,
	}

	// Encode response as JSON and write to response
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func createSession(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}

	// Parse request body
	var req CreateSessionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate token
	token := generateRandomString(TokenLength)

	// Get session ID from cookie or generate a new one
	sessionID := generateRandomString(SessionIDLength)
	session := &Session{
		OwnerID: token,
		ID:      sessionID,
		Participants: map[string]Participant{
			req.ParticipantID: {
				ID:    req.ParticipantID,
				Token: token,
			},
		},
	}

	err = setSession(session)
	if err != nil {
		http.Error(w, "unable to create new session", http.StatusInternalServerError)
		return
	}

	// Create response
	res := CreateSessionResponse{
		SessionID: session.ID,
		Token:     token,
	}

	// Encode response as JSON and write to response
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
