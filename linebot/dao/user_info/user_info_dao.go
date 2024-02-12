package user_info

import "linebot/entity"

type UserInfoDao interface {
	GetUserByLineId(lineId string) *entity.UserInfo
	UpdateUserLastAccess(lineId string) bool
	UpsertInvalidUser(lineId string) bool
}

type emptyUserInfoDao struct{}

func (emptyUserInfoDao) GetUserByLineId(lineId string) *entity.UserInfo { return nil }
func (emptyUserInfoDao) UpdateUserLastAccess(lineId string) bool        { return false }
func (emptyUserInfoDao) UpsertInvalidUser(lineId string) bool           { return false }

func NewEmptyUserInfoDao() *emptyUserInfoDao {
	return &emptyUserInfoDao{}
}
