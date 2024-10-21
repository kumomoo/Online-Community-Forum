package models

import "time"

type Post struct {
	ID          int64     `json:"id,string" gorm:"column:post_id"`                             // 帖子ID，唯一
	AuthorID    int64     `json:"author_id" gorm:"column:author_id"`                           // 作者ID
	CommunityID int64     `json:"community_id" gorm:"column:community_id" binding:"required"`  // 所属社区ID
	Status      int32     `json:"status" gorm:"column:status"`                                 // 帖子状态，默认值为 1
	Title       string    `json:"title" gorm:"column:title" binding:"required"`                // 帖子标题
	Content     string    `json:"content" gorm:"column:content" binding:"required"`            // 帖子内容
	CreateTime  time.Time `json:"create_time" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"` // 创建时间，自动生成
}

// 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName       string `json:"author_name"`
	VoteNum          int64  `json:"vote_num"`
	*Post                   //嵌入帖子结构体
	*CommunityDetail `json:"community"`
}
