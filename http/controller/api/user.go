package api

import (
	"github.com/gin-gonic/gin"
	apiResp "github.com/RobertLesgros/rustdesk-interface/v2/http/response/api"
	"github.com/RobertLesgros/rustdesk-interface/v2/service"
	"net/http"
)

type User struct {
}

// currentUser Utilisateur actuel
// @Tags Utilisateur
// @Summary Informations utilisateur
// @Description Informations utilisateur
// @Accept  json
// @Produce  json
// @Success 200 {object} apiResp.UserPayload
// @Failure 500 {object} response.Response
// @Router /currentUser [get]
// @Security token
//func (u *User) currentUser(c *gin.Context) {
//	user := service.AllService.UserService.CurUser(c)
//	up := (&apiResp.UserPayload{}).FromName(user)
//	c.JSON(http.StatusOK, up)
//}

// Info Informations utilisateur
// @Tags Utilisateur
// @Summary Informations utilisateur
// @Description Informations utilisateur
// @Accept  json
// @Produce  json
// @Success 200 {object} apiResp.UserPayload
// @Failure 500 {object} response.Response
// @Router /currentUser [get]
// @Security token
func (u *User) Info(c *gin.Context) {
	user := service.AllService.UserService.CurUser(c)
	up := (&apiResp.UserPayload{}).FromUser(user)
	c.JSON(http.StatusOK, up)
}
