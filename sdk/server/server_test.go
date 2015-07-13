package server_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikebeyer/clc-sdk/clc"
	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/server"
	"github.com/stretchr/testify/assert"
)

var name = "va1testserver01"

func TestGetServer(t *testing.T) {
	assert := assert.New(t)

	ms := mockServerAPI()
	defer ms.Close()

	service := service(ms.URL)
	resp, err := service.Get(name)

	assert.Nil(err)
	assert.Equal(name, resp.Name)
}

func TestGetServerByUUID(t *testing.T) {
	assert := assert.New(t)

	ms := mockServerAPI()
	defer ms.Close()

	service := service(ms.URL)
	resp, err := service.Get("5404cf5ece2042dc9f2ac16ab67416bb")

	assert.Nil(err)
	assert.Equal("va1testserver01", resp.Name)
}

func TestCreateServer(t *testing.T) {
	assert := assert.New(t)

	ms := mockServerAPI()
	defer ms.Close()

	service := service(ms.URL)
	server := server.Server{
		Name:           "va1testserver01",
		CPU:            1,
		MemoryGB:       1,
		GroupID:        "group",
		SourceServerID: "UBUNTU",
		Type:           "standard",
	}
	s, err := service.Create(server, nil)

	assert.Nil(err)
	assert.True(s.IsQueued)
	assert.Equal(server.Name, s.Server)
}

func mockServerAPI() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/servers/test", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		server := &server.Server{}
		err := json.NewDecoder(r.Body).Decode(server)
		if err != nil {
			http.Error(w, "server err", http.StatusInternalServerError)
			return
		}

		if !server.Valid() {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `{"server":"va1testserver01","isQueued":true,"links":[{"rel":"status","href":"/v2/operations/test/status/12345","id":"12345"},{"rel":"self","href":"/v2/servers/test/12345?uuid=True","id":"12345","verbs":["GET"]}]}`)
	})
	mux.HandleFunc("/servers/test/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if len(r.URL.Query()) == 0 {
			server := &clc.ServerResponse{Name: name}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}

		if r.URL.Query().Get("uuid") == "true" {
			server := &clc.ServerResponse{Name: "va1testserver01"}
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(server)
			return
		}

		http.Error(w, "not found", http.StatusNotFound)
	})

	return httptest.NewServer(mux)
}
func service(url string) *server.Service {
	config := api.Config{
		User: api.User{
			Username: "test.user",
			Password: "s0s3cur3",
		},
		Alias:   "test",
		BaseURL: url,
	}

	client := api.New(config)
	client.Token = api.Token{Token: "validtoken"}
	return server.New(client)
}
