package my

import (
	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-api/v2/global"
	"github.com/RobertLesgros/rustdesk-api/v2/http/request/admin"
	"github.com/RobertLesgros/rustdesk-api/v2/http/response"
	"github.com/RobertLesgros/rustdesk-api/v2/model"
	"github.com/RobertLesgros/rustdesk-api/v2/service"
	"gorm.io/gorm"
)

type AddressBookCollectionRule struct {
}

// List Liste
// @Tags Mes règles de carnet d'adresses
// @Summary Liste des règles
// @Description Liste des règles
// @Accept  json
// @Produce  json
// @Param page query int false "Numéro de page"
// @Param page_size query int false "Taille de la page"
// @Param is_my query int false "Est-ce le mien"
// @Param user_id query int false "ID utilisateur"
// @Param collection_id query int false "ID de la collection"
// @Success 200 {object} response.Response{data=model.AddressBookCollectionList}
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection_rule/list [get]
// @Security token
func (abcr *AddressBookCollectionRule) List(c *gin.Context) {
	query := &admin.AddressBookCollectionRuleQuery{}
	if err := c.ShouldBindQuery(query); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	u := service.AllService.UserService.CurUser(c)
	query.UserId = int(u.Id)

	res := service.AllService.AddressBookService.ListRules(query.Page, query.PageSize, func(tx *gorm.DB) {
		tx.Where("user_id = ?", query.UserId)
		if query.CollectionId > 0 {
			tx.Where("collection_id = ?", query.CollectionId)
		}
	})
	response.Success(c, res)
}

// Create Créer une règle
// @Tags Mes règles de carnet d'adresses
// @Summary Créer une règle
// @Description Créer une règle
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollectionRule true "Informations sur la règle"
// @Success 200 {object} response.Response{data=model.AddressBookCollection}
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection_rule/create [post]
// @Security token
func (abcr *AddressBookCollectionRule) Create(c *gin.Context) {
	f := &model.AddressBookCollectionRule{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	if f.Type != model.ShareAddressBookRuleTypePersonal && f.Type != model.ShareAddressBookRuleTypeGroup {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	//t := f.ToAddressBookCollection()
	t := f
	u := service.AllService.UserService.CurUser(c)
	t.UserId = u.Id
	msg, res := abcr.CheckForm(u, t)
	if !res {
		response.Fail(c, 101, response.TranslateMsg(c, msg))
		return
	}
	err := service.AllService.AddressBookService.CreateRule(t)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

func (abcr *AddressBookCollectionRule) CheckForm(u *model.User, t *model.AddressBookCollectionRule) (string, bool) {
	if t.UserId != u.Id {
		return "NoAccess", false
	}
	if t.CollectionId > 0 && !service.AllService.AddressBookService.CheckCollectionOwner(t.UserId, t.CollectionId) {
		return "ParamsError", false
	}

	//check to_id
	if t.Type == model.ShareAddressBookRuleTypePersonal {
		if t.ToId == t.UserId {
			return "CannotShareToSelf", false
		}
		tou := service.AllService.UserService.InfoById(t.ToId)
		if tou.Id == 0 {
			return "ItemNotFound", false
		}
		// Les non-administrateurs ne peuvent pas partager avec des utilisateurs hors organisation
		//if tou.GroupId != u.GroupId {
		//	return "NoAccess", false
		//}
	} else if t.Type == model.ShareAddressBookRuleTypeGroup {
		// Les non-administrateurs ne peuvent pas partager avec d'autres groupes
		//if t.ToId != u.GroupId {
		//	return "NoAccess", false
		//}

		tog := service.AllService.GroupService.InfoById(t.ToId)
		if tog.Id == 0 {
			return "ItemNotFound", false
		}
	} else {
		return "ParamsError", false
	}
	// Vérification des doublons
	ex := service.AllService.AddressBookService.RuleInfoByToIdAndCid(t.Type, t.ToId, t.CollectionId)
	if t.Id == 0 && ex.Id > 0 {
		return "ItemExists", false
	}
	if t.Id > 0 && ex.Id > 0 && t.Id != ex.Id {
		return "ItemExists", false
	}
	return "", true
}

// Update Modifier
// @Tags Mes règles de carnet d'adresses
// @Summary Modifier la règle
// @Description Modifier la règle
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollectionRule true "Informations sur la règle"
// @Success 200 {object} response.Response{data=model.AddressBookCollection}
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection_rule/update [post]
// @Security token
func (abcr *AddressBookCollectionRule) Update(c *gin.Context) {
	f := &model.AddressBookCollectionRule{}
	if err := c.ShouldBindJSON(f); err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError")+err.Error())
		return
	}
	errList := global.Validator.ValidStruct(c, f)
	if len(errList) > 0 {
		response.Fail(c, 101, errList[0])
		return
	}
	if f.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}
	u := service.AllService.UserService.CurUser(c)

	ex := service.AllService.AddressBookService.RuleInfoById(f.Id)
	if ex.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	if ex.UserId != u.Id {
		response.Fail(c, 101, response.TranslateMsg(c, "NoAccess"))
		return
	}
	t := f
	msg, res := abcr.CheckForm(u, t)
	if !res {
		response.Fail(c, 101, response.TranslateMsg(c, msg))
		return
	}
	err := service.AllService.AddressBookService.UpdateRule(t)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
		return
	}
	response.Success(c, nil)
}

// Delete Supprimer
// @Tags Mes règles de carnet d'adresses
// @Summary Supprimer la règle
// @Description Supprimer la règle
// @Accept  json
// @Produce  json
// @Param body body model.AddressBookCollectionRule true "Informations sur la règle"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/my/address_book_collection_rule/delete [post]
// @Security token
func (abcr *AddressBookCollectionRule) Delete(c *gin.Context) {
	f := &model.AddressBookCollectionRule{}
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
	ex := service.AllService.AddressBookService.RuleInfoById(f.Id)
	if ex.Id == 0 {
		response.Fail(c, 101, response.TranslateMsg(c, "ItemNotFound"))
		return
	}
	u := service.AllService.UserService.CurUser(c)
	if ex.UserId != u.Id {
		response.Fail(c, 101, response.TranslateMsg(c, "NoAccess"))
		return
	}

	err := service.AllService.AddressBookService.DeleteRule(ex)
	if err == nil {
		response.Success(c, nil)
		return
	}
	response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed")+err.Error())
}
