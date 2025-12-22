package my

import (
	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-interface/v2/global"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/request/admin"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/response"
	"github.com/RobertLesgros/rustdesk-interface/v2/model"
	"github.com/RobertLesgros/rustdesk-interface/v2/service"
	"gorm.io/gorm"
)

type LoginLog struct {
}

// List Liste
// @Tags Mon journal de connexion
// @Summary Liste des journaux de connexion
// @Description Liste des journaux de connexion
// @Accept  json
// @Produce  json
// @Param page query int false "Numéro de page"
// @Param page_size query int false "Taille de la page"
// @Param user_id query int false "ID utilisateur"
// @Success 200 {object} response.Response{data=model.LoginLogList}
// @Failure 500 {object} response.Response
// @Router /admin/my/login_log/list [get]
// @Security token
func (ct *LoginLog) List(c *gin.Context) {
	query := &admin.LoginLogQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	u := service.AllService.UserService.CurUser(c)
	res := service.AllService.LoginLogService.List(query.Page, query.PageSize, func(tx *gorm.DB) {
		tx.Where("user_id = ? and is_deleted = ?", u.Id, model.IsDeletedNo)
		tx.Order("id desc")
	})
	response.Success(c, res)
}

// Delete Supprimer
// @Tags Mon journal de connexion
// @Summary Supprimer un journal de connexion
// @Description Supprimer un journal de connexion
// @Accept  json
// @Produce  json
// @Param body body model.LoginLog true "Informations sur le journal de connexion"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/my/login_log/delete [post]
// @Security token
func (ct *LoginLog) Delete(c *gin.Context) {
	f := &model.LoginLog{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	id := f.Id
	errList := global.Validator.ValidVar(c, id, "required,gt=0")
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	l := service.AllService.LoginLogService.InfoById(f.Id)
	if l.Id == 0 || l.IsDeleted == model.IsDeletedYes {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	u := service.AllService.UserService.CurUser(c)
	if l.UserId != u.Id {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	err := service.AllService.LoginLogService.SoftDelete(l)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, err.Error())
}

// BatchDelete Supprimer (Lot)
// @Tags Mon journal de connexion
// @Summary Suppression par lot des journaux de connexion
// @Description Suppression par lot des journaux de connexion
// @Accept  json
// @Produce  json
// @Param body body admin.LoginLogIds true "Journaux de connexion"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/my/login_log/batchDelete [post]
// @Security token
func (ct *LoginLog) BatchDelete(c *gin.Context) {
	f := &admin.LoginLogIds{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	u := service.AllService.UserService.CurUser(c)
	err := service.AllService.LoginLogService.BatchSoftDelete(u.Id, f.Ids)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, err.Error())
	return
}
