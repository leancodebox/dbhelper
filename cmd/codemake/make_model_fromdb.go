package codemake

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/leancodebox/dbhelper/util"
	"github.com/leancodebox/dbhelper/util/eh"
	"github.com/leancodebox/dbhelper/util/stropt"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func init() {
	appendCommand(GModel)
}

var GModel = &cobra.Command{
	Use:   "make:model",
	Short: "从db创建gorm",
	Run:   runGModel,
	//Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
}

type modelEm struct {
	Name  string // 字段名
	Notes string // 注释
	Type  string // 字段类型
}
type genColumns struct {
	Field   string  `gorm:"column:Field"`
	Type    string  `gorm:"column:Type"`
	Key     string  `gorm:"column:Key"`
	Desc    string  `gorm:"column:Comment"`
	Null    string  `gorm:"column:Null"`
	Default *string `gorm:"column:Default"`
}

func runGModel(_ *cobra.Command, _ []string) {

	// init
	//dataSourceName := config.GetString("dbTool.originUrl")
	//dbStd := fmt.Sprintf(`"%v"`, config.GetString("dbTool.dbConnect", `thh/conf/dbconnect`))
	//outputRoot := config.GetString("dbTool.output", "./storage/tmp/model/")

	runConfig(func(targetUrl, originUrl, dbConnect, output string) {
		makeModel(originUrl, dbConnect, output)
	})

}

func makeModel(dataSourceName, dbStd, outputRoot string) {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{PrepareStmt: false,
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, // 全局禁用表名复数
		Logger:         newLogger})
	if eh.PrIF(err) {
		return
	}

	rows, err := db.Raw("show tables").Rows()
	tbDesc := make(map[string]string)
	if eh.PrIF(err) {
		return
	}
	for rows.Next() {
		var table string
		eh.PrIF(rows.Scan(&table))
		tbDesc[table] = table
	}
	eh.PrIF(rows.Close())

	for tmpTableName, _ := range tbDesc {
		fmt.Println("start gen ", tmpTableName)
		var list []genColumns

		// Get table annotations.获取表注释
		db.Raw("show FULL COLUMNS from " + tmpTableName).Scan(&list)

		modelStr, connectStr, repStr := buildModelContent(tmpTableName, list, dbStd)
		modelPath := stropt.LowerCamel(stropt.Camel(tmpTableName))

		modelEntityPath := outputRoot + modelPath + "/" + modelPath + ".go"
		connectPath := outputRoot + modelPath + "/" + modelPath + "_connect.go"
		repPath := outputRoot + modelPath + "/" + modelPath + "_rep.go"

		util.PutContent(modelEntityPath, modelStr)
		util.PutContent(connectPath, connectStr)

		if !util.IsExist(repPath) {
			util.PutContent(repPath, repStr)
		}
	}
	o, e := exec.Command("gofmt", "-w", outputRoot).Output()
	if e != nil {
		fmt.Println("fmt error", e)
	} else {
		fmt.Println("fmt success", string(o))
	}
	fmt.Println("build end")
}

type Field struct {
	Name            string
	DBFieldName     string
	FieldName       string
	Type            string
	Index           string
	StructFieldName string
	JsonName        string
	Pid             bool
	NullTag         string
	DefaultTag      string
	TypeTag         string
	Comment         string
	Field           string
}

func buildModelContent(tmpTableName string, list []genColumns, dbStd string) (string, string, string) {

	importList := map[string]string{}
	var hasPid = false
	var pidFiledName = ""
	var fieldList []Field

	for _, value := range list {
		var field string
		if IsNum(string(value.Field[0])) {
			field = "Column" + stropt.Camel(value.Field)
		} else {
			field = stropt.Camel(value.Field)
		}
		if pkgName, ok := EImportsHead[getTypeName(value.Type, false)]; ok {
			importList[pkgName] = pkgName
		}
		nullTag := ""
		if value.Null == "NO" {
			nullTag = "not null;"
		}

		defaultStr := ""
		if value.Default != nil {
			defaultStr = "default:"
			if len(*value.Default) == 0 {
				defaultStr += "''"
			} else {
				defaultStr += *value.Default
			}
			defaultStr += ";"
		}

		fieldName := "field" + stropt.Camel(stropt.LowerCamel(value.Field))
		typeString := `type:` + value.Type
		if value.Key == "PRI" {
			typeString = `autoIncrement`
			fieldName = `pid`
			hasPid = true
			pidFiledName = field
		}
		typeString += ";"

		fieldList = append(fieldList, Field{
			Name:            value.Field,
			DBFieldName:     value.Field,
			Field:           value.Field,
			FieldName:       fieldName,
			StructFieldName: field,
			JsonName:        stropt.LowerCamel(value.Field),
			Index:           value.Key,
			Type:            getTypeName(value.Type, value.Null != "NO"),
			Pid:             value.Key == "PRI",
			NullTag:         nullTag,
			DefaultTag:      defaultStr,
			TypeTag:         typeString,
			Comment:         strings.ReplaceAll(value.Desc, "\n", ""),
		})

	}
	if IsNum(string(tmpTableName[0])) {
		tmpTableName = "M" + tmpTableName
	}
	modelStr := buildByTmpl(
		map[string]any{
			"TableName":  tmpTableName,
			"pkgName":    stropt.LowerCamel(tmpTableName),
			"ModelName":  "Entity", //str.Camel(tmpTableName),
			"importList": importList,
			"fieldList":  fieldList,
		},
		"tmpl/db/entity.tmpl",
	)
	connectStr := buildByTmpl(
		map[string]any{
			"TableName":  tmpTableName,
			"pkgName":    stropt.LowerCamel(tmpTableName),
			"ModelName":  "Entity", //str.Camel(tmpTableName),
			"importList": importList,
			"fieldList":  fieldList,
			"DBPkg":      "\"" + dbStd + "\"",
		},
		"tmpl/db/connect.tmpl",
	)
	repStr := buildByTmpl(
		map[string]any{
			"TableName":    tmpTableName,
			"pkgName":      stropt.LowerCamel(tmpTableName),
			"ModelName":    "Entity", //str.Camel(tmpTableName),
			"importList":   importList,
			"fieldList":    fieldList,
			"hasPid":       hasPid,
			"pidFiledName": pidFiledName,
		},
		"tmpl/db/rep.tmpl",
	)
	return modelStr, connectStr, repStr
}

func getFileLine(path string) int {
	file, _ := os.Open(path)
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	lineCount := 0
	for fileScanner.Scan() {
		lineCount++
	}
	return lineCount
}

// IsNum 判断是否是数字 用来处理数字开头的字段
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// getTypeName Type acquisition filtering.类型获取过滤
func getTypeName(name string, isNull bool) string {
	// 优先匹配自定义类型

	// Precise matching first.先精确匹配
	if v, ok := TypeMysqlDicMp[name]; ok {
		return fixNullToPointer(v, isNull)
	}

	// Fuzzy Regular Matching.模糊正则匹配
	for _, l := range TypeMysqlMatchList {
		if ok, _ := regexp.MatchString(l.Key, name); ok {
			return fixNullToPointer(l.Value, isNull)
		}
	}

	panic(fmt.Sprintf("type (%v) not match in any way.maybe need to add on (https://github.com/xxjwxc/gormt/blob/master/data/view/cnf/def.go)", name))
}

// fixNullToPointer 修复为指针 目前针对所在项目，仅time，设置可设置为nil
func fixNullToPointer(name string, isNull bool) string {
	if isNull {
		//if strings.HasPrefix(name, "uint") {
		//	return "*" + name
		//}
		//if strings.HasPrefix(name, "int") {
		//	return "*" + name
		//}
		//if strings.HasPrefix(name, "float") {
		//	return "*" + name
		//}
		//if strings.HasPrefix(name, "date") {
		//	return "*" + name
		//}
		if strings.HasPrefix(name, "time") {
			return "*" + name
		}
		//if strings.HasPrefix(name, "bool") {
		//	return "*" + name
		//}
		//if strings.HasPrefix(name, "string") {
		//	return "*" + name
		//}
	}

	return name
}

var EImportsHead = map[string]string{
	"stirng":         `"string"`,
	"time.Time":      `"time"`,
	"gorm.Model":     `"gorm.io/gorm"`,
	"fmt":            `"fmt"`,
	"datatypes.JSON": `"gorm.io/datatypes"`,
	"datatypes.Date": `"gorm.io/datatypes"`,
}

var TypeMysqlDicMp = map[string]string{
	"smallint":            "int16",
	"smallint unsigned":   "uint16",
	"int":                 "int",
	"int unsigned":        "uint",
	"bigint":              "int64",
	"bigint unsigned":     "uint64",
	"mediumint":           "int32",
	"mediumint unsigned":  "uint32",
	"varchar":             "string",
	"char":                "string",
	"date":                "time.Time",
	"datetime":            "time.Time",
	"bit(1)":              "[]uint8",
	"tinyint":             "int8",
	"tinyint unsigned":    "uint8",
	"tinyint(1)":          "int", // tinyint(1) 默认设置成bool
	"tinyint(1) unsigned": "int", // tinyint(1) 默认设置成bool
	"json":                "string",
	"text":                "string",
	"timestamp":           "time.Time",
	"double":              "float64",
	"double unsigned":     "float64",
	"mediumtext":          "string",
	"longtext":            "string",
	"float":               "float32",
	"float unsigned":      "float32",
	"tinytext":            "string",
	"enum":                "string",
	"time":                "time.Time",
	"tinyblob":            "[]byte",
	"blob":                "[]byte",
	"mediumblob":          "[]byte",
	"longblob":            "[]byte",
	"integer":             "int64",
	"numeric":             "float64",
	"smalldatetime":       "time.Time", //sqlserver
	"nvarchar":            "string",
	"real":                "float32",
	"binary":              "[]byte",
}

var TypeMysqlMatchList = []struct {
	Key   string
	Value string
}{
	{`^(tinyint)[(]\d+[)] unsigned`, "uint8"},
	{`^(smallint)[(]\d+[)] unsigned`, "uint16"},
	{`^(int)[(]\d+[)] unsigned`, "uint32"},
	{`^(bigint)[(]\d+[)] unsigned`, "uint64"},
	{`^(float)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(double)[(]\d+,\d+[)] unsigned`, "float64"},
	{`^(tinyint)[(]\d+[)]`, "int8"},
	{`^(smallint)[(]\d+[)]`, "int16"},
	{`^(int)[(]\d+[)]`, "int"},
	{`^(bigint)[(]\d+[)]`, "int64"},
	{`^(char)[(]\d+[)]`, "string"},
	{`^(enum)[(](.)+[)]`, "string"},
	{`^(varchar)[(]\d+[)]`, "string"},
	{`^(varbinary)[(]\d+[)]`, "[]byte"},
	{`^(blob)[(]\d+[)]`, "[]byte"},
	{`^(binary)[(]\d+[)]`, "[]byte"},
	{`^(decimal)[(]\d+,\d+[)]`, "float64"},
	{`^(mediumint)[(]\d+[)]`, "int16"},
	{`^(mediumint)[(]\d+[)] unsigned`, "uint16"},
	{`^(double)[(]\d+,\d+[)]`, "float64"},
	{`^(float)[(]\d+,\d+[)]`, "float64"},
	{`^(datetime)[(]\d+[)]`, "time.Time"},
	{`^(bit)[(]\d+[)]`, "[]uint8"},
	{`^(text)[(]\d+[)]`, "string"},
	{`^(integer)[(]\d+[)]`, "int"},
	{`^(timestamp)[(]\d+[)]`, "time.Time"},
	{`^(geometry)[(]\d+[)]`, "[]byte"},
	{`^(set)[(][\s\S]+[)]`, "string"},
	{`^(point)`, "[]byte"},
}
