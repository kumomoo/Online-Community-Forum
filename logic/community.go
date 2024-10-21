package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	//查找到所有的community并返回
	return mysql.GetCommunityList()
}

// 查询社区详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	//查找到所有的community并返回
	return mysql.GetCommunityDetailByID(id)
}
