package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-interface/v2/global"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/request/admin"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/response"
	"github.com/RobertLesgros/rustdesk-interface/v2/service"
	"strconv"
)

type DeviceGroup struct {
}

// Detail Groupe de périphériques
// @Tags Groupe de périphériques
// @Summary Détails du groupe de périphériques
// @Description Détails du groupe de périphériques
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} response.Response{data=model.Group}
// @Failure 500 {object} response.Response
// @Router /admin/device_group/detail/{id} [get]
// @Security token
func (ct *DeviceGroup) Detail(c *gin.Context) {
	id := c.Param("id")
	iid, _ := strconv.Atoi(id)
	u := service.AllService.GroupService.DeviceGroupInfoById(uint(iid))
	if u.Id > 0 {
		response.Success(c, u)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
	return
}

// Create Créer un groupe de périphériques
// @Tags Groupe de périphériques
// @Summary Créer un groupe de périphériques
// @Description Créer un groupe de périphériques
// @Accept  json
// @Produce  json
// @Param body body admin.DeviceGroupForm true "Informations sur le groupe de périphériques"
// @Success 200 {object} response.Response{data=model.DeviceGroup}
// @Failure 500 {object} response.Response
// @Router /admin/device_group/create [post]
// @Security token
func (ct *DeviceGroup) Create(c *gin.Context) {
	f := &admin.DeviceGroupForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := f.ToDeviceGroup()
	err := service.AllService.GroupService.DeviceGroupCreate(u)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// List Liste
// @Tags Groupe
// @Summary Liste des groupes
// @Description Liste des groupes
// @Accept  json
// @Produce  json
// @Param page query int false "Numéro de page"
// @Param page_size query int false "Taille de la page"
// @Success 200 {object} response.Response{data=model.GroupList}
// @Failure 500 {object} response.Response
// @Router /admin/device_group/list [get]
// @Security token
func (ct *DeviceGroup) List(c *gin.Context) {
	query := &admin.PageQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	res := service.AllService.GroupService.DeviceGroupList(query.Page, query.PageSize, nil)
	response.Success(c, res)
}

// Update Modifier
// @Tags Groupe de périphériques
// @Summary Modifier le groupe de périphériques
// @Description Modifier le groupe de périphériques
// @Accept  json
// @Produce  json
// @Param body body admin.DeviceGroupForm true "Informations sur le groupe"
// @Success 200 {object} response.Response{data=model.Group}
// @Failure 500 {object} response.Response
// @Router /admin/device_group/update [post]
// @Security token
func (ct *DeviceGroup) Update(c *gin.Context) {
	f := &admin.DeviceGroupForm{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	if f.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	u := f.ToDeviceGroup()
	err := service.AllService.GroupService.DeviceGroupUpdate(u)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete Supprimer
// @Tags Groupe de périphériques
// @Summary Supprimer le groupe de périphériques
// @Description Supprimer le groupe de périphériques
// @Accept  json
// @Produce  json
// @Param body body admin.DeviceGroupForm true "Informations sur le groupe"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/device_group/delete [post]
// @Security token
func (ct *DeviceGroup) Delete(c *gin.Context) {
	f := &admin.DeviceGroupForm{}
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
	u := service.AllService.GroupService.DeviceGroupInfoById(f.Id)
	if u.Id > 0 {
		err := service.AllService.GroupService.DeviceGroupDelete(u)
		if err == nil {
			response.Success(c, nil)
			return
		}
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
}
