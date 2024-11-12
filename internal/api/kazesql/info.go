package kazesql

import (
	"KazeFrame/pkg/util"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 包含数据库连接所需的信息
type DBConnectInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// 根据数据库连接信息生成DSN
func dsnFromDBConnectInfo(info DBConnectInfo) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		info.User, info.Password, info.Host, info.Port, info.Database)
}

// 包含表的字段信息
type TableInfo struct {
	TableName string       `json:"table_name"`
	Columns   []ColumnInfo `json:"columns"`
	RowCount  int64        `json:"row_count"`
}

// 包含字段名、类型和长度
type ColumnInfo struct {
	ColumnName string `json:"column_name"`
	ColumnType string `json:"column_type"`
}

// 获取数据库中的所有表信息
func getDatabaseTableNamesAndInfo(db *gorm.DB) ([]TableInfo, error) {
	var tables []TableInfo
	rows, err := db.Raw("SHOW TABLES").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tableInfo, err := getTableInfo(db, tableName)
		if err != nil {
			return nil, err
		}
		tables = append(tables, *tableInfo)
	}
	return tables, nil
}

// 获取单个表的信息
func getTableInfo(db *gorm.DB, tableName string) (*TableInfo, error) {
	var columns []ColumnInfo
	rows, err := db.Raw("SHOW COLUMNS FROM `" + tableName + "`").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var columnName, columnType string
		if err := rows.Scan(&columnName, &columnType); err != nil {
			return nil, err
		}
		columns = append(columns, ColumnInfo{ColumnName: columnName, ColumnType: columnType})
	}

	var count int64
	if err := db.Table(tableName).Count(&count).Error; err != nil {
		return nil, err
	}

	return &TableInfo{TableName: tableName, Columns: columns, RowCount: count}, nil
}

// PostDBConnectInfo 处理 POST 请求，连接数据库并返回表信息
func PostDBConnectInfo(c *gin.Context) {
	var dbInfo DBConnectInfo
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}

	db, err := ConnectToDB(dbInfo)
	if err != nil {
		util.Rsp(c, 500, "数据库连接失败: "+err.Error())
		return
	}

	tableInfos, err := getDatabaseTableNamesAndInfo(db)
	if err != nil {
		util.Rsp(c, 500, "获取数据库信息失败: "+err.Error())
		return
	}

	c.JSON(200, tableInfos)
}

// ConnectToDB 连接到数据库
func ConnectToDB(info DBConnectInfo) (*gorm.DB, error) {
	dsn := dsnFromDBConnectInfo(info)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}
