package codemake

import (
	"fmt"
	"github.com/purerun/dbhelper/util/app"
	"github.com/purerun/dbhelper/util/config"
	"github.com/purerun/dbhelper/util/eh"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

func init() {
	appendCommand(&cobra.Command{
		Use:   "migrate:tabledata",
		Short: "导表助手（表数据）",
		Long:  ".env 文件中 TARGET_DATABASE_URL(local) 为导入db地址， ORIGIN_DATABASE_URL(tmp) 为被导入db地址。会获取tmp中的表信息导入至local中",
		Run:   runMTableDataFromDb,
		//Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	})
}

func runMTableDataFromDb(_ *cobra.Command, _ []string) {
	// init
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,         // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,        // 禁用彩色打印
		},
	)
	dataSourceName := config.GetString("ORIGIN_DATABASE_URL")
	localSourceName := config.GetString("TARGET_DATABASE_URL")
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

	limit := 600
	//根据规则提取关键信息
	for tmpTableName, _ := range tbDesc {
		var list []genColumns
		// Get table annotations.获取表注释
		db.Raw("show FULL COLUMNS from " + tmpTableName).Scan(&list)
		primaryKey := ""
		for _, column := range list {
			if column.Key == "PRI" {
				primaryKey = column.Field
			}
		}
		lastId := 0
		offset := 0
		for {
			selectSql := ""
			if primaryKey == "" {
				selectSql = fmt.Sprintf("select * from %v  LIMIT %v,%v",
					tmpTableName, offset, limit,
				)
			} else {
				selectSql = fmt.Sprintf("select * from %v where %v > %v order by %v LIMIT %v",
					tmpTableName, primaryKey, lastId, primaryKey, limit,
				)
			}
			var dataList []map[string]any
			db.Raw(selectSql).Find(&dataList)
			fmt.Println(tmpTableName + ":" + cast.ToString(len(dataList)))
			if len(dataList) > 1 {
				itemData := dataList[len(dataList)-1]
				lastId = cast.ToInt(itemData[primaryKey])
				offset += limit
				cErr := localDb.Table(tmpTableName).Clauses(clause.Insert{Modifier: "IGNORE"}).Create(dataList).Error
				if cErr != nil {
					fmt.Println(cErr)
				}
			}
			if len(dataList) < limit {
				break
			}

		}
	}
	fmt.Println("迁移完毕")
	fmt.Println(app.GetUnitTime())

}
