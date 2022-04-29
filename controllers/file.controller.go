package controllers

import (
	"fmt"
	"github.com/andrwnv/event-aggregator/core/repo"
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
	userRepo        *repo.UserRepo
	httpContentType map[string]string
}

func handleSaveError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
		"error": "Unable to save the file(s)",
	})
}

func NewFileController(rootPath string, userRepo *repo.UserRepo) *FileController {
	return &FileController{
		DownloadPath: rootPath,
		userRepo:     userRepo,
		httpContentType: map[string]string{
			".jpg": "image/jpeg",
			".png": "image/png",
			".gif": "image/gif",
		},
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

	ctx.Status(http.StatusNotModified)
}

func (c *FileController) UploadAvatar(ctx *gin.Context) {
	payload, parseErr := misc.ExtractJwtPayload(ctx)
	if parseErr {
		misc.FailedClaimsExtractResponse(ctx)
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		misc.IncorrectRequestBodyResponse(ctx)
		return
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s%s", uuid.New(), ext)

	pathForSave := path.Join(c.DownloadPath, payload.ID)
	if os.MkdirAll(pathForSave, os.ModePerm) != nil {
		handleSaveError(ctx)
		return
	}

	if _, ok := c.httpContentType[ext]; ok {
		if err := ctx.SaveUploadedFile(file, path.Join(pathForSave, newFileName)); err != nil {
			handleSaveError(ctx)
			return
		}
	}

	user, err := c.userRepo.GetByEmail(payload.Email)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Cant extract user from database",
		})
		return
	}

	updateDto := repo.ToUpdateDto(user)
	updateDto.PhotoUrl = newFileName
	result, err := c.userRepo.Update(uuid.MustParse(payload.ID), updateDto)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Cant update user info!",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"result": result,
	})
}

func (c *FileController) UploadImagesMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, parseErr := misc.ExtractJwtPayload(ctx)
		if parseErr {
			misc.FailedClaimsExtractResponse(ctx)
			return
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			misc.IncorrectRequestBodyResponse(ctx)
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
