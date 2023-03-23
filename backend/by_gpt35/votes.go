package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Vote struct {
	SessionID string `json:"session_id"`
	PromptID  string `json:"prompt_id"`
	VoterID   string `json:"voter_id"`
	Vote      string `json:"vote"`
}

type VoteRequest struct {
	Vote string `json:"vote"`
}

type VoteSummary struct {
	ParticipantID string `json:"participant_id"`
	Vote          string `json:"vote"`
}

func watchVotes(w http.ResponseWriter, r *http.Request) {
	// Parse session and prompt IDs from URL path
	sessionID := mux.Vars(r)["session_id"]
	promptID := mux.Vars(r)["prompt_id"]

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

	// Read votes from file
	file, err := os.Open("votes.json")
	if err != nil {
		http.Error(w, "failed to open votes file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	votes := make([]VoteSummary, 0)
	for scanner.Scan() {
		var vote Vote
		if err := json.Unmarshal(scanner.Bytes(), &vote); err != nil {
			continue // skip invalid votes
		}
		if vote.SessionID == sessionID && vote.PromptID == promptID {
			voteSummary := VoteSummary{
				ParticipantID: vote.VoterID,
				Vote:          vote.Vote,
			}
			votes = append(votes, voteSummary)
		}
	}

	// Return response
	res := struct {
		Votes []VoteSummary `json:"votes"`
	}{
		Votes: votes,
	}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func vote(w http.ResponseWriter, r *http.Request) {
	// Parse session and prompt IDs from URL path
	sessionID := mux.Vars(r)["session_id"]
	promptID := mux.Vars(r)["prompt_id"]

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

	// Parse voter ID from cookie
	voterIDCookie, err := r.Cookie("voter_id")
	if err != nil {
		http.Error(w, "voter ID cookie not found", http.StatusBadRequest)
		return
	}
	var voterID string
	if err = cookieStore.Decode("voter_id", voterIDCookie.Value, &voterID); err != nil {
		http.Error(w, "invalid voter ID cookie", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req VoteRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate vote value
	validVotes := []string{"0", "1", "2", "3", "5", "8", "13", "20", "40", "100", "?", "☕️"}
	validVote := false
	for _, vote := range validVotes {
		if req.Vote == vote {
			validVote = true
			break
		}
	}
	if !validVote {
		http.Error(w, "invalid vote value", http.StatusBadRequest)
		return
	}

	// Record vote in file
	vote := Vote{
		SessionID: sessionID,
		PromptID:  promptID,
		VoterID:   voterID,
		Vote:      req.Vote,
	}
	file, err := os.OpenFile("votes.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "failed to open votes file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	jsonBytes, err := json.Marshal(vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonBytes = append(jsonBytes, '\n')
	_, err = file.Write(jsonBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	res := struct{}{}
	jsonBytes, err = json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}
