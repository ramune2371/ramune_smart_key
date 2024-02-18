package user_info

import "linebot/entity"

type UserInfoDao interface {
	GetUserByLineId(lineId string) (*entity.UserInfo, error)
	UpdateUserLastAccess(lineId string) (bool, error)
	UpsertInvalidUser(lineId string) (bool, error)
}

type emptyUserInfoDao struct{}

func (emptyUserInfoDao) GetUserByLineId(lineId string) (*entity.UserInfo, error) { return nil, nil }
func (emptyUserInfoDao) UpdateUserLastAccess(lineId string) (bool, error)        { return false, nil }
func (emptyUserInfoDao) UpsertInvalidUser(lineId string) (bool, error)           { return false, nil }

func NewEmptyUserInfoDao() *emptyUserInfoDao {
	return &emptyUserInfoDao{}
}
