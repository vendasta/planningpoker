package main

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Prompt struct {
	SessionID string `json:"session_id"`
	PromptID  string `json:"prompt_id"`
	Prompt    string `json:"prompt"`
}

type PromptWaitRequest struct {
	LastPromptID string `json:"last_prompt_id"`
}

type PromptWaitResponse struct {
	PromptID string `json:"prompt_id"`
	Prompt   string `json:"prompt"`
}

type CreatePromptRequest struct {
	Prompt string `json:"prompt"`
}

type CreatePromptResponse struct {
	PromptID string
}

func createPrompt(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}

	token, err := getBearerTokenFromHTTP(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
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

	// Verify that token matches session owner ID
	if token != session.OwnerID {
		http.Error(w, "token does not match session owner ID", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var req CreatePromptRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate new prompt ID
	promptID := generateRandomString(8)

	// Create prompt object
	prompt := Prompt{
		SessionID: sessionID,
		PromptID:  promptID,
		Prompt:    req.Prompt,
	}

	// Write prompt to file
	file, err := os.OpenFile("prompts.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, "failed to open prompts file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	jsonBytes, err := json.Marshal(prompt)
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
	res := CreatePromptResponse{
		PromptID: promptID,
	}
	jsonBytes, err = json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func promptWait(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}

	// Parse session ID from URL path
	sessionID := mux.Vars(r)["session_id"]

	// Read session ID from cookie
	_, err := getSession(sessionID)
	if err != nil {
		http.Error(w, "session not found", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req PromptWaitRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Read prompts from file
	file, err := os.Open("prompts.json")
	if err != nil {
		http.Error(w, "failed to open prompts file", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var prompt Prompt
		if err := json.Unmarshal(scanner.Bytes(), &prompt); err != nil {
			continue // skip invalid prompts
		}
		if prompt.SessionID == sessionID {
			if req.LastPromptID == "" {
				// Return first prompt
				res := PromptWaitResponse{
					PromptID: prompt.PromptID,
					Prompt:   prompt.Prompt,
				}
				jsonBytes, err := json.Marshal(res)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
				return
			}

			// Return next prompt
			for scanner.Scan() {
				var nextPrompt Prompt
				if err := json.Unmarshal(scanner.Bytes(), &nextPrompt); err != nil {
					continue // skip invalid prompts
				}
				if nextPrompt.SessionID == sessionID {
					res := PromptWaitResponse{
						PromptID: nextPrompt.PromptID,
						Prompt:   nextPrompt.Prompt,
					}
					jsonBytes, err := json.Marshal(res)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					w.Write(jsonBytes)
					return
				}
			}
		}
	}

	// No prompt is available, return empty response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{}"))
}
