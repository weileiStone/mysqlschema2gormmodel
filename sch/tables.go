package sch

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
	"unicode"

	"gorm.io/gorm"
)

type Table struct {
	TableName string `json:"TABLE_NAME" gorm:"column:TABLE_NAME"`
}

type TableColumnSchema struct {
	TableName string   `json:"TABLE_NAME" gorm:"column:TABLE_NAME"`
	Cols      []Column `json:"columns"`
}

func GetAllTable(db *gorm.DB, schema string) ([]Table, error) {
	tables := make([]Table, 0)
	if err := db.Where(`TABLE_SCHEMA = ?`, schema).Table(`TABLES`).Select(`TABLE_NAME`).Find(&tables).Error; err != nil {
		return tables, err
	}
	return tables, nil
}

func (tcs *TableColumnSchema) String() string {
	return fmt.Sprintf(`table name %s, colums length %d`, tcs.TableName, len(tcs.Cols))
}

func (tcs TableColumnSchema) Title() string {
	reg := regexp.MustCompile(`^\d+`)
	s := reg.ReplaceAllString(tcs.TableName, "")
	if s == "" {
		return s
	}
	sslice := strings.Split(s, `_`)
	newss := make([]string, len(sslice))
	for _, news := range sslice {
		newRunes := []rune(news)
		newRunes[0] = unicode.ToUpper(newRunes[0])
		newss = append(newss, string(newRunes))
	}
	return strings.Join(newss, ``)
}
func (c *Column) Title() string {
	reg := regexp.MustCompile(`^\d+`)
	s := reg.ReplaceAllString(c.ColumnName, "")
	if s == "" {
		return s
	}
	sslice := strings.Split(s, `_`)
	newss := make([]string, len(sslice))
	for _, news := range sslice {
		newRunes := []rune(news)
		newRunes[0] = unicode.ToUpper(newRunes[0])
		newss = append(newss, string(newRunes))
	}
	return strings.Join(newss, ``)
}

func (tcs *TableColumnSchema) CreateTemplate(topPath string, templateStr string) error {
	fileName := topPath + `/` + tcs.TableName + `.go`
	fd, err := os.Create(fileName)
	defer fd.Close()
	if err != nil {
		fmt.Printf("create file err %s", err)
	}
	template, err := template.New("tableTemplateTcs").Parse(templateStr)
	if err != nil {
		return err
	}
	if err := template.Execute(fd, *tcs); err != nil {
		return err
	}

	return nil
}
