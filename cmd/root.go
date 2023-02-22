package console

import (
	"github.com/purerun/dbhelper/cmd/codemake"
	"github.com/purerun/dbhelper/util"
	"github.com/purerun/dbhelper/util/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dbhelper",
	Short: "数据库助手/database helper",
	Long:  `数据库助手/database helper`,
	PersistentPreRun: func(command *cobra.Command, args []string) {
		if !util.IsExist("./.env") {
			err := util.Put([]byte(app.GetEnvExample()), "./.env")
			if err != nil {
				panic(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(codemake.GetCommands()...)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
