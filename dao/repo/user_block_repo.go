package repo

import (
	"app/dao/model"
	"app/dao/query"
	"context"
	"errors"

	"github.com/zjutjh/mygo/ndb"
	"gorm.io/gorm"
)

type UserBlockRepo struct {
	query *query.Query
}

//初始化
func NewUserBlockRepo() *UserBlockRepo{
	return &UserBlockRepo{
		query: query.Use(ndb.Pick()),
	}
}
//根据用户id拉黑
func (r *UserBlockRepo) Block(ctx context.Context, blockerID, blockedID int64) error {
	// 不能拉黑自己
	if blockerID == blockedID {
		return errors.New("不能拉黑自己")
	}

	ub := r.query.UserBlock
	record := &model.UserBlock{
		BlockerID: blockerID,
		BlockedID: blockedID,
	}

	// 重复拉黑时会直接返回数据库唯一键冲突错误
	return ub.WithContext(ctx).Create(record)
}

// Unblock 按用户ID取消拉黑（不存在时返回 gorm.ErrRecordNotFound）
func (r *UserBlockRepo) Unblock(ctx context.Context, blockerID, blockedID int64) error {
	ub := r.query.UserBlock

	// 先查是否存在关系
	_, err := ub.WithContext(ctx).
		Where(
			ub.BlockerID.Eq(blockerID),
			ub.BlockedID.Eq(blockedID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound
		}
		return err
	}

	// 存在则删除
	_, err = ub.WithContext(ctx).
		Where(
			ub.BlockerID.Eq(blockerID),
			ub.BlockedID.Eq(blockedID),
		).
		Delete()
	return err
}

// IsBlocked 按用户ID查询是否已拉黑
func (r *UserBlockRepo) IsBlocked(ctx context.Context, blockerID, blockedID int64) (bool, error) {
	ub := r.query.UserBlock

	_, err := ub.WithContext(ctx).
		Where(
			ub.BlockerID.Eq(blockerID),
			ub.BlockedID.Eq(blockedID),
		).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//BlockByUsername按用户名进行拉黑
func (r *UserBlockRepo)BlockByUsername(ctx context.Context,blockerUsername,blockedUsername string)error{
	if blockerUsername == blockedUsername {
	return errors.New("不能拉黑自己")
	}
	userRepo := NewUserRepo()
	blocker,err := userRepo.FindByUsername(ctx,blockerUsername)
	if err != nil{
		return err
	}
	if blocker == nil{
		return gorm.ErrRecordNotFound
	}
	blocked, err := userRepo.FindByUsername(ctx, blockedUsername)
	if err != nil {
		return err
	}
	if blocked == nil {
		return gorm.ErrRecordNotFound
	}

	return r.Block(ctx, blocker.ID, blocked.ID)
}



// UnblockByUsername 按用户名取消拉黑（内部转ID）
func (r *UserBlockRepo) UnblockByUsername(ctx context.Context, blockerUsername, blockedUsername string) error {
	userRepo := NewUserRepo()

	blocker, err := userRepo.FindByUsername(ctx, blockerUsername)
	if err != nil {
		return err
	}
	if blocker == nil {
		return gorm.ErrRecordNotFound
	}

	blocked, err := userRepo.FindByUsername(ctx, blockedUsername)
	if err != nil {
		return err
	}
	if blocked == nil {
		return gorm.ErrRecordNotFound
	}

	return r.Unblock(ctx, blocker.ID, blocked.ID)
}


// IsBlockedByUsername 按用户名查询是否已拉黑（内部转ID）
func (r *UserBlockRepo) IsBlockedByUsername(ctx context.Context, blockerUsername, blockedUsername string) (bool, error) {
	userRepo := NewUserRepo()

	blocker, err := userRepo.FindByUsername(ctx, blockerUsername)
	if err != nil {
		return false, err
	}
	if blocker == nil {
		return false, gorm.ErrRecordNotFound
	}

	blocked, err := userRepo.FindByUsername(ctx, blockedUsername)
	if err != nil {
		return false, err
	}
	if blocked == nil {
		return false, gorm.ErrRecordNotFound
	}

	return r.IsBlocked(ctx, blocker.ID, blocked.ID)
}