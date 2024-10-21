package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignup) (err error) {
	//1.判断用户存不存在
	var exist bool
	exist, err = mysql.CheckUserExist(p.Username)
	if err != nil {
		//数据库查询出错
		return err
	}
	if exist {
		//用户已存在的错误
		return mysql.ErrorUserExist
	}
	//2.生成UID
	userID := snowflake.GenID()
	//构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3.保存进数据库
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	//生成jwt
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
