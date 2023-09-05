package db

import (
	databasesql "database/sql"
	"douyin/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var orm *gorm.DB
var sql *databasesql.DB

func ORM() *gorm.DB {
	return orm
}

func SQL() *databasesql.DB {
	return sql
}

func Init(mysqlConfig *config.MySQL) error {
	var err error

	if err = CreateDatabase(mysqlConfig); err != nil {
		return err
	}
	if orm, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       getDSN(mysqlConfig), // DSN data source name
		DefaultStringSize:         256,                 // string 类型字段的默认长度
		DisableDatetimePrecision:  false,               // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,               // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt: true,
	}); err != nil {
		return fmt.Errorf("连接数据库失败：%w", err)
	}
	if sql, err = orm.DB(); err != nil {
		return fmt.Errorf("获取数据库 *sql.DB 对象失败：%w", err)
	}
	log.Println("数据库连接成功")

	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='用户表'",
	).AutoMigrate(&User{}); err != nil {
		return fmt.Errorf("初始化用户表失败：%w", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='用户登录状态表'",
	).AutoMigrate(&UserToken{}); err != nil {
		return fmt.Errorf("初始化用户登录状态表失败：%w", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='视频表'",
	).AutoMigrate(&Video{}); err != nil {
		return fmt.Errorf("初始化视频表失败：%w", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='消息表'",
	).AutoMigrate(&Message{}); err != nil {
		return fmt.Errorf("初始化消息表失败：%w", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='关注关系表'",
	).AutoMigrate(&Follow{}); err != nil {
		return fmt.Errorf("初始化关注关系表失败：%w", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='点赞记录表'",
	).AutoMigrate(&Favorite{}); err != nil {
		return fmt.Errorf("初始化点赞记录表失败：%w", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='评论表'",
	).AutoMigrate(&Comment{}); err != nil {
		return fmt.Errorf("初始化评论表失败：%w", err)
	}

	return nil
}

func getDSN(mysqlConfig *config.MySQL) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Dbname,
		mysqlConfig.Options,
	)
}

func getDSNWithoutDBName(mysqlConfig *config.MySQL) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/?%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Options,
	)
}

func CreateDatabase(mysqlConfig *config.MySQL) error {
	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;",
		mysqlConfig.Dbname,
	)

	tmpSQL, err := databasesql.Open("mysql", getDSNWithoutDBName(mysqlConfig))
	if nil != err {
		return fmt.Errorf("连接临时数据库失败：%w", err)
	}
	defer func() {
		_ = tmpSQL.Close()
	}()

	if result, err := tmpSQL.Exec(query); nil != err {
		return fmt.Errorf("创建数据库失败：%w", err)
	} else if affected, err := result.RowsAffected(); nil != err || affected == 0 {
		return fmt.Errorf("创建数据库失败：%w", err)
	}

	return nil
}

func DropDatabase(mysqlConfig *config.MySQL) error {
	query := fmt.Sprintf(
		"DROP DATABASE IF EXISTS %s;",
		mysqlConfig.Dbname,
	)

	tmpSQL, err := databasesql.Open("mysql", getDSNWithoutDBName(mysqlConfig))
	if nil != err {
		return fmt.Errorf("连接临时数据库失败：%w", err)
	}
	defer func() {
		_ = tmpSQL.Close()
	}()

	if result, err := tmpSQL.Exec(query); nil != err {
		return fmt.Errorf("删除数据库失败：%w", err)
	} else if affected, err := result.RowsAffected(); nil != err || affected == 0 {
		return fmt.Errorf("删除数据库失败：%w", err)
	}

	return nil
}
