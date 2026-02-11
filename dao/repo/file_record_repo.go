package repo

import (
	"app/dao/model"
	"app/dao/query"
	"context"
	"errors"

	"github.com/zjutjh/mygo/ndb"
	"gorm.io/gorm"
)

type FileRecordRepo struct {
	query *query.Query
}

// NewFileRecordRepo 初始化
func NewFileRecordRepo() *FileRecordRepo {
	return &FileRecordRepo{
		query: query.Use(ndb.Pick()),
	}
}

// FindByHash 根据文件 MD5 哈希值查找记录
// 返回值：*model.FileRecord (记录指针), error (错误)
func (r *FileRecordRepo) FindByHash(ctx context.Context, hash string) (*model.FileRecord, error) {
	f := r.query.FileRecord
	
	// 使用 GORM Gen 生成的查询方法
	record, err := f.WithContext(ctx).Where(f.FileHash.Eq(hash)).First()
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 没找到记录，属于正常业务逻辑
		}
		return nil, err // 数据库连接等异常
	}
	
	return record, nil
}

// Create 插入新的文件去重记录
func (r *FileRecordRepo) Create(ctx context.Context, hash string, url string) error {
	f := r.query.FileRecord
	
	newRecord := &model.FileRecord{
		FileHash: hash,
		FileURL:  url, // 匹配你生成的 model 中的 FileURL 字段
	}
	
	return f.WithContext(ctx).Create(newRecord)
}