package my

import (
	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/request/admin"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/response"
	"github.com/RobertLesgros/rustdesk-interface/v2/service"
	"gorm.io/gorm"
	"time"
)

type Peer struct {
}

// List Liste
// @Tags Mes appareils
// @Summary Liste des appareils
// @Description Liste des appareils
// @Accept  json
// @Produce  json
// @Param page query int false "Numéro de page"
// @Param page_size query int false "Taille de la page"
// @Param time_ago query int false "Temps"
// @Param id query string false "ID"
// @Param hostname query string false "Nom d'hôte"
// @Param uuids query string false "uuids séparés par des virgules"
// @Success 200 {object} response.Response{data=model.PeerList}
// @Failure 500 {object} response.Response
// @Router /admin/my/peer/list [get]
// @Security token
func (ct *Peer) List(c *gin.Context) {
	query := &admin.PeerQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	u := service.AllService.UserService.CurUser(c)
	res := service.AllService.PeerService.List(query.Page, query.PageSize, func(tx *gorm.DB) {
		tx.Where("user_id = ?", u.Id)
		if query.TimeAgo > 0 {
			lt := time.Now().Unix() - int64(query.TimeAgo)
			tx.Where("last_online_time < ?", lt)
		}
		if query.TimeAgo < 0 {
			lt := time.Now().Unix() + int64(query.TimeAgo)
			tx.Where("last_online_time > ?", lt)
		}
		if query.Id != "" {
			tx.Where("id like ?", "%"+query.Id+"%")
		}
		if query.Hostname != "" {
			tx.Where("hostname like ?", "%"+query.Hostname+"%")
		}
		if query.Uuids != "" {
			tx.Where("uuid in (?)", query.Uuids)
		}
	})
	response.Success(c, res)
}
