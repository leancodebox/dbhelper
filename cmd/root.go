package console

import (
	"fmt"
	"github.com/leancodebox/dbhelper/cmd/codemake"
	"github.com/leancodebox/dbhelper/util"
	"github.com/leancodebox/dbhelper/util/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbhelper",
	Short: "数据库助手/database helper",
	Long:  `数据库助手/database helper`,
	PersistentPreRun: func(command *cobra.Command, args []string) {
		if !util.IsExist("./config.toml") {
			err := util.Put(app.GetDefaultConfig(), "./config.toml")
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "初始化配置文件",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("配置文件初始化完成，你可以查看当前目录下 config.toml 文件")
		},
	})
	rootCmd.AddCommand(codemake.GetCommands()...)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
