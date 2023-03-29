package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type Session struct {
	ID           string
	Participants map[string]Participant
	Prompts      []Prompt
	OwnerID      string
}

const (
	SessionIDLength = 6
	TokenLength     = 6
	Charset         = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	store      = NewInMemorySessionStore()
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/session/create", createSession).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{session_id}/join", joinSession).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{session_id}/prompt/wait", promptWait).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{session_id}/prompt/create", createPrompt).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote", vote).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/watch", watchVotes).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	r.HandleFunc("/session/{session_id}/close", closeSession).Methods(http.MethodPost, http.MethodOptions)

	// Allow CORS for all routes
	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":8080", r)
}

func getSession(sessionID string) (*Session, error) {
	s, err := store.getSession(sessionID)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func setSession(session *Session) error {
	return store.setSession(session)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[seededRand.Intn(len(Charset))]
	}
	return string(b)
}

// getBearerTokenFromHTTP extracts the bearer token out of the authorization header
func getBearerTokenFromHTTP(r *http.Request) (string, error) {
	auth := r.Header.Get("authorization")
	if auth == "" {
		return "", fmt.Errorf("no authorization header")
	}

	pieces := strings.Split(auth, " ")
	if len(pieces) != 2 || pieces[0] != "Bearer" {
		return "", fmt.Errorf("no authorization header")
	}
	return pieces[1], nil
}

func handleCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With,X-Grpc-Web,X-User-Agent")

	if r.Method == http.MethodOptions {
		w.Header().Add("Access-Control-Max-Age", "1728000")
		w.Header().Add("Content-Type", "text/plain; charset=UTF-8")
		w.Header().Add("Content-Length", "0")
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}
