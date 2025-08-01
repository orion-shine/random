package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"random/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Rds    *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds := redis.MustNewRedis(c.Redis)
	return &ServiceContext{
		Config: c,
		Rds:    rds,
	}
}
