package codemake

import (
	"fmt"

	"github.com/leancodebox/dbhelper/util"
	"github.com/leancodebox/dbhelper/util/config"
	"github.com/leancodebox/dbhelper/util/jsonopt"
	"github.com/leancodebox/dbhelper/util/stropt"
	"github.com/spf13/cobra"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "make:modelfromjson",
		Short: "从db创建gorm",
		Run:   runGModelFromJson,
		//Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

type ModelMakeItem struct {
	ModelName string       `json:"modelName"`
	KeyList   []genColumns `json:"keyList"`
}

func runGModelFromJson(_ *cobra.Command, _ []string) {
	data, _ := util.FileGetContents("model.json")
	genColumnsList := jsonopt.Decode[[]ModelMakeItem](data)
	for _, item := range genColumnsList {
		outPutModel(item.ModelName, item.KeyList)
	}
}

func outPutModel(modelName string, list []genColumns) {
	outputRoot := config.GetString("default.output", "./storage/model/")
	connect := config.GetString("default.connect", "connect")
	outputRoot = `./storage/model/`
	modelPath := stropt.LowerCamel(modelName)
	modelEntityPath := outputRoot + modelPath + "/" + modelPath + ".go"
	connectPath := outputRoot + modelPath + "/" + modelPath + "_connect.go"
	repPath := outputRoot + modelPath + "/" + modelPath + "_rep.go"

	modelStr, repStr := buildModelContent(stropt.Snake(modelName), list, connect)

	fmt.Println(modelStr, repStr)
	fmt.Println(modelEntityPath)
	fmt.Println(connectPath)
	fmt.Println(repPath)
	util.FilePutContents(modelEntityPath, modelStr)
	util.IsExistOrCreate(repPath, repStr)
}
