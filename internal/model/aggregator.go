package model

import "gorm.io/gorm"

// 聚合器Aggregator对应的规则数据库字段信息
type Aggregator struct {
	Id         uint32 `gorm:"primaryKey;column:id" bson:"id"`
	Name       string `gorm:"column:name" bson:"name"`
	Metric     string `gorm:"column:metric" bson:"metric"`
	FunctionId uint32 `gorm:"column:function_id" bson:"function_id"`
	RuleId     uint32 `gorm:"column:rule_id" bson:"rule_id"`
}

// 创建聚合器至规则数据库
func (a Aggregator) Create(db *gorm.DB) error {
	return db.Create(&a).Error
}

// 删除规则数据库中特定聚合器
func (a Aggregator) Delete(db *gorm.DB) error {
	return db.Where("id = ?", a.Id).Delete(&a).Error
}

// 获取规则数据库中特定聚合器
func (a Aggregator) Get(db *gorm.DB) (Aggregator, error) {
	var Aggregator Aggregator
	err := db.Where("id = ?", a.Id).First(&Aggregator).Error
	return Aggregator, err

}

// 更新规则数据库中特定聚合器
func (a Aggregator) Update(db *gorm.DB) error {
	var temp Aggregator
	temp.Id = a.Id
	return db.Model(&temp).Updates(a).Error
}

// 根据指标返回规则数据库所有聚合器
func (a Aggregator) Search(db *gorm.DB, metric string) ([]Aggregator, error) {
	var result []Aggregator
	err := db.Where("metric = ?", metric).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 返回规则数据库所有聚合器
func (a Aggregator) List(db *gorm.DB) ([]Aggregator, error) {
	var aggregators []Aggregator
	err := db.Find(&aggregators).Error
	if err != nil {
		return nil, err
	}
	return aggregators, nil
}
