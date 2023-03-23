package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
	"sync"
)

type (
	Vote struct {
		ParticipantID string
		PromptID      string
		Vote          string
	}

	Prompt struct {
		ID    string
		Text  string
		Votes []Vote
	}

	Participant struct {
		ID    string
		Token string
	}

	Session struct {
		ID           string
		Participants map[string]Participant
		Prompts      []Prompt
		OwnerID      string
	}

	VoteSubmitRequest struct {
		Vote string `json:"vote"`
	}

	VoteGetResponse struct {
		Votes []Vote `json:"vote"`
	}

	VoteWatchResponse struct {
		Votes []Vote `json:"vote"`
	}
)

var sessions map[string]Session
var sessionLock sync.Mutex

var prompts chan NewPromptMessage

func main() {
	sessions = make(map[string]Session)
	prompts = make(chan NewPromptMessage)
	go PromptHandler(prompts)
	r := mux.NewRouter()

	r.HandleFunc("/session/create", SessionCreateHandler)
	r.HandleFunc("/session/{session_id}/join", SessionJoinHandler)
	r.HandleFunc("/session/{session_id}/close", SessionCloseHandler)

	r.HandleFunc("/session/{session_id}/prompt/create", PromptCreateHandler)
	r.HandleFunc("/session/{session_id}/prompt/wait", PromptWaitHandler)

	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote/submit", VoteSubmitHandler)
	r.HandleFunc("/sesssion/{session_id}/prompt/{prompt_id}/vote/{vote_id}/get", VoteGetHandler)
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote/{vote_id}/watch", VoteWatchHandler)

	fmt.Printf("serving on port 9000\n")
	http.ListenAndServe(":9000", r)
}

//----------------------------------------------

//----------------------------------------------

func VoteSubmitHandler(w http.ResponseWriter, r *http.Request) {
	var req VoteSubmitRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Error submitting vote: %v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	switch req.Vote {
	case "0", "1/2", "1", "2", "3", "5", "8", "13", "20", "40", "100", "?", "☕️":
	default:
		fmt.Printf("Invalid vote: %v\n", req.Vote)
		http.Error(w, "invalid vote", http.StatusBadRequest)
		return
	}

	reqToken, err := GetToken(r)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	promptID := vars["prompt_id"]

	sessionLock.Lock()
	defer sessionLock.Unlock()

	s, ok := sessions[sessionID]
	if !ok {
		fmt.Printf("Session not found: %v\n", sessionID)
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	participantID := GetParticipantID(s, reqToken)
	if participantID == "" {
		fmt.Printf("Participant not found for token: %v\n", reqToken)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	voted := false
	for _, prompt := range s.Prompts {
		if prompt.ID == promptID {
			prompt.Votes = append(prompt.Votes, Vote{
				ParticipantID: participantID,
				PromptID:      promptID,
				Vote:          req.Vote,
			})
			voted = true
		}
	}

	if !voted {
		fmt.Printf("Prompt not found: %v\n", promptID)
		http.Error(w, "prompt not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func VoteGetHandler(w http.ResponseWriter, r *http.Request) {
	reqToken, err := GetToken(r)
	if err != nil {
		fmt.Printf("Error getting token: %v\n", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	sessionID := vars["session_id"]
	promptID := vars["prompt_id"]

	sessionLock.Lock()
	defer sessionLock.Unlock()

	s, ok := sessions[sessionID]
	if !ok {
		fmt.Printf("Session not found: %v\n", sessionID)
		http.Error(w, "session not found", http.StatusNotFound)
		return
	}

	participantID := GetParticipantID(s, reqToken)
	if participantID == "" {
		fmt.Printf("Participant not found for token: %v\n", reqToken)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	for _, prompt := range s.Prompts {
		if prompt.ID == promptID {
			response := VoteGetResponse{Votes: prompt.Votes}
			w.Header().Set("content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Printf("Error encoding response: %v\n", err)
				return
			}
			return
		}
	}

	fmt.Printf("Prompt not found: %v\n", promptID)
	http.Error(w, "prompt not found", http.StatusNotFound)
}

func VoteWatchHandler(w http.ResponseWriter, r *http.Request) {

}

//----------------------------------------------

func GetToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("authorization")
	regex := regexp.MustCompile(`^Bearer (.*)$`)
	matches := regex.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", fmt.Errorf("Invalid authorization header: %v", authHeader)
	}
	return matches[1], nil
}

func GetParticipantID(s Session, token string) string {
	for _, p := range s.Participants {
		if p.Token == token {
			return p.ID
		}
	}

	return ""
}

func GetParticipantIDWithLock(sid string, token string) string {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	s, ok := sessions[sid]
	if !ok {
		return ""
	}

	return GetParticipantID(s, token)
}

func GetPromptWithLock(sessionID, LastPromptID string) *Prompt {
	sessionLock.Lock()
	defer sessionLock.Unlock()

	s, ok := sessions[sessionID]
	if !ok {
		return nil
	}

	ready := LastPromptID == ""
	for _, p := range s.Prompts {
		if ready {
			return &p
		}
		if p.ID == LastPromptID {
			ready = true
		}
	}

	return nil

}
