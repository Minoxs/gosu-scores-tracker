package config

import "osu-phantom/src/utils"

const (
	client_id_key     = "CLIENT_ID"
	client_secret_key = "CLIENT_SECRET"

	server_port_key = "SERVER_PORT"
)

var (
	CLIENT_ID     = utils.GetEnv(client_id_key).Integer(0)
	CLIENT_SECRET = utils.GetEnv(client_secret_key).String()

	SERVER_PORT = utils.GetEnv(server_port_key).Integer(4242)
)
