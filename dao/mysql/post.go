package mysql

import (
	"bluebell/models"
	"strings"

	"gorm.io/gorm"
)

func CreatePost(p *models.Post) (err error) {
	// 使用 GORM 创建记录
	if err = db.Create(p).Error; err != nil {
		return err
	}
	return nil
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)

	// 使用 GORM 通过主键查找记录
	if err = db.Where("post_id = ?", pid).First(post).Error; err != nil {
		return nil, err
	}
	return post, nil
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	posts = make([]*models.Post, 0)

	// 使用 GORM 进行分页查询
	if err = db.Limit(int(size)).Offset(int((page - 1) * size)).Find(&posts).Error; err != nil {
		return nil, err
	}

	return posts, nil
}

// 根据给定的id列表查询帖子数据
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	// 构建查询
	err = db.Where("post_id IN ?", ids).
		Order(gorm.Expr("FIND_IN_SET(post_id, ?)", strings.Join(ids, ","))).
		Find(&postList).Error

	if err != nil {
		return nil, err
	}

	return postList, nil
}
