package controller

import (
	"net/http"
	"net/url"
	"strings"

	"x-ui/logger"
	"x-ui/web/locale"
	"x-ui/web/session"

	"github.com/gin-gonic/gin"
)

type BaseController struct{}

func (a *BaseController) checkLogin(c *gin.Context) {
	// Browser-originated cross-site mutations are not allowed.
	// Requests without an Origin header (e.g. existing API clients) retain
	// compatibility; SameSite=Lax protects standard browser sessions.
	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead && c.Request.Method != http.MethodOptions {
		if origin := c.GetHeader("Origin"); origin != "" {
			parsed, err := url.Parse(origin)
			if err != nil || parsed.Host == "" || !strings.EqualFold(parsed.Host, c.Request.Host) || (c.Request.TLS != nil && parsed.Scheme != "https") {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		}
	}
	if !session.IsLogin(c) {
		if isAjax(c) {
			pureJsonMsg(c, http.StatusUnauthorized, false, I18nWeb(c, "pages.login.loginAgain"))
		} else {
			c.Redirect(http.StatusTemporaryRedirect, c.GetString("base_path"))
		}
		c.Abort()
	} else {
		c.Next()
	}
}

func I18nWeb(c *gin.Context, name string, params ...string) string {
	anyfunc, funcExists := c.Get("I18n")
	if !funcExists {
		logger.Warning("I18n function not exists in gin context!")
		return ""
	}
	i18nFunc, _ := anyfunc.(func(i18nType locale.I18nType, key string, keyParams ...string) string)
	msg := i18nFunc(locale.Web, name, params...)
	return msg
}
