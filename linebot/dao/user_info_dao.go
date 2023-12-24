package dao

import (
	"linebot/entity"
	"time"

	"github.com/google/uuid"
)

/*
LineのIDを元に、ユーザレコードを取得
*/
func (g *GormDB) GetUserByLineId(lineId string) *entity.UserInfo {
	var ret *entity.UserInfo

	g.getTable(entity.UserInfoTable).Where("line_id = ?", lineId).Find(&ret)
	return ret
}

/*
LineのIDを元に、最終アクセス時間を更新
UI-A-01
*/
func (g *GormDB) UpdateUserLastAccess(lineId string) bool {

	table := g.getTable(entity.UserInfoTable).Where("line_id = ?", lineId).Update("last_access", time.Now())
	if table == nil {
		return false
	} else {
		return true
	}
}

/*
LineのIDを元に、不正なユーザレコードを作成 or 最終アクセス時間を更新
UI-E-01
*/
func (g *GormDB) UpsertInvalidUser(lineId string) bool {

	tx := g.getTable(entity.UserInfoTable)

	if tx == nil {
		return false
	}

	tx.Begin()
	var ret *entity.UserInfo
	tx.Where("line_id = ?", lineId).Find(&ret)
	if ret.UserUuid == "" {
		ret.UserUuid = uuid.New().String()
	}
	ret.LineId = lineId
	if ret.UserName == "" {
		ret.UserName = "Unknown"
	}
	tx.Save(&ret)
	tx.Commit()
	return true
}
