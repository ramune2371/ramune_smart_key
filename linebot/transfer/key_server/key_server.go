package key_server

import "linebot/entity"

type KeyServerTransfer interface {
	OpenKey() (entity.KeyServerResponse, error)
	CloseKey() (entity.KeyServerResponse, error)
	CheckKey() (entity.KeyServerResponse, error)
}
