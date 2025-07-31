package api

import (
	"chat/errcodes"
	"chat/global"
	"chat/logic"
	"chat/middlewares"
	"chat/model"
	"chat/model/request"
	"github.com/XYYSWK/Lutils/pkg/app"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type file struct {
}

func (file) PublishFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamPublishFile)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid)
		return
	}
	//fileType, myErr := gtype.GetFileType(params.File)
	//if myErr != nil {
	//	global.Logger.Error(myErr.Error(), middlewares.ErrLogMsg(ctx)...)
	//	reply.Reply(errcode.ErrServer)
	//	return
	//}
	//if fileType != "img" && fileType != "png" && fileType != "jpg" {
	//	fileType = "file"
	//}
	//fmt.Println("Publishfiletype", fileType)
	if params.File == nil {
		reply.Reply(errcodes.FileIsEmpty)
		return
	}
	if params.File.Size > global.PublicSetting.Rules.BiggestFileSize {
		reply.Reply(errcodes.FileTooBig)
		return
	}

	result, err := logic.Logics.File.PublishFile(ctx, model.PublishFile{
		File:       params.File,
		RelationID: params.RelationID,
		AccountID:  params.AccountID,
	})
	reply.Reply(err, result)
}

func (file) DeleteFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteFile)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	err := logic.Logics.File.DeleteFile(ctx, params.FileID)
	reply.Reply(err)
}

func (file) GetRelationFile(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetRelationFile)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.File.GetRelationFile(ctx, params.RelationID)
	reply.Reply(err, result)
}

func (file) UploadAccountAvatar(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUploadAccountAvatar)
	if err := ctx.ShouldBind(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.File.UploadAccountAvatar(ctx, content.ID, params.File)
	reply.Reply(err, result)
}

func (file) GetFileDetailsByID(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetFileDetailsByID)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	result, err := logic.Logics.File.GetFileDetailsByID(ctx, params.FileID)
	reply.Reply(err, result)
}
