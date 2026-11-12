package session

import (
    "crypto/sha256"
    "crypto/subtle"
    "encoding/gob"
    "net/http"

    "x-ui/database"
    "x-ui/database/model"

    "github.com/gin-contrib/sessions"
    "github.com/gin-gonic/gin"
)

const (
    loginUserKey = "LOGIN_USER"
    defaultPath = "/"
)

// LoginToken deliberately omits the bcrypt password hash and any user record.
// CredentialVersion invalidates existing sessions when the password changes.
type LoginToken struct {
    UserID int
    CredentialVersion [32]byte
}

func init() {
    gob.Register(LoginToken{})
    // Support decoding old signed cookies only to invalidate their sessions safely.
    gob.Register(model.User{})
}

func SetLoginUser(c *gin.Context, user *model.User) {
    if user == nil {
        return
    }
    sessions.Default(c).Set(loginUserKey, LoginToken{
        UserID: user.Id,
        CredentialVersion: sha256.Sum256([]byte(user.Password)),
    })
}

func SetMaxAge(c *gin.Context, maxAge int) {
    sessions.Default(c).Options(sessions.Options{
        Path: defaultPath,
        MaxAge: maxAge,
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        Secure: c.Request.TLS != nil,
    })
}

func GetLoginUser(c *gin.Context) *model.User {
    s := sessions.Default(c)
    obj := s.Get(loginUserKey)
    if obj == nil {
        return nil
    }
    token, ok := obj.(LoginToken)
    if !ok || token.UserID <= 0 {
        s.Delete(loginUserKey)
        return nil
    }
    user := &model.User{}
    if err := database.GetDB().First(user, token.UserID).Error; err != nil {
        s.Delete(loginUserKey)
        return nil
    }
    version := sha256.Sum256([]byte(user.Password))
    if subtle.ConstantTimeCompare(version[:], token.CredentialVersion[:]) != 1 {
        s.Delete(loginUserKey)
        return nil
    }
    return user
}

func IsLogin(c *gin.Context) bool {
    return GetLoginUser(c) != nil
}

func ClearSession(c *gin.Context) {
    s := sessions.Default(c)
    s.Clear()
    s.Options(sessions.Options{
        Path: defaultPath,
        MaxAge: -1,
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
        Secure: c.Request.TLS != nil,
    })
}
