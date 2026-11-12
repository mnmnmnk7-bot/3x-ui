package session

import (
    "net/http"
    "net/http/httptest"
    "path/filepath"
    "testing"

    "x-ui/database"
    "x-ui/database/model"
    "x-ui/util/crypto"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/gin-gonic/gin"
)

func TestSessionTokenDoesNotStorePasswordHashAndExpiresOnPasswordChange(t *testing.T) {
    if err := database.InitDB(filepath.Join(t.TempDir(), "test.db")); err != nil {
        t.Fatalf("init test db: %v", err)
    }
    defer database.CloseDB()

    var initialUser model.User
    if err := database.GetDB().First(&initialUser).Error; err != nil {
        t.Fatalf("get initial user: %v", err)
    }

    gin.SetMode(gin.TestMode)
    router := gin.New()
    store := cookie.NewStore([]byte("0123456789abcdef0123456789abcdef"))
    router.Use(sessions.Sessions("3x-ui", store))

    router.POST("/login", func(c *gin.Context) {
        SetMaxAge(c, 3600)
        SetLoginUser(c, &initialUser)
        if _, ok := sessions.Default(c).Get(loginUserKey).(LoginToken); !ok {
            t.Error("session contains more than the minimal login token")
        }
        if err := sessions.Default(c).Save(); err != nil {
            t.Errorf("save login: %v", err)
        }
        c.Status(http.StatusNoContent)
    })
    router.GET("/whoami", func(c *gin.Context) {
        user := GetLoginUser(c)
        if user == nil {
            c.Status(http.StatusUnauthorized)
            return
        }
        c.Status(http.StatusNoContent)
    })

    loginRecorder := httptest.NewRecorder()
    router.ServeHTTP(loginRecorder, httptest.NewRequest(http.MethodPost, "/login", nil))
    if loginRecorder.Code != http.StatusNoContent {
        t.Fatalf("login returned %d", loginRecorder.Code)
    }
    cookies := loginRecorder.Result().Cookies()
    if len(cookies) != 1 {
        t.Fatalf("expected a session cookie; got %d", len(cookies))
    }
    if cookies[0].HttpOnly != true || cookies[0].SameSite != http.SameSiteLaxMode {
        t.Fatal("session cookie missing HttpOnly / SameSite=Lax")
    }

    access := func() int {
        req := httptest.NewRequest(http.MethodGet, "/whoami", nil)
        req.AddCookie(cookies[0])
        recorder := httptest.NewRecorder()
        router.ServeHTTP(recorder, req)
        return recorder.Code
    }
    if got := access(); got != http.StatusNoContent {
        t.Fatalf("valid session returned %d", got)
    }

    updated, err := crypto.HashPasswordAsBcrypt("different-long-secret")
    if err != nil {
        t.Fatal(err)
    }
    if err := database.GetDB().Model(&model.User{}).Where("id = ?", initialUser.Id).Update("password", updated).Error; err != nil {
        t.Fatal(err)
    }
    if got := access(); got != http.StatusUnauthorized {
        t.Fatalf("old session after password change returned %d; want 401", got)
    }
}
