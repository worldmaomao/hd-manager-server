package service

import (
	"crypto/md5"
	"fmt"
	"github.com/sarulabs/di"
	"log"
	config "worldmaomao/harddisk/internal/config"
	"worldmaomao/harddisk/internal/constant"
	"worldmaomao/harddisk/internal/dao"
	"worldmaomao/harddisk/internal/dao/model"
	"worldmaomao/harddisk/internal/rest/middlewares"
)

type UserService struct {
	userDao   *dao.UserDao
	container di.Container
}

func NewUserService(container di.Container) *UserService {
	return &UserService{
		container: container,
		userDao:   container.Get(constant.UserDao).(*dao.UserDao),
	}
}

func (service *UserService) Initial() error {
	var (
		username = "admin"
		password = "123456"
	)
	user, _ := service.userDao.GetByUsername(username)
	if len(user.Username) != 0 {
		return nil
	}
	log.Println("to add user[" + username + "]")
	return service.userDao.Add(model.User{
		Username: username,
		Password: service.md5(password + username),
		Roles:    []string{constant.ADMIN},
	})
}

func (service *UserService) md5(password string) string {
	md5Hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", md5Hash)
}

func (service *UserService) Login(username string, password string, platform string) (string, error) {
	if len(password) == 0 || len(username) == 0 || len(platform) == 0 {
		return "", fmt.Errorf("参数错误")
	}
	user, err := service.userDao.GetByUsername(username)
	if err != nil {
		log.Println(err.Error())
		return "", fmt.Errorf("登录失败,用户名或者密码错误")
	}
	if user.Password == service.md5(password+username) {
		configuration := service.container.Get(constant.Configuration).(*config.Configuration)
		return middlewares.GenerateJWTToken(*user, configuration.GetJwtKey(), platform), nil
	} else {
		log.Println("password is wrong")
		return "", fmt.Errorf("登录失败,用户名或者密码错误")
	}
	return "", nil
}
