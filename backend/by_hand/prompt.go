package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type (
	NewPromptMessage struct {
		PromptID  string `json:"prompt_id"`
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
	}

	PromptCreateRequest struct {
		Text string `json:"text"`
	}

	PromptCreateResponse struct {
		PromptID string `json:"prompt_id"`
	}

	PromptWaitRequest struct {
		LastPromptID string `json:"last_prompt_id"`
	}

	PromptWaitResponse struct {
		PromptID string `json:"prompt_id"`
		Text     string `json:"text"`
	}
)

func PromptHandler(ch chan NewPromptMessage) {
	for {
		p := <-ch
		fmt.Printf("Prompt: %v\n", p)
	}
}

func PromptCreateHandler(w http.ResponseWriter, r *http.Request) {
	var pcr PromptCreateRequest
	err := json.NewDecoder(r.Body).Decode(&pcr)
	if err != nil {
		fmt.Printf("Error creating prompt: %v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
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

	if participantID != s.OwnerID {
		fmt.Printf("Participant is not session owner: %v\n", participantID)
		http.Error(w, "only the session owner can create prompts", http.StatusUnauthorized)
		return
	}

	newPrompt := Prompt{
		ID:   uuid.New().String(),
		Text: pcr.Text,
	}

	s.Prompts = append(s.Prompts, newPrompt)
	prompts <- NewPromptMessage{
		PromptID:  newPrompt.ID,
		SessionID: s.ID,
		Text:      newPrompt.Text,
	}

	response := PromptCreateResponse{
		PromptID: newPrompt.ID,
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

func PromptWaitHandler(w http.ResponseWriter, r *http.Request) {
	var req PromptWaitRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("Error creating prompt: %v\n", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
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

	participantID := GetParticipantIDWithLock(sessionID, reqToken)
	if participantID == "" {
		fmt.Printf("Error getting participant: %v\n", err)
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	availablePrompt := GetPromptWithLock(sessionID, req.LastPromptID)
	if availablePrompt != nil {
		pwr := &PromptWaitResponse{
			PromptID: availablePrompt.ID,
			Text:     availablePrompt.Text,
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(pwr)
		if err != nil {
			fmt.Printf("Error encoding response: %v\n", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		return
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()
	var pwr *PromptWaitResponse
	for wait := true; wait; {
		select {
		case p := <-prompts:
			if p.SessionID != sessionID {
				break
			}
			pwr = &PromptWaitResponse{
				PromptID: p.PromptID,
				Text:     p.Text,
			}

			wait = false
			break
		case <-ctx.Done():
			wait = false
		}
	}

	if pwr == nil {
		fmt.Printf("No prompts found before timeout\n")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(pwr)
	if err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

}
