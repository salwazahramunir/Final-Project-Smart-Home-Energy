package model

import "mime/multipart"

type UploadFileRequest struct {
	File  *multipart.FileHeader `form:"file" binding:"required"`
	Query string                `form:"question" binding:"required"`
}
