package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	//1.获取参数与参数校验
	var p models.ParamSignup
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//请求参数有误
		ResponseError(c, CodeInvalidParam)
		return
	}
	//手动对请求参数进行详细的业务规则校验
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.RePassword != p.Password {
	// 	zap.L().Error("SignUp with invalid param")
	// 	//请求参数有误
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }
	fmt.Println(p)

	//2.业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	//3.返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	//1.获取参数与参数校验
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		//请求参数有误
		ResponseError(c, CodeInvalidParam)
		return
	}
	fmt.Println(p)

	//2.业务处理
	user, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	//3.返回响应
	ResponseSuccess(c, gin.H{
		"user_id":   fmt.Sprintf("%d", user.UserID),
		"user_name": user.Username,
		"token":     user.Token,
	})
}
