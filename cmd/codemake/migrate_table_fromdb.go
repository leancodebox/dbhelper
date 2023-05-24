package codemake

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/leancodebox/dbhelper/util/config"
	"github.com/leancodebox/dbhelper/util/eh"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "migrate:table",
		Short: "导表助手（表结构）",
		Long:  ".env 文件中 TARGET_DATABASE_URL(local) 为导入db地址， ORIGIN_DATABASE_URL(tmp) 为被导入db地址。会获取tmp中的表信息导入至local中",
		Run:   runMTableFromDb,
		//Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runMTableFromDb(_ *cobra.Command, _ []string) {
	// init
	dataSourceName := config.GetString("dbTool.originUrl")
	localSourceName := config.GetString("dbTool.targetUrl")
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	)
	localDb, err := gorm.Open(mysql.Open(localSourceName), &gorm.Config{PrepareStmt: false,
		NamingStrategy: schema.NamingStrategy{SingularTable: true}, // 全局禁用表名复数
		Logger:         newLogger})
	if eh.PrIF(err) {
		return
	}
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

	reg, _ := regexp.Compile("AUTO_INCREMENT=[0-9]+")
	// (datetime[A-Za-z0-1\s]+')(0000-00-00 00:00:00)'
	//根据规则提取关键信息
	for tmpTableName, _ := range tbDesc {
		var list []ShowTable

		// Get table annotations.获取表注释
		sctSql := "show create table " + tmpTableName
		db.Raw(sctSql).Scan(&list)
		if len(list) == 1 {
			//list[0].CreateTable
			createSql := reg.ReplaceAllString(list[0].CreateTable, `AUTO_INCREMENT=1`)
			createSql = strings.ReplaceAll(createSql, "0000-00-00 00:00:00", "1970-01-01 00:00:01")
			createSql = strings.ReplaceAll(createSql, "0000-00-00", "1970-01-01")
			cErr := localDb.Exec(createSql).Error
			if cErr != nil {
				fmt.Println()
			} else {
				fmt.Println(tmpTableName + " 迁移完毕")
			}
		} else {
			fmt.Println("存在异常查询", sctSql)
		}
	}

	fmt.Println("迁移完毕")

}

type ShowTable struct {
	Table       string `gorm:"column:Table"`
	CreateTable string `gorm:"column:Create Table"`
}
