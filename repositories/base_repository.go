//@Title base_repository.go
//@Description 仓库基本包。实现部分重复代码
//@Author tb

package repositories

import "gorm.io/gorm"

//IBaseRepo 基础仓库接口
type IBaseRepo interface {
	//Add 增
	//@model 模型的接口类型
	Add(model interface{}) error

	//DeleteByID 软删除
	//@model 模型的接口类型
	//@id 需要删除的数据的id 因为本系统认为默认都使用了gorm.Model所以使用uint, 也可以不使用，转为uint也可以
	DeleteByID(model interface{}, id uint) error

	//Update 修改数据
	//@model 模型的接口类型
	//@columns 需要修改的字段，不填则为默认保存模型中的更新非零值
	Update(model interface{}, columns ...string) error

	//根据id查找
	//@model 模型的接口类型
	//@id 需要查找的数据的id 因为本系统认为默认都使用了gorm.Model所以使用uint, 也可以不使用，转为uint也可以
	//@preloads 需要预加载的字段
	FindByID(model interface{}, id uint, preloads ...string) error
}

//BaseRepository 定义的基础结构体 用于继承
type BaseRepository struct {
	Db *gorm.DB
}


//Add 增
//@model 模型的接口类型
func (b *BaseRepository) Add(model interface{}) error {
	return b.Db.Create(model).Error
}

//DeleteByID 软删除
//@model 模型的接口类型
//@id 需要删除的数据的id 因为本系统认为默认都使用了gorm.Model所以使用uint, 也可以不使用，转为uint也可以
func (b *BaseRepository) DeleteByID(model interface{}, id uint) error {
	return b.Db.Delete(model, id).Error
}

//Update 修改数据
//@model 模型的接口类型
//@columns 需要修改的字段，不填则为默认保存模型中的更新非零值
func (b *BaseRepository) Update(model interface{}, columns ...string) error {
	if len(columns) == 0 {
		return b.Db.Model(model).Updates(model).Error
	}
	updateFields := make([]interface{}, len(columns))
	for k := range columns {
		updateFields[k] = columns[k]
	}
	return b.Db.Model(model).Select(updateFields).Updates(model).Error
}

//根据id查找
//@model 模型的接口类型
//@id 需要查找的数据的id 因为本系统认为默认都使用了gorm.Model所以使用uint, 也可以不使用，转为uint也可以
func (b *BaseRepository) FindByID(model interface{}, id uint, preloads ...string) error {
	db := b.Db
	//需要预加载
	if len(preloads) > 0 {
		for k := range preloads {
			db = db.Preload(preloads[k])
		}
	}

	return db.First(model).Error
}
