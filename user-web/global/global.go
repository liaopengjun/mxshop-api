package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
)

var (
	Trans         ut.Translator
	ServerConfig  = &config.ServerConfig{}
	UserSrvClient proto.UserClient
	Log           *zap.Logger
	Redis         *redis.Client
)
