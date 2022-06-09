package controllers

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/core/repo"
	"github.com/andrwnv/event-aggregator/core/usecases"
	"github.com/andrwnv/event-aggregator/middleware"
	"github.com/andrwnv/event-aggregator/misc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type FileController struct {
	DownloadPath    string
	userUsecase     *usecases.UserUsecase
	httpContentType map[string]string
}

func handleSaveError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"error": "Unable to save the file(s)",
	})
}

func NewFileController(rootPath string, userUsecase *usecases.UserUsecase) *FileController {
	return &FileController{
		DownloadPath: rootPath,
		userUsecase:  userUsecase,
		httpContentType: map[string]string{
			".jpg": "image/jpeg",
			".png": "image/png",
			".gif": "image/gif",
		},
	}
}

func (c *FileController) MakeRoutesV1(rootGroup *gin.RouterGroup) {
	group := rootGroup.Group("/file")
	{
		group.PATCH("/update_avatar", middleware.AuthorizeJWTMiddleware(), c.UploadAvatar)
		group.GET("/img/:filename", c.GetImage)
	}
}

func (c *FileController) GetImage(ctx *gin.Context) {
	fileName := ctx.Param("filename")
	dirName := ctx.Query("uuid")
	targetPath := filepath.Join(c.DownloadPath, dirName, fileName)

	if !strings.HasPrefix(filepath.Clean(targetPath), c.DownloadPath) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Look like you attacking me",
		})
		return
	}

	fileExt := filepath.Ext(fileName)
	if fileExt == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cant find file extension",
		})
		return
	}

	ctx.Header("Content-Disposition", "inline")
	ctx.Header("Content-Type", c.httpContentType[fileExt])
	ctx.File(targetPath)

	ctx.Status(http.StatusOK)
}

func (c *FileController) UploadAvatar(ctx *gin.Context) {
	payload, err := misc.ExtractJwtPayload(ctx)
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	file, err := ctx.FormFile("file")
	if misc.HandleError(ctx, err, http.StatusBadRequest) {
		return
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New(), ext)

	pathForSave := path.Join(c.DownloadPath, payload.ID)
	if misc.HandleError(ctx, os.MkdirAll(pathForSave, os.ModePerm), http.StatusServiceUnavailable) {
		return
	}

	if _, ok := c.httpContentType[ext]; ok {
		if misc.HandleError(ctx, ctx.SaveUploadedFile(file, path.Join(pathForSave, newFileName)), http.StatusServiceUnavailable) {
			return
		}
	}

	user, err := c.userUsecase.GetFull(payload)
	if misc.HandleError(ctx, err, http.StatusInternalServerError, "Cant extract user from database.") {
		return
	}

	updateDto := repo.UserToUpdateUserDto(user)
	updateDto.PhotoUrl = newFileName
	result := c.userUsecase.Update(uuid.MustParse(payload.ID), updateDto)

	if misc.HandleError(ctx, result.Error, http.StatusInternalServerError, "Cant update user info.") {
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"result": result.Value,
	})
}

func (c *FileController) UploadImagesMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, err := misc.ExtractJwtPayload(ctx)
		if misc.HandleError(ctx, err, http.StatusBadRequest) {
			return
		}

		form, err := ctx.MultipartForm()
		if misc.HandleError(ctx, err, http.StatusBadRequest) {
			return
		}

		files := form.File["files"]
		if len(files) == 0 {
			ctx.Next()
			return
		}

		pathForSave := path.Join(c.DownloadPath, payload.ID)
		if os.MkdirAll(pathForSave, os.ModePerm) != nil {
			handleSaveError(ctx)
			return
		}

		var filesUrls []string
		for _, file := range files {
			ext := filepath.Ext(file.Filename)
			newFileName := fmt.Sprintf("%s%s", uuid.New(), ext)

			if _, ok := c.httpContentType[ext]; ok {
				if err := ctx.SaveUploadedFile(file, path.Join(pathForSave, newFileName)); err != nil {
					handleSaveError(ctx)
					return
				}

				filesUrls = append(filesUrls, newFileName)
			}
		}

		ctx.Set("file-names", filesUrls)
		ctx.Next()
	}
}
