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

type setting struct {
}

func (setting) GetFriends(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Setting.GetFriends(ctx, content.ID)
	reply.Reply(err, result)
}

func (setting) UpdateNickName(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateNickName)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	//fmt.Println("1---------", params.RelationID, params.NickName)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdateNickName(ctx, content.ID, params.RelationID, params.NickName)
	reply.Reply(err)
}

// UpdateSettingPin 更新好友或群组的pin（置顶）状态
// @Tags     setting
// @Summary  更新好友或群组的pin（置顶）状态
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                  true  "Bearer 账户令牌"
// @Param    data           body      request.ParamUpdateSettingPin  true  "关系ID，pin状态"
// @Success  200            {object}  common.State{}          "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 4002:关系不存在"
// @Router   /api/setting/update/pin [put]
func (setting) UpdateSettingPin(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateSettingPin)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdateSettingPin(ctx, content.ID, params.RelationID, *params.IsPin)
	reply.Reply(err)
}

// UpdateSettingDisturb 更新好友或群组的免打扰状态
// @Tags     setting
// @Summary  更新好友或群组的免打扰状态
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                  true  "Bearer 账户令牌"
// @Param    data           body      request.ParamUpdateSettingDisturb  true  "关系ID，免打扰状态"
// @Success  200            {object}  common.State{}          "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 4002:关系不存在"
// @Router   /api/setting/update/disturb [put]
func (setting) UpdateSettingDisturb(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateSettingDisturb)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdateSettingDisturb(ctx, content.ID, params.RelationID, *params.IsNotDisturb)
	reply.Reply(err)
}

// UpdateSettingShow 更新好友或群组的是否展示的状态
// @Tags     setting
// @Summary  更新好友或群组的是否展示的状态
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                  true  "Bearer 账户令牌"
// @Param    data           body      request.ParamUpdateSettingShow  true  "关系ID，展示状态"
// @Success  200            {object}  common.State{}          "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 4002:关系不存在"
// @Router   /api/setting/update/show [put]
func (setting) UpdateSettingShow(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamUpdateSettingShow)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.UpdateSettingShow(ctx, content.ID, params.RelationID, *params.IsShow)
	reply.Reply(err)
}

// GetPins 获取当前账户所有pin的好友和群组列表
// @Tags     setting
// @Summary  获取当前账户所有pin的好友和群组列表
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                  true  "Bearer 账户令牌"
// @Success  200            {object}  common.State{data=reply.ParamGetPins}          "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 4002:关系不存在"
// @Router   /api/setting/pins [get]
func (setting) GetPins(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Setting.GetPins(ctx, content.ID)
	reply.Reply(err, result)
}

// GetShows 获取当前账户首页显示的好友和群组列表
// @Tags     setting
// @Summary  获取当前账户首页显示的好友和群组列表
// @accept   application/json
// @Produce  application/json
// @Param    Authorization  header    string                  true  "Bearer 账户令牌"
// @Success  200            {object}  common.State{data=reply.ParamGetShows}          "1001:参数有误 1003:系统错误 2007:身份不存在 2008:身份验证失败 2010:账号不存在 4002:关系不存在"
// @Router   /api/setting/shows [get]
func (setting) GetShows(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	result, err := logic.Logics.Setting.GetShows(ctx, content.ID)
	reply.Reply(err, result)
}

func (setting) GetFriendsByName(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamGetFriendsByName)
	//打印是空的但是运行结果是对的
	//fmt.Println("参数：", params.Name)
	if err := ctx.ShouldBindQuery(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	limit, offset := global.Page.GetPageSizeAndOffset(ctx.Request)
	result, err := logic.Logics.Setting.GetFriendsByName(ctx, content.ID, params.Name, limit, offset)
	reply.Reply(err, result)
}

func (setting) DeleteFriend(ctx *gin.Context) {
	reply := app.NewResponse(ctx)
	params := new(request.ParamDeleteFriend)
	if err := ctx.ShouldBindJSON(params); err != nil {
		reply.Reply(errcode.ErrParamsNotValid.WithDetails(err.Error()))
		return
	}
	content, ok := middlewares.GetTokenContent(ctx)
	if !ok || content.TokenType != model.AccountToken {
		reply.Reply(errcodes.AuthNotExist)
		return
	}
	err := logic.Logics.Setting.DeleteFriend(ctx, content.ID, params.RelationID)
	reply.Reply(err)
}
