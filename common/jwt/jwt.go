//@Package jwt
//@File jwt.go
//@Description jwt常用配置

package jwt

import (
	"time"

	"github.com/bertang/bert/common/config/application"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/jwt"
)

//UserClaim jwt中保存的数据
type UserClaim struct {
	ID     uint
	Name   string
	Mobile string
	Email  string
	Avatar string
}

//Sign jwt 签名
func Sign(claim *UserClaim) (token []byte, err error) {
	appConf := application.GetAppConf()
	age := 7 * time.Hour * 24
	if appConf.JwtMaxAge > 0 {
		age = time.Hour * time.Duration(appConf.JwtMaxAge)
	}
	token, err = jwt.Sign(jwt.HS256, []byte(application.GetAppConf().JwtSecret), claim, jwt.MaxAge(age))
	return
}

//JwtMiddleware 返回jwt中间件。用于验证jwt是否存在
func JwtMiddleware() context.Handler {
	verifier := jwt.NewVerifier(jwt.HS256, []byte(application.GetAppConf().JwtSecret))
	// Enable server-side token block feature (even before its expiration time):
	verifier.WithDefaultBlocklist()
	// Enable payload decryption with:
	// verifier.WithDecryption(encKey, nil)
	return verifier.Verify(func() interface{} {
		return new(UserClaim)
	})
}
