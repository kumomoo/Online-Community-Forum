package models

import "time"

type User struct {
	ID         int64     `gorm:"primaryKey;autoIncrement"`                                             // 主键，自增
	UserID     int64     `gorm:"uniqueIndex:idx_user_id"`                                              // user_id 字段，唯一索引
	Username   string    `gorm:"type:varchar(64);uniqueIndex:idx_username"`                            // username 字段，唯一索引
	Password   string    `gorm:"type:varchar(64)"`                                                     // 密码字段
	Email      string    `gorm:"type:varchar(64)"`                                                     // 邮箱字段，允许为空
	Gender     int8      `gorm:"type:tinyint(4);default:0"`                                            // 性别字段，默认值为 0
	CreateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`                             // 创建时间，默认当前时间
	UpdateTime time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"` // 更新时间，默认当前时间，且在更新时自动更新
	Token      string
}
