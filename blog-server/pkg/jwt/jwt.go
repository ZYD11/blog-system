package jwt

import (
	"blog-server/config"
	"blog-server/dao"
	"blog-server/models/response"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

// JWTAuth 中间件，检查token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//调用过滤去将放行的请求先放行
		if doSquare(c) {
			return
		}
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusOK, gin.H{
				"status": 401,
				"msg":    "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		s := strings.Split(token, " ")
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(s[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 401,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		appServer := config.GetServerCfg()
		lock := appServer.Lock
		if lock == "0" {
			get, err := dao.RedisDB.GET(claims.UserInfo.UserName)
			if err == nil {
				if !(get == s[1]) {
					c.JSON(http.StatusOK, gin.H{
						"status": 401,
						"msg":    "您的账号已在其他终端登录，请重新登录",
					})
					c.Abort()
					return
				}
			}
		}
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": 401,
				"msg":    err.Error(),
			})
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	TokenExpired     error  = errors.New("授权已过期")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("令牌非法")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "0df9b8db-6f7c-d713-eeab-ecb317696042"
)

// 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	UserInfo response.UserResponse `json:"userInfo"`
	jwt.StandardClaims
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateUserToken 生成含有用户信息的token
func (j *JWT) CreateUserToken(u *response.UserResponse) (string, error) {
	jwtConfig := config.GetJwtConfig()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserInfo: *u,
		StandardClaims: jwt.StandardClaims{
			//设置一小时时效
			ExpiresAt: time.Now().Add(jwtConfig.TimeOut * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    jwtConfig.Issuer,
		},
	})
	return claims.SignedString(j.SigningKey)
}

// 解析Tokne
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
