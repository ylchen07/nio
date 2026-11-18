package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPersistence(t *testing.T) {
	s := Server{path: filepath.Join(t.TempDir(), "nio.json")}
	s.db.Intros = []Intro{{ID: "abc", Name: "alice", Status: "pending"}}
	if err := s.save(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(s.path)
	if err != nil || !strings.Contains(string(b), "alice") {
		t.Fatal("failed to persist")
	}
}
func TestAuth(t *testing.T) {
	s := Server{adminKey: "correct-secret", db: Data{Agents: []Agent{{ID: "agent1", TokenHash: digest("agent-secret")}}}}
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer agent-secret")
	if a, ok := s.authenticate(r); !ok || a.ID != "agent1" {
		t.Fatal("agent auth failed")
	}
	r.Header.Set("Authorization", "Bearer wrong")
	if _, ok := s.authenticate(r); ok {
		t.Fatal("invalid token accepted")
	}
	r.Header.Set("Authorization", "Bearer correct-secret")
	if !s.admin(r) {
		t.Fatal("admin auth failed")
	}
}
func TestRateLimit(t *testing.T) {
	s := Server{hits: map[string][]time.Time{}}
	if s.limited("a", 1, time.Minute) {
		t.Fatal("first rejected")
	}
	if !s.limited("a", 1, time.Minute) {
		t.Fatal("second accepted")
	}
}
func TestIntro(t *testing.T) {
	s := Server{path: filepath.Join(t.TempDir(), "nio.json"), hits: map[string][]time.Time{}}
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/v1/introductions", strings.NewReader(`{"name":"AliceAgent","purpose":"Want to arrange meeting"}`))
	s.intro(w, r)
	if w.Code != 202 || len(s.db.Intros) != 1 {
		t.Fatalf("unexpected response: %d %s", w.Code, w.Body.String())
	}
}
func TestSkill(t *testing.T) {
	s := Server{}
	w := httptest.NewRecorder()
	s.skill(w, httptest.NewRequest(http.MethodGet, "/SKILL.md", nil))
	if !strings.Contains(w.Body.String(), `{"name":"YourAgent"`) {
		t.Fatal("invalid sample JSON")
	}
}
