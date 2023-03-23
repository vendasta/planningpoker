package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
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
	seededRand  = rand.New(rand.NewSource(time.Now().UnixNano()))
	cookieStore = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/create_session", createSession).Methods("POST")
	r.HandleFunc("/session/{session_id}/join", joinSession).Methods("POST")
	r.HandleFunc("/session/{session_id}/prompt/wait", promptWait).Methods("POST")
	r.HandleFunc("/session/{session_id}/prompt/create", createPrompt).Methods("POST")
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/vote", vote).Methods("POST")
	r.HandleFunc("/session/{session_id}/prompt/{prompt_id}/watch", watchVotes).Methods("POST")
	r.HandleFunc("/session/{session_id}/close`", closeSession).Methods("POST")

	http.ListenAndServe(":8080", r)
}

func getCookie(r *http.Request) (*http.Cookie, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func getSession(r *http.Request) (Session, error) {
	cookie, err := getCookie(r)
	if err != nil {
		return Session{}, err
	}

	var session Session
	err = cookieStore.Decode("session_id", cookie.Value, &session)
	if err != nil {
		return Session{}, err
	}
	return session, nil
}

func setSession(w http.ResponseWriter, session Session) {
	encoded, err := cookieStore.Encode("session_id", session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = Charset[seededRand.Intn(len(Charset))]
	}
	return string(b)
}
