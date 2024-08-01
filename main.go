package main

import (
	"encoding/json"
	"fmt"
	"os"
	"text/template"
	"unicode"

	"template/sch"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var tableTemplate = `
	package model

	type {{ .Title }} struct{
		{{- range $k,$v := .Cols }}  
		{{ $v.Title }} {{ $v.Type }}` + " `json:" + `"` + "{{ $v.JsonName }}" + `" gorm:"column:{{ $v.GormName }}"` + ` form:"{{ $v.JsonName }}"` + ` {{- if $v.Binding}} binding:"{{ $v.Binding }}" {{- end}}` + "`" + `
		{{- end}}
	}

	func(m *{{ .Title }}) TableName ()string{
		return "{{ .TableName }}"
	}

`

var schemaTableTemplate = `
	package dto

	type {{ .Title }} struct{
		{{- range $k,$v := .Cols }}  
		{{ $v.Title }} {{ $v.Schema2GoType }}` + " `json:" + `"` + "{{ $v.ColumnName }}" + `" gorm:"{{ $v.GormFormat }}"` + ` form:"{{ $v.ColumnName }}"` + "`" + `
		{{- end}}
	}

	func(m *{{ .Title }}) TableName ()string{
		return "{{ .TableName }}"
	}

`

type Table struct {
	TableName string `json:"table_name"`
	Cols      []Col  `json:"cols"`
}

type Col struct {
	ColName  string `json:"col_name"`
	Type     string `json:"type"`
	JsonName string `json:"json_name"`
	GormName string `json:"gorm_name"`
	Binding  string `json:"binding"`
}

func main() {
	//read json config
	// tables := loadJsonConfig("./tables.json")
	// for _, tab := range tables.Tables {
	// 	createTemplate(tab)
	// }
	config := loadConfigJson("./tables.json")
	dsnSpf := fmt.Sprintf(`%s:%s@tcp(%s)/information_schema?charset=utf8mb4&parseTime=True&loc=Local`, config.MysqlUser, config.MysqlPwd, config.MysqlIP)
	CreateGormDb(dsnSpf)
	tables, err := sch.GetAllTable(DB, `cactus`)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, table := range tables {
		// fmt.Println(table.TableName)
		tcs := sch.TableColumnSchema{}
		tcs.TableName = table.TableName
		columns, err := sch.GetColumns(DB, table.TableName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tcs.Cols = columns
		for _, col := range columns {
			fmt.Println(col.String())
		}
		err = tcs.CreateTemplate(`./dto`, schemaTableTemplate)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

type tabel_infor struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}

func DescTable() {
	var Pricerecord []tabel_infor
	if err := DB.Raw(`desc m2_status`).Scan(&Pricerecord).Error; err != nil {
		fmt.Println(err)
	}

	for _, v := range Pricerecord {
		fmt.Println(v)
	}
}

func CreateGormDb(dsn string) {
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("open mysql failed,", err)
	}
	DB = d
}

type ConfigJson struct {
	MysqlIP        string `json:"mysql_ip"`
	MysqlPwd       string `json:"mysql_pwd"`
	MysqlUser      string `json:"mysql_user"`
	TargetDatabase string `json:"target_database"`
}

func loadConfigJson(path string) ConfigJson {
	readConfig, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("readConfig err %s", err)
		os.Exit(1)
	}
	config := ConfigJson{}
	if err := json.Unmarshal(readConfig, &config); err != nil {
		fmt.Printf("Unmarshal err %s", err)
		os.Exit(1)
	}
	return config
}

// ````````````````````````````````````````````````````````````````````````````
func createTemplate(table Table) {
	fileName := "./genfile/" + table.TableName + ".go"
	fd, err := os.Create(fileName)
	defer fd.Close()
	if err != nil {
		fmt.Printf("create file err %s", err)
		os.Exit(1)
	}
	template, err := template.New("tableTemplate").Parse(tableTemplate)
	if err != nil {
		fmt.Printf("template err %s", err)
		os.Exit(1)
	}
	if err := template.Execute(fd, table); err != nil {
		fmt.Printf("template err %s", err)
		os.Exit(1)
	}
}

func (t Table) Title() string {
	s := t.TableName
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
func (t Col) Title() string {
	s := t.ColName
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

type Tables struct {
	Tables []Table `json:"tables"`
}

func loadJsonConfig(path string) Tables {
	readConfig, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("readConfig err %s", err)
		os.Exit(1)
	}
	tableItem := make([]Table, 0)
	tables := Tables{
		Tables: tableItem,
	}
	if err := json.Unmarshal(readConfig, &tables); err != nil {
		fmt.Printf("Unmarshal err %s", err)
		os.Exit(1)
	}
	return tables
}
