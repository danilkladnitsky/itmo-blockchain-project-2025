package session

import (
	"fmt"
	"time"

	"memory-book/internal/lib/logger/sl"
	"memory-book/internal/lib/utils"
	"memory-book/internal/models"
	"memory-book/internal/storage/redis"
)

const (
	day   = time.Hour * 24
	week  = day * 7
	month = day * 30
)

func CreateAuthKey(userID int) string {
	return "UserID:" + string(rune(userID))
}
func CreateSessionIfNotExist(cache *redis.RedisDB, ui models.UserRedis) error {

	ctx, _ := utils.GetContext()

	var redisUser models.UserRedis
	err := cache.Get(ctx, CreateAuthKey(ui.UserID), &redisUser)
	if err == nil {
		err = cache.Set(ctx, CreateAuthKey(ui.UserID), redisUser, month)
		if err != nil {
			return fmt.Errorf("failed to update cache expiration", sl.Err(err))
		}
	} else {
		user := models.UserRedis{
			Username:       ui.Username,
			UserID:         ui.UserID,
			IsCompanyAdmin: ui.IsCompanyAdmin,
			CompanyID:      ui.CompanyID,
			IsServiceAdmin: false,
		}

		err = cache.Set(ctx, CreateAuthKey(user.UserID), user, month)
		if err != nil {
			return fmt.Errorf("failed to set cache", sl.Err(err))

		}
	}
	return nil
}
