package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func createTestSession(t *testing.T, testServer *httptest.Server) (string, string) {
	scr := SessionCreateRequest{ParticipantID: "Test User"}
	buf := bytes.NewBufferString("")
	err := json.NewEncoder(buf).Encode(scr)
	if err != nil {
		t.Fatalf("Error encoding request: %v", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/session/create", testServer.URL), buf)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Error sending request: %s", resp.Status)
	}

	var responseBody SessionCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	return responseBody.SessionID, responseBody.Token
}

func joinTestSession(t *testing.T, testServer *httptest.Server, sessionID string) string {
	sjr := SessionJoinRequest{ParticipantID: "Test User2"}
	buf := bytes.NewBufferString("")
	err := json.NewEncoder(buf).Encode(sjr)
	if err != nil {
		t.Fatalf("Error encoding request: %v", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/session/%s/join", testServer.URL, sessionID), buf)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Error sending request: %s", resp.Status)
	}

	var responseBody SessionJoinResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	return responseBody.Token
}

func createTestPrompt(t *testing.T, testServer *httptest.Server, sessionID, token, prompt string) string {
	pr := PromptCreateRequest{Text: prompt}
	buf := bytes.NewBufferString("")
	err := json.NewEncoder(buf).Encode(pr)
	if err != nil {
		t.Fatalf("Error encoding request: %v", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/session/%s/prompt/create", testServer.URL, sessionID), buf)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Error sending request: %s", resp.Status)
	}

	var responseBody PromptCreateResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	return responseBody.PromptID
}

func createTestVote(t *testing.T, testServer *httptest.Server, sessionID, token, promptID, vote string) {
	vr := VoteSubmitRequest{Vote: vote}
	buf := bytes.NewBufferString("")
	err := json.NewEncoder(buf).Encode(vr)
	if err != nil {
		t.Fatalf("Error encoding request: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/session/%s/prompt/%s/vote/submit", testServer.URL, sessionID, promptID), buf)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Error sending request: %s", resp.Status)
	}
}

func getTestVote(t *testing.T, testServer *httptest.Server, sessionID, token, promptID string) []Vote {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/session/%s/prompt/%s/vote", testServer.URL, sessionID, promptID), nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Error sending request: %s", resp.Status)
	}

	var responseBody VoteGetResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	return responseBody.Votes
}

func testWaitForPrompt(t *testing.T, testServer *httptest.Server, sessionID, token, lastPrompt string) string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/session/%s/prompt/wait?last_prompt_id=%s", testServer.URL, sessionID, lastPrompt), nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := testServer.Client().Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if resp.StatusCode == http.StatusNoContent {
		return ""
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Error sending request: %s", resp.Status)
	}

	var responseBody PromptWaitResponse
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	return responseBody.PromptID
}

func TestUpToVote(t *testing.T) {
	initialize()
	testServer := httptest.NewServer(CreateHandler())

	sessionID, token := createTestSession(t, testServer)
	promptID := createTestPrompt(t, testServer, sessionID, token, "Test Prompt?")
	createTestVote(t, testServer, sessionID, token, promptID, "1")
	votes := getTestVote(t, testServer, sessionID, token, promptID)

	if len(votes) != 1 {
		t.Fatalf("Expected 1 vote, got %d", len(votes))
	}

	if votes[0].Vote != "1" {
		t.Fatalf("Expected vote to be 1, got %s", votes[0].Vote)
	}

	createTestVote(t, testServer, sessionID, token, promptID, "2")
	votes = getTestVote(t, testServer, sessionID, token, promptID)

	if len(votes) != 1 {
		t.Fatalf("Expected 1 vote, got %d", len(votes))
	}

	if votes[0].Vote != "2" {
		t.Fatalf("Expected vote to be 2, got %s", votes[0].Vote)
	}
}

func TestTwoVotes(t *testing.T) {
	initialize()
	testServer := httptest.NewServer(CreateHandler())

	sessionID1, token1 := createTestSession(t, testServer)
	token2 := joinTestSession(t, testServer, sessionID1)
	promptID := createTestPrompt(t, testServer, sessionID1, token1, "Test Prompt?")

	createTestVote(t, testServer, sessionID1, token1, promptID, "1")
	createTestVote(t, testServer, sessionID1, token2, promptID, "2")

	votes := getTestVote(t, testServer, sessionID1, token1, promptID)
	if len(votes) != 2 {
		t.Fatalf("Expected 2 vote, got %d", len(votes))
	}
}

func TestWaitForPrompt(t *testing.T) {
	promptText := "Test Prompt?"
	initialize()
	testServer := httptest.NewServer(CreateHandler())

	sessionID, token := createTestSession(t, testServer)

	var promptID string
	go func() {
		time.Sleep(1 * time.Second)
		promptID = createTestPrompt(t, testServer, sessionID, token, promptText)
	}()

	newPromptID := testWaitForPrompt(t, testServer, sessionID, token, "")

	if newPromptID != promptID {
		t.Fatalf("Expected prompt ID to be %s, got %s", promptID, newPromptID)
	}
}
