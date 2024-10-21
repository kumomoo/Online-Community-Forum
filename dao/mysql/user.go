package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

const secret = "kumomo"

func CheckUserExist(username string) (bool, error) {
	var count int64
	if err := db.Model(&models.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		fmt.Println("查询失败")
		return false, err
	}
	return count > 0, nil
}

func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)

	// 使用 GORM 插入用户
	if err = db.Create(&user).Error; err != nil {
		fmt.Println("插入失败")
		return err
	}
	return nil
}

func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户输入的密码

	err = db.Where("username = ?", user.Username).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorUserNotExist
		}
		return err // 查询数据库时出错
	}

	// 判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserById 根据id获取用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)

	// 使用 GORM 查询
	if err = db.Where("user_id = ?", uid).First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}
