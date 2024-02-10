package user_info

import "linebot/entity"

type UserInfoDao interface {
	GetUserByLineId(lineId string) *entity.UserInfo
	UpdateUserLastAccess(lineId string) bool
	UpsertInvalidUser(lineId string) bool
}
