// Package gen 命令行的 gen 命令
package codemake

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/leancodebox/dbhelper/util/config"
	"github.com/spf13/cast"
	"os"
	"path/filepath"
	"text/template"

	"github.com/leancodebox/dbhelper/util"
	"github.com/leancodebox/dbhelper/util/stropt"

	"github.com/iancoleman/strcase"
)

type Model struct {
	StructName  string
	PackageName string
	ClientName  string
}

//go:embed tmpl
var tmplFS embed.FS

// makeModelFromString 格式化用户输入的内容
func makeModelFromString(name string) Model {
	model := Model{}
	model.StructName = stropt.Singular(strcase.ToCamel(name))
	model.PackageName = stropt.Snake(model.StructName)
	model.ClientName = stropt.Camel(model.StructName)
	return model
}

func buildWithOutput(data map[string]any, filePath string, tmplPath string) {
	outputData := buildByTmpl(data, tmplPath)
	dirPath := filepath.Dir(filePath)
	if !util.IsExist(filePath) {
		if err := os.MkdirAll(dirPath, 0666); err != nil {
		}
	}
	err := util.Put([]byte(outputData), filePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 提示成功
	fmt.Println(fmt.Sprintf("[%s] created.", filePath))
}

func buildByTmpl(data map[string]any, tmplPath string) string {

	modelData, err := tmplFS.ReadFile(tmplPath)
	if err != nil {
		fmt.Println("err", err)
	}
	modelStub := string(modelData)

	var b bytes.Buffer
	t := template.New("test")
	t.Funcs(template.FuncMap{
		"AddOne": func(p int) int { return p + 1 },
	})
	t = template.Must(t.Parse(modelStub))

	err = t.Execute(&b, data)

	if err != nil {
		fmt.Println(err)
	}
	return b.String()
}

func runConfig(action func(targetUrl, originUrl, dbConnect, output string)) {
	oConfig := config.GetAny("db")
	if oConfig == nil {
		fmt.Println("当前尚未有配置")
		return
	}
	if configList, ok := oConfig.([]any); ok {
		for index, itemConfig := range configList {
			if itemConfigMap, ok := itemConfig.(map[string]any); ok {
				targetUrl := cast.ToString(itemConfigMap["target_url"])
				originUrl := cast.ToString(itemConfigMap["origin_url"])
				dbConnect := cast.ToString(itemConfigMap["connect"])
				output := cast.ToString(itemConfigMap["output"])
				action(targetUrl, originUrl, dbConnect, output)
			} else {
				fmt.Println("db 配置有误，分组：", index)
			}
		}
	} else {
		fmt.Println("配置有误请检查")
	}
}
