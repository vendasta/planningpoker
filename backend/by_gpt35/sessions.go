package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

	// Get session ID from cookie
	sessionIDCookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "session ID cookie not found", http.StatusBadRequest)
		return
	}
	var cookieSessionID string
	if err = cookieStore.Decode("session_id", sessionIDCookie.Value, &cookieSessionID); err != nil {
		http.Error(w, "invalid session ID cookie", http.StatusBadRequest)
		return
	}

	// Verify that session ID in cookie matches session ID in URL path
	if cookieSessionID != sessionID {
		http.Error(w, "session ID in URL path does not match session ID in cookie", http.StatusBadRequest)
		return
	}

	// Delete session ID cookie
	sessionIDCookie.MaxAge = -1
	http.SetCookie(w, sessionIDCookie)

	// Delete voter ID cookie
	voterIDCookie, err := r.Cookie("voter_id")
	if err == nil {
		voterIDCookie.MaxAge = -1
		http.SetCookie(w, voterIDCookie)
	}

	// Delete prompt and vote files for this session
	promptFilepath := fmt.Sprintf("prompts-%s.json", sessionID)
	voteFilepath := fmt.Sprintf("votes-%s.json", sessionID)
	os.Remove(promptFilepath)
	os.Remove(voteFilepath)

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
	// Parse session ID from URL path
	sessionID := mux.Vars(r)["session_id"]

	// Read session ID from cookie
	session, err := getSession(r)
	if err != nil {
		http.Error(w, "session ID cookie not found", http.StatusBadRequest)
		return
	}

	// Verify that session ID in cookie matches session ID in URL path
	if session.ID != sessionID {
		http.Error(w, "session ID in URL path does not match session ID in cookie", http.StatusBadRequest)
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
	session = Session{
		ID: sessionID,
		Participants: map[string]Participant{
			req.ParticipantID: {
				ID:    req.ParticipantID,
				Token: token,
			},
		},
	}

	// Create response
	res := JoinSessionResponse{
		Token: generateRandomString(TokenLength),
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
	session, err := getSession(r)
	if err != nil {
		sessionID := generateRandomString(SessionIDLength)
		session = Session{
			OwnerID: token,
			ID:      sessionID,
			Participants: map[string]Participant{
				req.ParticipantID: {
					ID:    req.ParticipantID,
					Token: token,
				},
			},
		}
	}

	// write the session to the cookie
	setSession(w, session)

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
