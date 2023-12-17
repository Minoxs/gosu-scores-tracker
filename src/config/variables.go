package config

import "osu-phantom/src/utils"

const (
	clientIdKey     = "CLIENT_ID"
	clientSecretKey = "CLIENT_SECRET"

	serverPortKey = "SERVER_PORT"
)

var (
	ClientId     = utils.GetEnv(clientIdKey).Integer(0)
	ClientSecret = utils.GetEnv(clientSecretKey).String()

	ServerPort = utils.GetEnv(serverPortKey).Integer(4242)
)
