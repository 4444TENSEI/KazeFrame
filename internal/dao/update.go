package dao

// 单条更新
func (br *customRepo[T]) Update() error {
	return br.DB.Save(new(T)).Error
}

// 条件更新
// 1.指定字段名和值的字典(值如果不唯一会导致更新多条符合的数据) 2.要更新的字段和对应的值
// 使用示例请全局搜索函数: UpdateProfile
func (br *customRepo[T]) UpdateByField(field string, value interface{}, updateList map[string]interface{}) error {
	return br.DB.Model(new(T)).Where(field+" = ?", value).Updates(updateList).Error
}
