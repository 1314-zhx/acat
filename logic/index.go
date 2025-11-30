package logic

import (
	"acat/dao"
	"acat/dao/db"
	"context"
	"go.uber.org/zap"
)

func Check(ctx context.Context) ([][]int64, error) {
	indexDao := dao.NewIndexDao(db.DB)
	schedule, err := indexDao.Check(ctx) // 注意顺序：(data, error)
	if err != nil {
		zap.L().Error("logic/index.go Check() failed", zap.Error(err))
		return nil, err // 查询失败才是 error
	}
	// 即使 schedule 为空，也返回成功
	return schedule, nil
}
