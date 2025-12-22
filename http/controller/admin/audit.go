package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-interface/v2/global"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/request/admin"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/response"
	"github.com/RobertLesgros/rustdesk-interface/v2/model"
	"github.com/RobertLesgros/rustdesk-interface/v2/service"
	"gorm.io/gorm"
)

type Audit struct {
}

// ConnList Liste
// @Tags Journaux de connexion
// @Summary Liste des journaux de connexion
// @Description Liste des journaux de connexion
// @Accept  json
// @Produce  json
// @Param page query int false "Numéro de page"
// @Param page_size query int false "Taille de la page"
// @Param peer_id query int false "Appareil cible"
// @Param from_peer query int false "Appareil source"
// @Success 200 {object} response.Response{data=model.AuditConnList}
// @Failure 500 {object} response.Response
// @Router /admin/audit_conn/list [get]
// @Security token
func (a *Audit) ConnList(c *gin.Context) {
	query := &admin.AuditQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.AuditService.AuditConnList(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.PeerId != "" {
			tx.Where("peer_id like ?", "%"+query.PeerId+"%")
		}
		if query.FromPeer != "" {
			tx.Where("from_peer like ?", "%"+query.FromPeer+"%")
		}
		tx.Order("id desc")
	})
	response.Success(c, res)
}

// ConnDelete Supprimer
// @Tags Journaux de connexion
// @Summary Supprimer un journal de connexion
// @Description Supprimer un journal de connexion
// @Accept  json
// @Produce  json
// @Param body body model.AuditConn true "Informations sur le journal de connexion"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_conn/delete [post]
// @Security token
func (a *Audit) ConnDelete(c *gin.Context) {
	f := &model.AuditConn{}
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
	l := service.AllService.AuditService.ConnInfoById(f.Id)
	if l.Id > 0 {
		err := service.AllService.AuditService.DeleteAuditConn(l)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// BatchConnDelete Supprimer (Lot)
// @Tags Journaux de connexion
// @Summary Suppression par lot des journaux de connexion
// @Description Suppression par lot des journaux de connexion
// @Accept  json
// @Produce  json
// @Param body body admin.AuditConnLogIds true "Journaux de connexion"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_conn/batchDelete [post]
// @Security token
func (a *Audit) BatchConnDelete(c *gin.Context) {
	f := &admin.AuditConnLogIds{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}

	err := service.AllService.AuditService.BatchDeleteAuditConn(f.Ids)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, err.Error())
	return
}

// FileList Liste
// @Tags Journaux de fichiers
// @Summary Liste des journaux de fichiers
// @Description Liste des journaux de fichiers
// @Accept  json
// @Produce  json
// @Param page query int false "Numéro de page"
// @Param page_size query int false "Taille de la page"
// @Param peer_id query int false "Appareil cible"
// @Param from_peer query int false "Appareil source"
// @Success 200 {object} response.Response{data=model.AuditFileList}
// @Failure 500 {object} response.Response
// @Router /admin/audit_file/list [get]
// @Security token
func (a *Audit) FileList(c *gin.Context) {
	query := &admin.AuditQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.AuditService.AuditFileList(query.Page, query.PageSize, func(tx *gorm.DB) {
		if query.PeerId != "" {
			tx.Where("peer_id like ?", "%"+query.PeerId+"%")
		}
		if query.FromPeer != "" {
			tx.Where("from_peer like ?", "%"+query.FromPeer+"%")
		}
		tx.Order("id desc")
	})
	response.Success(c, res)
}

// FileDelete Supprimer
// @Tags Journaux de fichiers
// @Summary Supprimer un journal de fichier
// @Description Supprimer un journal de fichier
// @Accept  json
// @Produce  json
// @Param body body model.AuditFile true "Informations sur le journal de fichier"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_file/delete [post]
// @Security token
func (a *Audit) FileDelete(c *gin.Context) {
	f := &model.AuditFile{}
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
	l := service.AllService.AuditService.FileInfoById(f.Id)
	if l.Id > 0 {
		err := service.AllService.AuditService.DeleteAuditFile(l)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}

// BatchFileDelete Supprimer (Lot)
// @Tags Journaux de fichiers
// @Summary Suppression par lot des journaux de fichiers
// @Description Suppression par lot des journaux de fichiers
// @Accept  json
// @Produce  json
// @Param body body admin.AuditFileLogIds true "Journaux de fichiers"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/audit_file/batchDelete [post]
// @Security token
func (a *Audit) BatchFileDelete(c *gin.Context) {
	f := &admin.AuditFileLogIds{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if len(f.Ids) == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}

	err := service.AllService.AuditService.BatchDeleteAuditFile(f.Ids)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, err.Error())
	return
}
