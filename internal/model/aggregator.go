package model

import "gorm.io/gorm"

// 聚合函数和告警规则分来存，保证扩展性
// Name要unique
type Aggregator struct {
	Id         uint32 `gorm:"primaryKey;column:id" json:"id"`
	Name       string `gorm:"column:name" json:"name"`
	Metric     string `gorm:"column:metric" json:"metric"`
	FunctionId uint32 `gorm:"column:function_id" json:"function_id"`
	RuleId     uint32 `gorm:"column:rule_id" json:"rule_id"`
}

func (a Aggregator) Create(db *gorm.DB) error {
	return db.Create(&a).Error
}

func (a Aggregator) Delete(db *gorm.DB) error {
	return db.Where("id = ?", a.Id).Delete(&a).Error
}

func (a Aggregator) Get(db *gorm.DB) (Aggregator, error) {
	var Aggregator Aggregator
	err := db.Where("id = ?", a.Id).First(&Aggregator).Error
	return Aggregator, err

}

func (a Aggregator) Update(db *gorm.DB) error {
	// 根据 `struct` 更新属性，只会更新非零值的字段
	// db.Model(&user).Updates(User{Name: "hello", Age: 18, Active: false})
	// UPDATE users SET name='hello', age=18, updated_at = '2013-11-17 21:34:10' WHERE id = 111;
	var temp Aggregator
	temp.Id = a.Id
	return db.Model(&temp).Updates(a).Error
}

// 返回该指标的所有聚合器
func (a Aggregator) Search(db *gorm.DB, metric string) ([]Aggregator, error) {
	var result []Aggregator
	// SELECT * FROM users WHERE name LIKE '%jin%';
	err := db.Where("metric = ?", metric).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 返回所有聚合器
func (a Aggregator) List(db *gorm.DB) ([]Aggregator, error) {
	// SELECT * FROM users WHERE name LIKE '%jin%';
	var aggregators []Aggregator
	err := db.Find(&aggregators).Error
	if err != nil {
		return nil, err
	}
	return aggregators, nil

}
