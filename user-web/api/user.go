package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/middlewares"
	"mxshop-api/user-web/models"
	"mxshop-api/user-web/proto"
	"mxshop-api/user-web/request"
	"net/http"
	"strings"
	"time"
)

func Register(c *gin.Context) {
	// 用户绑定请求参数
	var userRegister = new(request.RegisterForm)
	if err := c.ShouldBind(&userRegister); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 验证码

	// 请求服务
	user, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: userRegister.Mobile,
		PassWord: userRegister.PassWord,
		Mobile:   userRegister.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】失败: %s", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 生成用户token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "codelpj",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

func PassWordLogin(c *gin.Context) {

	// 登录参数校验
	var userLogin = new(request.PassWordLoginForm)
	if err := c.ShouldBind(&userLogin); err != nil {
		HandleValidatorError(c, err)
		return
	}
	// 验证码校验

	// 登录逻辑
	// 用户是否存在
	user, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: userLogin.Mobile})
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"mobile": "用户不存在",
		})
		return
	}
	// 用户校验密码
	passRes, err := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
		Password:          userLogin.PassWord,
		EncryptedPassword: user.PassWord,
	})
	if err != nil || !passRes.Success {
		c.JSON(http.StatusBadRequest, map[string]string{
			"password": "登录失败",
		})
	}
	// 生成用户token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "codelpj",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"nick_name":  user.NickName,
		"token":      token,
		"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

func GetUserDetail(c *gin.Context) {
	claims, _ := c.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问用户: %d", currentUser.ID)
	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.IdRequest{
		Id: int32(currentUser.ID),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name":     rsp.NickName,
		"birthday": time.Unix(int64(rsp.BirthDay), 0).Format("2006-01-02"),
		"gender":   rsp.Gender,
		"mobile":   rsp.Mobile,
	})
}
func UpdateUser(c *gin.Context) {

	var userUpdate = new(request.UpdateUserForm)
	if err := c.ShouldBind(&userUpdate); err != nil {
		HandleValidatorError(c, err)
		return
	}

	claims, _ := c.Get("claims")
	currentUser := claims.(*models.CustomClaims)

	//将前端传递过来的日期格式转换成int
	loc, _ := time.LoadLocation("Local") //local的L必须大写
	birthDay, _ := time.ParseInLocation("2006-01-02", userUpdate.Birthday, loc)
	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(currentUser.ID),
		NickName: userUpdate.Name,
		Gender:   userUpdate.Gender,
		BirthDay: uint64(birthDay.Unix()),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, gin.H{})

}
func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	//将grpc的code转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}
