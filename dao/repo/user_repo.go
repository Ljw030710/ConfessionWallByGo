package repo

import (
	"context"
	"errors"

	"github.com/zjutjh/mygo/ndb"
	"gorm.io/gorm"

	"app/dao/model"
	"app/dao/query"
)

type UserRepo struct {
	query *query.Query
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		query: query.Use(ndb.Pick()),
	}
}

// FindById 根据ID查询用户
func (r *UserRepo) FindById(ctx context.Context, id int64) (*model.User, error) {
	u := r.query.User
	record, err := u.WithContext(ctx).Where(u.ID.Eq(id)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}
//登录
func(r *UserRepo) Login(ctx context.Context,username,password string)(*model.User,error){
	u := r.query.User

	// 直接在数据库查询时同时匹配用户名、密码和未删除状态
	user, err := u.WithContext(ctx).Where(
		u.Username.Eq(username),
		u.Password.Eq(password), // 直接比对明文
	).First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 此时可能是用户名不存在，也可能是密码错误
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	return user, nil
}

//注册
// Create 注册用户并返回用户对象
func (r *UserRepo) Create(ctx context.Context, username, password, nickname string) (*model.User, error) {
    u := r.query.User
    newUser := &model.User{
        Username: username,
        Password: password,
        Nickname: nickname,
    }
    
    // 临时增加 Debug() 打印，看看控制台有没有输出 INSERT 语句
    err := u.WithContext(ctx).Debug().Create(newUser) 
    return newUser, err
}
//FindByUsername
func (r *UserRepo) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	u := r.query.User
	record, err := u.WithContext(ctx).Where(u.Username.Eq(username)).First()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return record, err
}
//根据用户ID更新密码
func(r *UserRepo) UpdatePasswordByID(ctx context.Context,id int64,newPassword string)error{
	u:=r.query.User
	_,err :=u.WithContext(ctx).Where(u.ID.Eq(id)).Update(u.Password,newPassword)
	return  err
}
//根据username进行
func(r *UserRepo) UpdatePasswordByUsername(ctx context.Context,username string,newPassword string)error{
	u := r.query.User
	_,err := u.WithContext(ctx).Where(u.Username.Eq(username)).Update(u.Password,newPassword)
	return err
}
//更新nickname
func(r *UserRepo) UpdateNicknameByUsername(ctx context.Context,username string,newNickname string)error{
	u:= r.query.User
	_, err := u.WithContext(ctx).Where(u.Username.Eq(username)).Update(u.Nickname, newNickname)
	return err
}