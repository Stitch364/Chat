package operate

import (
	"chat/global"
	"github.com/XYYSWK/Lutils/pkg/utils"
	"github.com/gin-gonic/gin"
)

var UserKey = "user"

// SaveUserToken 将token存进redist
func (r *RDB) SaveUserToken(ctx *gin.Context, userID int64, tokens []string) error {
	//以 a : b 的形式拼接
	key := utils.LinkStr(UserKey, utils.IDToString(userID))
	for _, token := range tokens {
		if err := r.rdb.SAdd(ctx, key, token).Err(); err != nil {
			return err
		}
		r.rdb.Expire(ctx, key, global.PrivateSetting.Token.AccessTokenExpire)
	}
	return nil
}

// CheckUserTokenValid 判断token是否有效
func (r *RDB) CheckUserTokenValid(ctx *gin.Context, userID int64, token string) bool {
	key := utils.LinkStr(UserKey, utils.IDToString(userID))
	ok := r.rdb.SIsMember(ctx, key, token).Val()
	return ok
}

// DeleteAllTokenByUser 删除用户所有的token
func (r *RDB) DeleteAllTokenByUser(ctx *gin.Context, userID int64) error {
	key := utils.LinkStr(UserKey, utils.IDToString(userID))
	if err := r.rdb.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}
