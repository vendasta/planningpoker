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
	initialize()

	r := CreateHandler()

	fmt.Printf("serving on port 9000\n")
	http.ListenAndServe(":9000", r)
}

//----------------------------------------------

func initialize() {
	sessions = make(map[string]Session)
	prompts = make(chan NewPromptMessage)
	go PromptHandler(prompts)
}

func CreateHandler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/session/create", SessionCreateHandler)
	r.HandleFunc("/session/{session_id}/join", SessionJoinHandler)
	r.HandleFunc("/session/{session_id}/close", SessionCloseHandler)

	r.HandleFunc("/session/{session_id}/prompt/create", PromptCreateHandler)
	r.HandleFunc("/session/{session_id}/prompt/wait", PromptWaitHandler)

	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote/submit", VoteSubmitHandler)
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote", VoteGetHandler)
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote/watch", VoteWatchHandler)
	return r
}

//----------------------------------------------

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
