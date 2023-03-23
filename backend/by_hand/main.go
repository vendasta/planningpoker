package main

import (
	"github.com/gorilla/mux"
	"net/http"
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

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/session/create", SessionCreateHandler)
	r.HandleFunc("/session/{session_id}/join", SessionJoinHandler)
	r.HandleFunc("/session/{session_id}/close", SessionCloseHandler)

	r.HandleFunc("/prompt/create", PromptCreateHandler)
	r.HandleFunc("/prompt/{prompt_id}/wait", PromptWaitHandler)

	r.HandleFunc("/vote/submit", VoteSubmitHandler)
	r.HandleFunc("/vote/{vote_id}/watch", VoteWatchHandler)
}

//----------------------------------------------

func PromptCreateHandler(w http.ResponseWriter, r *http.Request) {

}

func PromptWaitHandler(w http.ResponseWriter, r *http.Request) {

}

//----------------------------------------------

func VoteSubmitHandler(w http.ResponseWriter, r *http.Request) {

}

func VoteWatchHandler(w http.ResponseWriter, r *http.Request) {

}

//----------------------------------------------
