package utils

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
	"strings"
)

// GetSeedAndSecret
func GetSeedAndSecret(ctx context.Context, rds *redis.Redis, secret string) (string, string, error) {
	seed, err := rds.Get("seed")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			seed = "javaisthebestlanguageintheworld"
		} else {
			return "", "", err
		}
	}

	if secret == "" {
		secret = "helloworld"
	}

	return seed, secret, nil
}

// Int64SliceToString int64 转化为字符串
func Int64SliceToString(s []int64) string {
	strs := make([]string, len(s))
	for i, v := range s {
		strs[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(strs, ",")
}
