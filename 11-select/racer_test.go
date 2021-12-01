package racer

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "time"
)

func TestRacer(t *testing.T) {

  t.Run("compare speed of server, return fatest one", func(t *testing.T) {
    slowServer := makeDelayServer(20 * time.Millisecond)
    fastServer := makeDelayServer(0 * time.Millisecond)

    defer slowServer.Close()
    defer fastServer.Close()

    slowURL := slowServer.URL
    fastURL := fastServer.URL

    want := fastURL
    got, err := Racer(slowURL, fastURL)

    if err != nil {
      t.Fatalf("did not expect an error but got one %v", err)
    }
    if got != want {
      t.Errorf("got %q but wnat %q", got, want)
    }
  })

  t.Run("return an error if a server doesn't respond within 10s", func(t *testing.T) {
    server := makeDelayServer(25 * time.Millisecond)

    defer server.Close()

    _, err := ConfigurableRacer(server.URL, server.URL, 10 * time.Millisecond)

    if err == nil {
      t.Error("expect an error but didn't get one")
    }
  })
}

func makeDelayServer(delay time.Duration) *httptest.Server {
  return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    time.Sleep(delay)
    w.WriteHeader(http.StatusOK)
  }))
}
