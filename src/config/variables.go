package config

import "osu-phantom/src/utils"

const (
	client_id_key     = "CLIENT_ID"
	client_secret_key = "CLIENT_SECRET"
)

var (
	CLIENT_ID     = utils.GetEnv(client_id_key).Integer(0)
	CLIENT_SECRET = utils.GetEnv(client_secret_key).String()
)
