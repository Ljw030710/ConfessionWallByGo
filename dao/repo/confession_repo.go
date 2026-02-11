package repo

import (
	"app/dao/model"
	"app/dao/query"
	"context"

	"github.com/zjutjh/mygo/ndb"
)

type ConfessionRepo struct {
	query *query.Query
}

func NewConfessionRepo() *ConfessionRepo{
	return &ConfessionRepo{
		query: query.Use(ndb.Pick()),
	}
}

//根据confesssionID 进行更新
func(r *ConfessionRepo) UpdateByID(
	ctx context.Context,
	id int64,
	receiverName *string,
	content *string,
	imageURL *string,
	isAnonymous *int8,
	status *int8,
)error{
	c := r.query.Confession
	//检查是否存在
	_,err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err !=nil{
		return err
	}
	// 2) 只收集前端传入的字段，未传的不更新
	updates := map[string]interface{}{}
	if receiverName != nil {
		updates["receiver_name"] = *receiverName
	}
	if content != nil {
		updates["content"] = *content
	}
	if imageURL != nil {
		updates["image_url"] = *imageURL
	}
	if isAnonymous != nil {
		updates["is_anonymous"] = *isAnonymous
	}
	if status != nil {
		updates["status"] = *status
	}

	// 3) 如果一个字段都没传，直接返回（由上层决定是否视为参数错误）
	if len(updates) == 0 {
		return nil
	}

	// 4) 执行更新
	_, err = c.WithContext(ctx).Where(c.ID.Eq(id)).Updates(updates)
	return err
}

// Create 新增一条表白记录
func (r *ConfessionRepo) Create(
	ctx context.Context,
	senderID int64,
	receiverName string,
	content string,
	imageURL string,
	isAnonymous int8,
	status int8,
) (*model.Confession, error) {
	c := r.query.Confession

	newConfession := &model.Confession{
		SenderID:     senderID,
		ReceiverName: receiverName,
		Content:      content,
		ImageURL:     imageURL,
		IsAnonymous:  isAnonymous,
		Status:       status,
	}

	err := c.WithContext(ctx).Create(newConfession)
	return newConfession, err
}
//删除表白记录
func (r *ConfessionRepo) DeleteByID(ctx context.Context,id int64)error{
	c := r.query.Confession
	_,err := c.WithContext(ctx).Where(c.ID.Eq(id)).First()
	if err != nil{
		return err
	}
	_,err = c.WithContext(ctx).Where(c.ID.Eq(id)).Update(c.Status,int8(3))
	return err
}
