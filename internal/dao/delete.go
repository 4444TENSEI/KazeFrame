// 删除操作相关, 默认全都是硬删除, 想要触发软删除的话, 去掉函数内的.Unscoped()即可
package dao

import (
	"KazeFrame/internal/model"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// 批量条件删除-返回详细删除成功失败数据
// 使用示例请全局搜索函数: DeleteUser
func (br *customRepo[T]) QuickHardDelete(fieldName string, values []interface{}) (model.DeleteUserResponse, error) {
	var response model.DeleteUserResponse
	for _, value := range values {
		query := br.DB.Unscoped().Where(fieldName+" = ?", value)
		result := query.Delete(new(T))
		switch {
		case result.Error != nil:
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				response.ErrCount++
				response.ErrValue = append(response.ErrValue, fmt.Sprintf("%v", value))
			} else {
				return response, result.Error
			}
		case result.RowsAffected == 0:
			response.ErrCount++
			response.ErrValue = append(response.ErrValue, fmt.Sprintf("%v", value))
		default:
			response.OkCount++
			response.OkValue = append(response.OkValue, fmt.Sprintf("%v", value))
		}
	}
	return response, nil
}

// 批量条件删除简易版, 使用方式和QuickHardDelete一样, 但是没有详细反馈(现在应该没有不需要详细反馈的情况吧)
func (br *customRepo[T]) HardDelete(fieldName string, values []interface{}) error {
	query := br.DB.Unscoped().Where(fieldName+" IN ?", values)
	return query.Delete(new(T)).Error
}

// 清空数据表所有数据, 危险操作请谨慎使用, 需要传递布尔值true代表永久删除, false代表软删除
// 使用示例请全局搜索函数: ClearAllRequestLog
func (br *customRepo[T]) ClearAllData(hardDeleteEnable bool) (int64, error) {
	db := br.DB
	if hardDeleteEnable {
		db = db.Unscoped()
	}
	var count int64
	if err := db.Model(new(T)).Count(&count).Error; err != nil {
		return 0, err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(new(T)).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 删除某个时间点之前的记录, 需要指定datetime类型的字段名以及一个时间点
// 使用示例全局搜索函数: DeleteRequestLogBytime
func (br *customRepo[T]) DeleteByTime(fieldName string, beforeTime time.Time) (int64, error) {
	query := br.DB.Unscoped().Where(fmt.Sprintf("%s < ?", fieldName), beforeTime)
	result := query.Delete(new(T))
	return result.RowsAffected, result.Error
}
