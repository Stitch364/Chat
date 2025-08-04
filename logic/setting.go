package logic

import (
	"chat/dao"
	db "chat/dao/mysql/sqlc"
	"chat/errcodes"
	"chat/global"
	"chat/middlewares"
	"chat/model"
	"chat/model/chat/server"
	"chat/model/reply"
	"chat/task"
	"context"
	"database/sql"
	"errors"
	"github.com/XYYSWK/Lutils/pkg/app/errcode"
	"github.com/gin-gonic/gin"
)

type setting struct{}

// ExistsSetting 是否存在 account 和 relation 关系的联系
// 参数：accountID，relationDI
// 成功：是否存在，nil
// 失败：打印错误日志 errcode.ErrServer
func ExistsSetting(ctx context.Context, accountID, relationID int64) (bool, errcode.Err) {
	ok, err := dao.Database.DB.ExistsSetting(ctx, &db.ExistsSettingParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return false, errcode.ErrServer
	}
	return ok, nil
}

func (setting) GetFriends(ctx *gin.Context, accountID int64) (*reply.ParamGetFriends, errcode.Err) {
	data, err := dao.Database.DB.GetFriendSettingsOrderByName(ctx, &db.GetFriendSettingsOrderByNameParams{
		AccountID: accountID,
		ID:        accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.SettingFriend, 0, len(data))
	for _, v := range data {
		result = append(result, &model.SettingFriend{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: "friend",
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			FriendInfo: &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			},
		})
	}
	return &reply.ParamGetFriends{
		List:  result,
		Total: int64(len(result)),
	}, nil
}
func (setting) UpdateNickName(ctx *gin.Context, accountID, relationID int64, nickName string) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.NickName == nickName {
			return nil
		}
		if err := dao.Database.DB.UpdateSettingNickName(ctx, &db.UpdateSettingNickNameParams{
			NickName:   nickName,
			AccountID:  accountID,
			RelationID: relationID,
		}); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改昵称的通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateNickName(accessToken, accountID, relationID, nickName))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (setting) UpdateSettingPin(ctx *gin.Context, accountId, relationID int64, isPin bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountId,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsPin == isPin {
			return nil
		}
		if err := dao.Database.DB.UpdateSettingPin(ctx, &db.UpdateSettingPinParams{
			IsPin:      isPin,
			AccountID:  accountId,
			RelationID: relationID,
		}); err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改置顶的通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingPin, accountId, relationID, isPin))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (setting) UpdateSettingDisturb(ctx *gin.Context, accountID, relationID int64, isNotDisturb bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsNotDisturb == isNotDisturb {
			return nil
		}
		err = dao.Database.DB.UpdateSettingDisturb(ctx, &db.UpdateSettingDisturbParams{
			IsNotDisturb: isNotDisturb,
			AccountID:    accountID,
			RelationID:   relationID,
		})
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改免打扰的通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingNotDisturb, accountID, relationID, isNotDisturb))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (setting) UpdateSettingShow(ctx *gin.Context, accountID, relationID int64, isShow bool) errcode.Err {
	settingInfo, err := dao.Database.DB.GetSettingByID(ctx, &db.GetSettingByIDParams{
		AccountID:  accountID,
		RelationID: relationID,
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return errcodes.RelationNotExists
	case errors.Is(err, nil):
		if settingInfo.IsShow == isShow {
			return nil
		}
		err = dao.Database.DB.UpdateSettingShow(ctx, &db.UpdateSettingShowParams{
			IsShow:     isShow,
			AccountID:  accountID,
			RelationID: relationID,
		})
		if err != nil {
			global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
			return errcode.ErrServer
		}
		// 向自己推送更改展示状态的通知
		accessToken, _ := middlewares.GetToken(ctx.Request.Header)
		global.Worker.SendTask(task.UpdateSettingState(accessToken, server.SettingShow, accountID, relationID, isShow))
		return nil
	default:
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
}

func (setting) GetPins(ctx *gin.Context, accountID int64) (*reply.ParamGetPins, errcode.Err) {
	friendData, err := dao.Database.DB.GetFriendPinSettingsOrderByPinTime(ctx, &db.GetFriendPinSettingsOrderByPinTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	groupData, err := dao.Database.DB.GetGroupPinSettingsOrderByPinTime(ctx, &db.GetGroupPinSettingsOrderByPinTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.SettingPin, 0, len(friendData)+len(groupData))
	for i, j := 0, 0; i < len(friendData) || j < len(groupData); {
		if i < len(friendData) && (j >= len(groupData) || friendData[i].PinTime.Before(groupData[j].PinTime)) {
			v := friendData[i]
			friendInfo := &model.SettingFriendInfo{
				AccountID: accountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			}
			result = append(result, &model.SettingPin{
				SettingPinInfo: model.SettingPinInfo{
					RelationID:   v.RelationID,
					RelationType: "friend",
					NickName:     v.NickName,
					PinTime:      v.PinTime,
				},
				FriendInfo: friendInfo,
			})
			i++
		} else {
			v := groupData[j]
			groupInfo := &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        v.Name.String, // 去掉前面的 "
				Description: v.Description.String,
				Avatar:      v.Avatar.String, // 去掉后面的 "
			}
			result = append(result, &model.SettingPin{
				SettingPinInfo: model.SettingPinInfo{
					RelationID:   v.RelationID,
					RelationType: "group",
					NickName:     v.NickName,
					PinTime:      v.PinTime,
				},
				GroupInfo: groupInfo,
			})
			j++
		}
	}
	return &reply.ParamGetPins{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func (setting) GetShows(ctx *gin.Context, accountID int64) (*reply.ParamGetShows, errcode.Err) {
	friendData, err := dao.Database.DB.GetFriendShowSettingsOrderByShowTime(ctx, &db.GetFriendShowSettingsOrderByShowTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	groupData, err := dao.Database.DB.GetGroupShowSettingsOrderByShowTime(ctx, &db.GetGroupShowSettingsOrderByShowTimeParams{
		AccountID:   accountID,
		AccountID_2: accountID,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	result := make([]*model.Setting, 0, len(friendData)+len(groupData))
	for i, j := 0, 0; i < len(friendData) || j < len(groupData); {
		if i < len(friendData) && (j >= len(groupData) || friendData[i].LastShow.After(groupData[j].LastShow)) {
			v := friendData[i]
			friendInfo := &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			}
			result = append(result, &model.Setting{
				SettingInfo: model.SettingInfo{
					RelationID:   v.RelationID,
					RelationType: "friend",
					NickName:     v.NickName,
					IsNotDisturb: v.IsNotDisturb,
					IsPin:        v.IsPin,
					IsShow:       v.IsShow,
					PinTime:      v.PinTime,
					LastShow:     v.LastShow,
				},
				FriendInfo: friendInfo,
			})
			i++
		} else {
			v := groupData[j]
			groupInfo := &model.SettingGroupInfo{
				RelationID:  v.RelationID,
				Name:        v.Name.String,
				Description: v.Description.String,
				Avatar:      v.Avatar.String,
			}
			result = append(result, &model.Setting{
				SettingInfo: model.SettingInfo{
					RelationID:   v.RelationID,
					RelationType: "group",
					NickName:     v.NickName,
					IsNotDisturb: v.IsNotDisturb,
					IsPin:        v.IsPin,
					IsShow:       v.IsShow,
					PinTime:      v.PinTime,
					LastShow:     v.LastShow,
					IsLeader:     v.IsLeader,
				},
				GroupInfo: groupInfo,
			})
			j++
		}
	}
	return &reply.ParamGetShows{
		List:  result,
		Total: int64(len(result)),
	}, nil
}

func getFriendRelationByID(ctx *gin.Context, relationID int64) (*db.GetFriendRelationByIDRow, errcode.Err) {
	relationInfo, err := dao.Database.DB.GetFriendRelationByID(ctx, relationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.RelationNotExists
		}
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	return relationInfo, nil
}

func (setting) DeleteFriend(ctx *gin.Context, accountID, relationID int64) errcode.Err {
	friendInfo, myErr := getFriendRelationByID(ctx, relationID)
	if myErr != nil {
		return myErr
	}
	if friendInfo.Account1ID.Int64 != accountID && friendInfo.Account2ID.Int64 != accountID {
		return errcodes.AuthPermissionsInsufficient
	}
	if friendInfo.Account1ID.Int64 == accountID && friendInfo.Account2ID.Int64 == accountID {
		return errcodes.DeleteNotLegal
	}
	if err := dao.Database.DB.DeleteRelationWithTx(ctx, dao.Database.Redis, relationID); err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return errcode.ErrServer
	}
	// 推送删除通知
	accessToken, _ := middlewares.GetToken(ctx.Request.Header)
	global.Worker.SendTask(task.DeleteRelation(accessToken, accountID, relationID))
	return nil
}

func (setting) GetFriendsByName(ctx *gin.Context, accountID int64, name string, limit, offset int32) (*reply.ParamGetFriendsByName, errcode.Err) {
	data, err := dao.Database.DB.GetFriendSettingsByName(ctx, &db.GetFriendSettingsByNameParams{
		AccountID:   accountID,
		AccountID_2: accountID,
		Limit:       limit,
		Offset:      offset,
		CONCAT:      name,
		CONCAT_2:    name,
	})
	if err != nil {
		global.Logger.Error(err.Error(), middlewares.ErrLogMsg(ctx)...)
		return nil, errcode.ErrServer
	}
	if len(data) == 0 {
		return &reply.ParamGetFriendsByName{List: []*model.SettingFriend{}}, nil
	}
	result := make([]*model.SettingFriend, 0, len(data))
	for _, v := range data {
		result = append(result, &model.SettingFriend{
			SettingInfo: model.SettingInfo{
				RelationID:   v.RelationID,
				RelationType: string(db.RelationsRelationFriend),
				NickName:     v.NickName,
				IsNotDisturb: v.IsNotDisturb,
				IsPin:        v.IsPin,
				IsShow:       v.IsShow,
				PinTime:      v.PinTime,
				LastShow:     v.LastShow,
			},
			FriendInfo: &model.SettingFriendInfo{
				AccountID: v.AccountID,
				Name:      v.AccountName,
				Avatar:    v.AccountAvatar,
			},
		})
	}
	return &reply.ParamGetFriendsByName{
		List:  result,
		Total: data[0].Total.(int64),
	}, nil
}
