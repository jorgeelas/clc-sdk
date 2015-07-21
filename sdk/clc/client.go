package clc

import (
	"github.com/mikebeyer/clc-sdk/sdk/aa"
	"github.com/mikebeyer/clc-sdk/sdk/alert"
	"github.com/mikebeyer/clc-sdk/sdk/api"
	"github.com/mikebeyer/clc-sdk/sdk/lb"
	"github.com/mikebeyer/clc-sdk/sdk/server"
	"github.com/mikebeyer/clc-sdk/sdk/status"
)

type Client struct {
	client *api.Client

	Server *server.Service
	Status *status.Service
	AA     *aa.Service
	Alert  *alert.Service
	LB     *lb.Service
}

func New(config api.Config) *Client {
	c := &Client{
		client: api.New(config),
	}

	c.Server = server.New(c.client)
	c.Status = status.New(c.client)
	c.AA = aa.New(c.client)
	c.Alert = alert.New(c.client)
	c.LB = lb.New(c.client)

	return c
}
