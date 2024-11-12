package dao

import (
	"fmt"
)

// 精准查询-需要字段和值, 返回完美匹配的所有数据切片
// 这里取唯一值的调用示例(根据精准获取用户数据):
// var userData *model.User
// users, _ := dao.UserRepo.FindByFieldExact("id", ID)
// userData = &users[0]
func (br *customRepo[T]) FindByFieldExact(field string, value any) ([]T, error) {
	var result []T
	err := br.DB.Where(field+" = ?", value).Find(&result).Error
	return result, err
}

// 模糊查询-需要字段和值
func (br *customRepo[T]) FindByFieldFuzzy(field string, value string) ([]T, error) {
	var result []T
	err := br.DB.Where(fmt.Sprintf("%s LIKE ?", field), "%"+value+"%").Find(&result).Error
	return result, err
}

// 最简单粗暴的全部查询-返回完整的整表数据
// 暂时不用先放着, 因为使用FindTableData不仅直接获取全部数据还能分页并且也能够全部查询
func (br *customRepo[T]) FindAll() ([]T, error) {
	var result []T
	err := br.DB.Find(new(T)).Error
	return result, err
}

// 集成分页查询、全部查询、整表数据量查询
// 使用示例请全局搜索函数: GetRequestLog
func (br *customRepo[T]) FindTableData(page, limit int) ([]T, int64, error) {
	// 计算总数
	var count int64
	if err := br.DB.Model(new(T)).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	// 如果page和limit小于1，则获取所有数据
	if page < 1 || limit < 1 {
		var allData []T
		if err := br.DB.Find(&allData).Error; err != nil {
			return nil, 0, err
		}
		return allData, count, nil
	}
	// 分页查询
	offset := (page - 1) * limit
	var dataSlice []T
	if err := br.DB.Limit(limit).Offset(offset).Find(&dataSlice).Error; err != nil {
		return nil, 0, err
	}
	return dataSlice, count, nil
}

// 数据量查询, 全部/条件皆可, 如果传入空字符串则返回整表计数
func (br *customRepo[T]) CountTableData(fieldName string, value any) (int64, error) {
	var count int64
	var err error
	if fieldName == "" || value == "" {
		err = br.DB.Model(new(T)).Count(&count).Error
	} else {
		err = br.DB.Model(new(T)).Where(fieldName+" = ?", value).Count(&count).Error
	}
	return count, err
}
