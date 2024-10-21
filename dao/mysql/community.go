package mysql

import (
	"bluebell/models"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	err = db.Select("community_id", "community_name").Find(&communityList).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.L().Warn("there is no community in db")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)

	// 使用 GORM 通过主键查找
	if err = db.Table("communities").Where("community_id = ?", id).First(&community).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrorInvalidID
		}
	}
	return community, err
}
