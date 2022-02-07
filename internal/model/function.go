package model

import "gorm.io/gorm"

// 聚合函数Function对应的规则数据库字段信息
type Function struct {
	Id          uint32  `gorm:"primaryKey;column:id" bson:"id"`
	Type        string  `gorm:"column:type" bson:"type"`
	Threshold   float64 `gorm:"column:threshold" bson:"threshold"`
	Description string  `gorm:"description" bson:"description"`
}

// 创建聚合函数至规则数据库
func (f Function) Create(db *gorm.DB) error {
	return db.Create(&f).Error
}

// 删除规则数据库中特定聚合函数
func (f Function) Delete(db *gorm.DB) error {
	return db.Where("id = ?", f.Id).Delete(&f).Error
}

// 获取规则数据库中特定聚合函数
func (f Function) Get(db *gorm.DB) (Function, error) {
	var Function Function
	err := db.Where("id = ?", f.Id).First(&Function).Error
	return Function, err

}

// 更新规则数据库中特定聚合函数
func (f Function) Update(db *gorm.DB) error {
	var temp Function
	temp.Id = f.Id
	return db.Model(&temp).Updates(f).Error
}

// 获取规则数据库中所有聚合函数
func (f Function) List(db *gorm.DB) ([]Function, error) {
	var functions []Function
	err := db.Find(&functions).Error
	if err != nil {
		return nil, err
	}
	return functions, nil
}
