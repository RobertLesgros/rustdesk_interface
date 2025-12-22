package admin

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/RobertLesgros/rustdesk-interface/v2/global"
	"github.com/RobertLesgros/rustdesk-interface/v2/http/response"
	"github.com/RobertLesgros/rustdesk-interface/v2/lib/upload"
	"os"
	"time"
)

type File struct {
}

// OssToken retrieves OSS token for file upload
// @Tags File
// @Summary Get OSS token
// @Description Get OSS token for file upload
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/file/oss_token [get]
// @Security token
func (f *File) OssToken(c *gin.Context) {
	token := global.Oss.GetPolicyToken("")
	response.Success(c, token)
}

type FileBack struct {
	upload.CallbackBaseForm
	Url string `json:"url"`
}

// Notify is the callback after successful upload
func (f *File) Notify(c *gin.Context) {

	res := global.Oss.Verify(c.Request)
	if !res {
		response.Fail(c, 101, response.TranslateMsg(c, "NoAccess"))
		return
	}
	fm := &FileBack{}
	if err := c.ShouldBind(fm); err != nil {
		fmt.Println(err)
	}
	fm.Url = global.Config.Oss.Host + "/" + fm.Filename
	response.Success(c, fm)

}

// Upload uploads a file to local storage
// @Tags File
// @Summary Upload file to local storage
// @Description Upload file to local storage
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "File to upload"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/file/upload [post]
// @Security token
func (f *File) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "ParamsError"))
		return
	}

	// SECURITY: Sanitize filename to prevent path traversal attacks
	originalFilename := file.Filename
	safeFilename := filepath.Base(originalFilename)

	// Reject files with suspicious patterns
	if strings.Contains(safeFilename, "..") ||
		strings.HasPrefix(safeFilename, ".") ||
		safeFilename == "" {
		response.Fail(c, 101, response.TranslateMsg(c, "InvalidFilename"))
		return
	}

	// Generate unique filename to prevent overwrites
	uniqueFilename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), safeFilename)

	timePath := time.Now().Format("20060102") + "/"
	webPath := "/upload/" + timePath
	path := global.Config.Gin.ResourcesPath + webPath
	dst := path + uniqueFilename

	err = os.MkdirAll(path, 0750)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed"))
		return
	}

	// Upload file to specified directory
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		response.Fail(c, 101, response.TranslateMsg(c, "OperationFailed"))
		return
	}

	// Return file web address
	response.Success(c, gin.H{
		"url": webPath + uniqueFilename,
	})
}
