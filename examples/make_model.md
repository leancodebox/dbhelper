## 自动生成 gorm 文件

### 创建如下数据表

例如在 `test` 数据库中创建:

```sql
CREATE TABLE `users`
(
    `id`            bigint unsigned AUTO_INCREMENT,
    `created_at`    datetime(3) NULL,
    `updated_at`    datetime(3) NULL,
    `deleted_at`    datetime(3) NULL,
    `name`          longtext,
    `email`         longtext,
    `age`           tinyint unsigned,
    `birthday`      datetime( 3) NULL,
    `member_number` longtext,
    `activated_at`  datetime(3) NULL,
    PRIMARY KEY (`id`),
    INDEX           idx_users_deleted_at (`deleted_at`)
)
```

### 初始化配置文件

```shell
$ dbhelper init

# 配置文件初始化完成，你可以查看当前目录下 .env 文件
```

### 修改配置文件

```shell
$ vim .env

# 你的数据库配置
ORIGIN_DATABASE_URL=root:123456@tcp(:3306)/test?charset=utf8mb4&parseTime=True&loc=Local
# 生成代码目录
MODEL_OUTPUT_DIR=./model/
```

### 生成 gorm 文件

```shell
$ dbhelper make:model
```

### 查看 gorm 文件列表

```shell
ls -l model/Users

total 12
-rw-r--r-- 1 stars stars 1969 Mar  2 13:58 Users.go
-rw-r--r-- 1 stars stars  316 Mar  2 13:58 Users_connect.go
-rw-r--r-- 1 stars stars  823 Mar  2 13:58 Users_rep.go
```

### 查看 gorm 文件内容

**Users.go**

```go
package Users

import (
	"time"
)

const tableName = "users"
const pid = "id"
const fieldCreatedAt = "created_at"
const fieldUpdatedAt = "updated_at"
const fieldDeletedAt = "deleted_at"
const fieldName = "name"
const fieldEmail = "email"
const fieldAge = "age"
const fieldBirthday = "birthday"
const fieldMemberNumber = "member_number"
const fieldActivatedAt = "activated_at"

type Users struct {
	Id           uint64     `gorm:"primaryKey;column:id;autoIncrement;not null;" json:"id"`   //
	CreatedAt    *time.Time `gorm:"column:created_at;type:datetime(3);" json:"createdAt"`     //
	UpdatedAt    *time.Time `gorm:"column:updated_at;type:datetime(3);" json:"updatedAt"`     //
	DeletedAt    *time.Time `gorm:"column:deleted_at;type:datetime(3);" json:"deletedAt"`     //
	Name         string     `gorm:"column:name;type:longtext;" json:"name"`                   //
	Email        string     `gorm:"column:email;type:longtext;" json:"email"`                 //
	Age          uint8      `gorm:"column:age;type:tinyint(3) unsigned;" json:"age"`          //
	Birthday     *time.Time `gorm:"column:birthday;type:datetime(3);" json:"birthday"`        //
	MemberNumber string     `gorm:"column:member_number;type:longtext;" json:"memberNumber"`  //
	ActivatedAt  *time.Time `gorm:"column:activated_at;type:datetime(3);" json:"activatedAt"` //
}

// func (itself *Users) BeforeSave(tx *gorm.DB) (err error) {}
// func (itself *Users) BeforeCreate(tx *gorm.DB) (err error) {}
// func (itself *Users) AfterCreate(tx *gorm.DB) (err error) {}
// func (itself *Users) BeforeUpdate(tx *gorm.DB) (err error) {}
// func (itself *Users) AfterUpdate(tx *gorm.DB) (err error) {}
// func (itself *Users) AfterSave(tx *gorm.DB) (err error) {}
// func (itself *Users) BeforeDelete(tx *gorm.DB) (err error) {}
// func (itself *Users) AfterDelete(tx *gorm.DB) (err error) {}
// func (itself *Users) AfterFind(tx *gorm.DB) (err error) {}

func (Users) TableName() string {
	return tableName
}
```

**Users_connect.go**

```go
package Users

import (
	"gorm.io/gorm"

	db "thh/conf/dbconnect"
)

// Prohibit manual changes
// 禁止手动更改本文件

func builder() *gorm.DB {
	return db.Std().Table(tableName)
}

func First(db *gorm.DB) (el Users) {
	db.First(&el)
	return
}

func List(db *gorm.DB) (el []Users) {
	db.Find(&el)
	return
}
```

**Users_rep.go**

```go
package Users

func Create(entity *Users) int64 {
	result := builder().Create(entity)
	return result.RowsAffected
}

func Save(entity *Users) int64 {
	result := builder().Save(entity)
	return result.RowsAffected
}

func SaveAll(entities *[]Users) int64 {
	result := builder().Save(entities)
	return result.RowsAffected
}

func Delete(entity *Users) int64 {
	result := builder().Delete(entity)
	return result.RowsAffected
}

func Get(id any) (entity *Users) {
	builder().Where(pid, id).First(entity)
	return
}

func GetBy(field, value string) (entity Users) {
	builder().Where(field+" = ?", value).First(&entity)
	return
}

func All() (entities []Users) {
	builder().Find(&entities)
	return
}

func IsExist(field, value string) bool {
	var count int64
	builder().Where(field+" = ?", value).Count(&count)
	return count > 0
}
```