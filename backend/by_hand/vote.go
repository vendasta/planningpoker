package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

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
	for k, prompt := range s.Prompts {
		if prompt.ID == promptID {
			skip := false
			for k, vote := range prompt.Votes {
				if vote.ParticipantID == participantID {
					vote.Vote = req.Vote
					prompt.Votes[k] = vote
					skip = true
					break
				}
			}
			if !skip {
				prompt.Votes = append(prompt.Votes, Vote{
					ParticipantID: participantID,
					PromptID:      promptID,
					Vote:          req.Vote,
				})
			}
			s.Prompts[k] = prompt
			voted = true
			break
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
