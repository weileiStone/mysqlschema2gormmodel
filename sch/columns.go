package sch

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Column struct {
	TableName   string `json:"TABLE_NAME" gorm:"column:TABLE_NAME"`
	ColumnName  string `json:"COLUMN_NAME" gorm:"column:COLUMN_NAME"`
	DataType    string `json:"DATA_TYPE" gorm:"column:DATA_TYPE"`
	Default     string `json:"COLUMN_DEFAULT" gorm:"column:COLUMN_DEFAULT"`
	COLUMN_KEY  string `json:"COLUMN_KEY" gorm:"column:COLUMN_KEY"`
	COLUMN_TYPE string `json:"COLUMN_TYPE" gorm:"column:COLUMN_TYPE"`
	IS_NULLABLE string `json:"IS_NULLABLE" gorm:"column:IS_NULLABLE"`
}

func GetColumns(db *gorm.DB, table string) ([]Column, error) {
	columns := make([]Column, 0)
	if err := db.Where(`TABLE_NAME = ?`, table).Table(`COLUMNS`).Find(&columns).Error; err != nil {
		return columns, err
	}
	return columns, nil
}

func (c *Column) String() string {
	return fmt.Sprintf(`table name %s, columns %s, data type %s default %s, column key %s, is null %s `, c.TableName, c.ColumnName, c.DataType, c.Default, c.COLUMN_KEY, c.IS_NULLABLE)
}

var SchemaGoTypeMap = map[string]string{
	`int`:        `int`,
	`varchar`:    `string`,
	`char`:       `string`,
	`mediuetext`: `string`,
	`longtext`:   `string`,
	`datetime`:   `time.Time`,
	`text`:       `string`,
	`smallint`:   `int16`,
	`tinyint`:    `int8`,
	`mediumint`:  `int32`,
	`bigint`:     `int64`,
}

func (c Column) Schema2GoType() string {
	changeType := `string`
	if t, ok := SchemaGoTypeMap[c.DataType]; ok {
		return t
	}
	return changeType
}

func (c Column) GormFormat() string {
	gormStringSlice := make([]string, 0)
	columnStr := fmt.Sprintf(`column:%s;`, c.ColumnName)
	gormStringSlice = append(gormStringSlice, columnStr)
	typeStr := fmt.Sprintf(`type:%s;`, c.COLUMN_TYPE)
	gormStringSlice = append(gormStringSlice, typeStr)
	if c.Default != `` {
		defaultStr := fmt.Sprintf(`default:%s;`, c.Default)
		gormStringSlice = append(gormStringSlice, defaultStr)
	}
	if c.COLUMN_KEY == `PRI` {
		gormStringSlice = append(gormStringSlice, `primarykey;`)
	}
	if c.IS_NULLABLE == `NO` { //gorm最后一个
		gormStringSlice = append(gormStringSlice, `NOT NULL;`)
	}
	return strings.Join(gormStringSlice, ``)
}
