package repository

import (
	"context"
	"github.com/38888/nunu-layout-advanced/internal/model"
	"github.com/38888/nunu-layout-advanced/internal/model/dao"
	"github.com/38888/nunu-layout-advanced/pkg/pagination"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
	FindOne(ctx context.Context, id int64) (*model.User, error)
	FindOneByOther(ctx context.Context, do *model.User) (*model.User, error)
	All(ctx context.Context, do *model.User) ([]*model.User, error)
	List(ctx context.Context, page, size int, do *model.User) ([]*model.User, error)
	Count(ctx context.Context, do *model.User) (int64, error)
}

func NewUserRepository(
	r *Repository,
) UserRepository {
	return &userRepository{
		Repository: r,
	}
}

type userRepository struct {
	*Repository
}

// Create 方法用于在数据库中创建一个新的记录
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 do 是指向 model.User 结构体的指针，包含了要创建的的相关数据
// 返回值 error 如果创建过程出现错误则返回相应的错误信息，否则返回 nil
func (r *userRepository) Create(ctx context.Context, do *model.User) error {
	user := r.DB(ctx).User
	err := user.WithContext(ctx).Create(do)
	if err != nil {
		return err
	}
	return nil
}

// Update 方法用于更新数据库中指定的记录信息
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 do 是指向 model.User 结构体的指针，包含了要更新的的相关数据（通过其 ID 来定位具体要更新的记录）
// 返回值 error 如果更新过程出现错误则返回相应的错误信息，若更新影响的行数小于 1（即未找到对应的记录）则返回 gorm.ErrRecordNotFound，否则返回 nil
func (r *userRepository) Update(ctx context.Context, do *model.User) error {
	user := r.DB(ctx).User
	resultInfo, err := user.WithContext(ctx).Where(user.ID.Eq(do.ID)).Updates(do)
	if err != nil {
		return err
	}
	if resultInfo.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete 方法用于从数据库中删除指定 ID 的记录
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 id 是要删除的记录的 ID（类型为 int64）
// 返回值 error 如果删除过程出现错误则返回相应的错误信息，若删除影响的行数小于 1（即未找到对应的记录）则返回 gorm.ErrRecordNotFound，否则返回 nil
func (r *userRepository) Delete(ctx context.Context, id int64) error {
	user := r.DB(ctx).User
	resultInfo, err := user.WithContext(ctx).Where(user.ID.Eq(id)).Delete()
	if err != nil {
		return err
	}
	if resultInfo.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindOne 方法用于从数据库中查找并获取指定 ID 的单个记录
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 id 是要查找的记录的 ID（类型为 int64）
// 返回值 *model.User 是指向找到的记录对应的 model.User 结构体的指针，如果未找到则返回 nil；error 如果查找过程出现错误则返回相应的错误信息，否则返回 nil
func (r *userRepository) FindOne(ctx context.Context, id int64) (*model.User, error) {
	var resp *model.User

	user := r.DB(ctx).User
	resp, err := user.WithContext(ctx).Where(user.ID.Eq(id)).Take()

	if err != nil {
		return nil, err
	}
	return resp, nil
}

// whereBuilder 方法用于构建查询条件的基础构建器，根据传入的 model.User 结构体指针来设置相应的查询条件（目前仅展示了对 ID 的条件设置，可能还有其他条件待补充）
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 do 是指向 model.User 结构体的指针，可从中获取构建查询条件的相关数据
// 返回值 dao.IUserDo 构建好的带有查询条件的基础构建器，后续可基于此继续添加更多查询相关操作
func (r *userRepository) whereBuilder(ctx context.Context, do *model.User) dao.IUserDo {
	user := r.DB(ctx).User
	qb := user.WithContext(ctx)
	if do != nil {
		if do.ID != 0 {
			qb = qb.Where(user.ID.Eq(do.ID))
		}
		//other where

	}
	return qb
}

// FindOneByOther 方法通过自定义的条件（由 whereBuilder 方法构建的条件）来查找并获取单个记录
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 do 是指向 model.User 结构体的指针，包含了构建查询条件的相关数据（用于定位具体要查找的记录）
// 返回值 *model.User 是指向找到的记录对应的 model.User 结构体的指针，如果未找到则返回 nil；error 如果查找过程出现错误则返回相应的错误信息，否则返回 nil
func (r *userRepository) FindOneByOther(ctx context.Context, do *model.User) (*model.User, error) {
	var resp *model.User
	qb := r.whereBuilder(ctx, do)
	resp, err := qb.Take()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// All 方法用于获取满足特定条件（由 whereBuilder 方法构建的条件）的所有记录列表
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 do 是指向 model.User 结构体的指针，包含了构建查询条件的相关数据
// 返回值 []*model.User 满足条件的记录列表（是一个指向 model.User 结构体的指针切片），error 如果获取过程出现错误则返回相应的错误信息，否则返回 nil
func (r *userRepository) All(ctx context.Context, do *model.User) ([]*model.User, error) {
	qb := r.whereBuilder(ctx, do)
	list, err := qb.
		Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// List 方法用于获取满足特定条件（由 whereBuilder 方法构建的条件）的记录列表，支持分页和指定每页数量的功能
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 page 表示页码（从 1 开始计数），用于分页获取数据
// 参数 size 表示每页显示的记录数量
// 参数 do 是指向 model.User 结构体的指针，包含了构建查询条件的相关数据
// 返回值 []*model.User 满足条件的记录列表（是一个指向 model.User 结构体的指针切片），error 如果获取过程出现错误则返回相应的错误信息，否则返回 nil
func (r *userRepository) List(ctx context.Context, page, size int, do *model.User) ([]*model.User, error) {
	qb := r.whereBuilder(ctx, do)
	if size != 0 {
		qb = qb.Limit(size)
	}
	if size != 0 && page != 0 {
		qb = qb.Offset(pagination.GetPageOffset(page, size))
	}
	list, err := qb.
		Find()
	if err != nil {
		return nil, err
	}
	return list, nil
}

// Count 方法用于统计满足特定条件（由 whereBuilder 方法构建的条件）的记录数量
// 参数 ctx 是上下文信息，用于控制操作的超时等情况
// 参数 do 是指向 model.User 结构体的指针，包含了构建查询条件的相关数据
// 返回值 int64 满足条件的记录数量，error 如果统计过程出现错误则返回相应的错误信息，否则返回 nil
func (r *userRepository) Count(ctx context.Context, do *model.User) (int64, error) {
	qb := r.whereBuilder(ctx, do)
	count, err := qb.Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}
