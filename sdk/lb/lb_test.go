package lb_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/lb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLB(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Get", "http://localhost/v2/sharedLoadBalancers/test/dc1/12345", mock.Anything).Return(nil)
	service := lb.New(client)

	id := "12345"
	resp, err := service.Get("dc1", id)

	assert.Nil(err)
	assert.Equal(id, resp.ID)
}

func TestGetAllLBs(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Get", "http://localhost/v2/sharedLoadBalancers/test/dc1", mock.Anything).Return(nil)
	service := lb.New(client)

	resp, err := service.GetAll("dc1")

	assert.Nil(err)
	assert.Equal(1, len(resp))
}

func TestCreateLB(t *testing.T) {
	assert := assert.New(t)

	client := NewMockClient()
	client.On("Post", "http://localhost/v2/sharedLoadBalancers/test/dc1", mock.Anything, mock.Anything).Return(nil)
	service := lb.New(client)

	lb := lb.LoadBalancer{
		Name:        "new",
		Description: "balancing load",
	}
	resp, err := service.Create("dc1", lb)

	assert.Nil(err)
	assert.Equal(lb.Name, resp.Name)
	assert.Equal("enabled", resp.Status)
	assert.NotEmpty(resp.ID)
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Get(url string, resp interface{}) error {
	if strings.HasSuffix(url, "12345") {
		json.Unmarshal([]byte(`{"id":"12345","name":"new","description":"balancing load","ipAddress":"10.10.10.10","status":"enabled","pools":[],"links":[{"rel":"self","href":"/v2/sharedLoadBalancers/test/dc1/12345","verbs":["GET","PUT","DELETE"]},{"rel":"pools","href":"/v2/sharedLoadBalancers/test/dc1/12345/pools","verbs":["GET","POST"]}]}`), resp)
	} else {
		json.Unmarshal([]byte(`[{"id":"12345","name":"new","description":"balancing load","ipAddress":"10.10.10.10","status":"enabled","pools":[],"links":[{"rel":"self","href":"/v2/sharedLoadBalancers/test/dc1/12345","verbs":["GET","PUT","DELETE"]},{"rel":"pools","href":"/v2/sharedLoadBalancers/test/dc1/12345/pools","verbs":["GET","POST"]}]}]`), resp)
	}

	args := m.Called(url, resp)
	return args.Error(0)
}

func (m *MockClient) Post(url string, body, resp interface{}) error {
	json.Unmarshal([]byte(`{"id":"12345","name":"new","description":"balancing load","ipAddress":"10.10.10.10","status":"enabled","pools":[],"links":[{"rel":"self","href":"/v2/sharedLoadBalancers/test/dc1/12345","verbs":["GET","PUT","DELETE"]},{"rel":"pools","href":"/v2/sharedLoadBalancers/test/dc1/12345/pools","verbs":["GET","POST"]}]}`), resp)
	args := m.Called(url, body, resp)
	return args.Error(0)
}

func (m *MockClient) Put(url string, body, resp interface{}) error {
	args := m.Called(url, body, resp)
	return args.Error(0)
}

func (m *MockClient) Patch(url string, body, resp interface{}) error {
	args := m.Called(url, body, resp)
	return args.Error(0)
}

func (m *MockClient) Delete(url string, resp interface{}) error {
	args := m.Called(url, resp)
	return args.Error(0)
}

func (m *MockClient) Config() *api.Config {
	return &api.Config{
		User: api.User{
			Username: "test.user",
			Password: "s0s3cur3",
		},
		Alias:   "test",
		BaseURL: "http://localhost/v2",
	}
}
