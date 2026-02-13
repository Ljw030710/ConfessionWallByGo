package repo

import (
	"app/dao/model"
	"app/dao/query"
	"context"
	"fmt"
	"strings"

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
		ParentCommentID: 0,
		Username: username,
		ReplyToUsername: "",
		Content: content,
		
	}
	err = ndb.Pick().WithContext(ctx).Create(newComment).Error
	if err != nil{
		return nil,err
	}
	return newComment,nil
}

//回复评论的操作
func (r *ConfessionCommentRepo) Reply(
	ctx context.Context,
	confessionID int64,
	parentCommentID int64,
	username string,
	replyToUsername string,
	content string,
) (*model.ConfessionComment, error) {
	c := r.query.Confession
	u := r.query.User

	if confessionID <= 0 {
		return nil, fmt.Errorf("invalid confession_id: %d", confessionID)
	}
	if parentCommentID <= 0 {
		return nil, fmt.Errorf("invalid parent_comment_id: %d", parentCommentID)
	}
	if strings.TrimSpace(username) == "" {
		return nil, fmt.Errorf("username is empty")
	}
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("content is empty")
	}

	// 1) 校验表白存在且未删除
	_, err := c.WithContext(ctx).Where(
		c.ID.Eq(confessionID),
		c.Status.Neq(int8(3)),
	).First()
	if err != nil {
		return nil, err
	}

	// 2) 校验回复人存在
	_, err = u.WithContext(ctx).Where(
		u.Username.Eq(username),
	).First()
	if err != nil {
		return nil, err
	}

	// 3) 校验父评论存在、属于同一 confession、且未软删除
	var parent struct {
		ID       int64  `gorm:"column:id"`
		Username string `gorm:"column:username"`
	}
	err = ndb.Pick().WithContext(ctx).
		Table("confession_comments").
		Select("id, username").
		Where("id = ? AND confession_id = ? AND deleted_at = 0", parentCommentID, confessionID).
		First(&parent).Error
	if err != nil {
		return nil, err
	}

	// 4) 没传 replyToUsername 时，默认回复父评论作者
	if strings.TrimSpace(replyToUsername) == "" {
		replyToUsername = parent.Username
	}

	// 5) 校验被回复用户存在
	_, err = u.WithContext(ctx).Where(
		u.Username.Eq(replyToUsername),
	).First()
	if err != nil {
		return nil, err
	}

	// 6) 入库回复评论
	newReply := &model.ConfessionComment{
		ConfessionID:    confessionID,
		ParentCommentID: parentCommentID,
		Username:        username,
		ReplyToUsername: replyToUsername,
		Content:         content,
	}
	err = ndb.Pick().WithContext(ctx).Create(newReply).Error
	if err != nil {
		return nil, err
	}

	return newReply, nil
}