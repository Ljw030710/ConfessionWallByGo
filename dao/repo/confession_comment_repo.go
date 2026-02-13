package repo

import (
	"app/dao/model"
	"app/dao/query"
	"context"

	"github.com/zjutjh/mygo/ndb"
)

type ConfessionCommentRepo struct {
	query *query.Query
}

func NewConfessionCommentRepo()*ConfessionCommentRepo{
	return &ConfessionCommentRepo{
		query: query.Use(ndb.Pick()),
	}
}

//创建评论
func (r *ConfessionCommentRepo)Create(
	ctx context.Context,
	confessionID int64,
	username string,
	content string,
)(*model.ConfessionComment,error){
	c := r.query.Confession
	u:= r.query.User
	_,err := c.WithContext(ctx).Where(
		c.ID.Eq(confessionID),
		c.Status.Neq(int8(3)),
	).First()
	if err != nil{
		return nil,err
	}
	_,err =  u.WithContext(ctx).Where(
		u.Username.Eq(username),
	).First()
	if err != nil{
		return nil,err
	}

	newComment := &model.ConfessionComment{
		ConfessionID: confessionID,
		Username: username,
		Content: content,
	}
	err = ndb.Pick().WithContext(ctx).Create(newComment).Error
	if err != nil{
		return nil,err
	}
	return newComment,nil
}
