package main

import (
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
)

var sessions map[string]Session
var sessionLock sync.Mutex

var prompts chan NewPromptMessage

func main() {
	prompts = make(chan NewPromptMessage)
	go PromptHandler(prompts)
	r := mux.NewRouter()

	r.HandleFunc("/session/create", SessionCreateHandler)
	r.HandleFunc("/session/{session_id}/join", SessionJoinHandler)
	r.HandleFunc("/session/{session_id}/close", SessionCloseHandler)

	r.HandleFunc("/session/{session_id}/prompt/create", PromptCreateHandler)
	r.HandleFunc("/session/{session_id}/prompt/wait", PromptWaitHandler)

	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote/submit", VoteSubmitHandler)
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote/{vote_id}/watch", VoteWatchHandler)
}

//----------------------------------------------

//----------------------------------------------

func VoteSubmitHandler(w http.ResponseWriter, r *http.Request) {

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

	ready := false
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
