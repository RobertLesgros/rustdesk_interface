package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-api/v2/global"
	"github.com/RobertLesgros/rustdesk-api/v2/http/request/admin"
	"github.com/RobertLesgros/rustdesk-api/v2/http/response"
	"github.com/RobertLesgros/rustdesk-api/v2/service"
	"gorm.io/gorm"
)

type ShareRecord struct {
}

// List Liste
// @Tags Enregistrements de partage
// @Summary Liste des enregistrements de partage
// @Description Liste des enregistrements de partage
// @Accept  json
// @Produce  json
// @Param user_id query int false "ID utilisateur"
// @Param page query int false "Numéro de page"
// @Param page_size query int false "Taille de la page"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/share_record/list [get]
// @Security token
func (sr *ShareRecord) List(c *gin.Context) {
	query := &admin.ShareRecordQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.ShareRecordService.List(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.UserId > 0 {
			tx.Where("user_id = ?", query.UserId)
		}
	})
	response.Success(c, res)
}

// Delete Supprimer
// @Tags Enregistrements de partage
// @Summary Supprimer un enregistrement de partage
// @Description Supprimer un enregistrement de partage
// @Accept  json
// @Produce  json
// @Param body body admin.ShareRecordForm true "Informations sur l'enregistrement de partage"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/share_record/delete [post]
// @Security token
func (sr *ShareRecord) Delete(c *gin.Context) {
	f := &admin.ShareRecordForm{}
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
	i := service.AllService.ShareRecordService.InfoById(f.Id)
	if i.Id > 0 {
		err := service.AllService.ShareRecordService.Delete(i)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// BatchDelete Suppression par lot
// @Tags Enregistrements de partage
// @Summary Suppression par lot des enregistrements de partage
// @Description Suppression par lot des enregistrements de partage
// @Accept  json
// @Produce  json
// @Param body body admin.PeerShareRecordBatchDeleteForm true "id"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/share_record/batchDelete [post]
// @Security token
func (sr *ShareRecord) BatchDelete(c *gin.Context) {
	f := &admin.PeerShareRecordBatchDeleteForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	err := service.AllService.ShareRecordService.BatchDelete(f.Ids)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}
