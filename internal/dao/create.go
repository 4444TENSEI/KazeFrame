package dao

// 新增数据, 需要传入一个结构体, 用于匹配要新建的数据所包含的字段名
// 使用示例请全局搜索函数: CreateUser
func (br *customRepo[T]) Create(entity T) error {
	return br.DB.Create(&entity).Error
}

// 同样是新增数据, 但方式不同, 通过直接传入一个map新建数据, 更加宽松和自由, 需要写死字段名
// 如果字段名没有匹配到数据库表的字段, 则会报错所以自行选择是否使用此方式
func (br *customRepo[T]) CreateEasy(data map[string]interface{}) error {
	return br.DB.Model(new(T)).Create(data).Error
}
