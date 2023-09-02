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

func init() {
	var err error

	createDatabase(config.Conf.MySQL.Dbname)
	if orm, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       getDSN(config.Conf.MySQL.Dbname), // DSN data source name
		DefaultStringSize:         256,                              // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                             // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                             // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                             // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                            // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt: true,
	}); err != nil {
		log.Panicln("数据库连接失败", err)
	}
	if sql, err = orm.DB(); err != nil {
		log.Panicln("获取数据库 *sql.DB 对象失败", err)
	}
	log.Println("数据库连接成功")

	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='用户表'",
	).AutoMigrate(&User{}); err != nil {
		log.Panicln("初始化用户表失败", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='用户登录状态表'",
	).AutoMigrate(&UserToken{}); err != nil {
		log.Panicln("初始化用户登录状态表失败", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='视频表'",
	).AutoMigrate(&Video{}); err != nil {
		log.Panicln("初始化视频表失败", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='消息表'",
	).AutoMigrate(&Message{}); err != nil {
		log.Panicln("初始化消息表失败", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='关注关系表'",
	).AutoMigrate(&Follow{}); err != nil {
		log.Panicln("初始化关注关系表失败", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='点赞记录表'",
	).AutoMigrate(&Favorite{}); err != nil {
		log.Panicln("初始化点赞记录表失败", err)
	}
	if err = orm.Set(
		"gorm:table_options",
		" COMMENT='评论表'",
	).AutoMigrate(&Comment{}); err != nil {
		log.Panicln("初始化评论表失败", err)
	}
}

func getDSN(dbname string) string {
	mysqlConfig := config.Conf.MySQL
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		dbname,
		mysqlConfig.Options,
	)
}

func createDatabase(dbname string) {
	query := fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_general_ci;",
		dbname,
	)

	tmpSQL, err := databasesql.Open("mysql", getDSN(""))
	if nil != err {
		log.Panicln("连接临时数据库失败", err)
	}
	defer func() {
		_ = tmpSQL.Close()
	}()

	if result, err := tmpSQL.Exec(query); nil != err {
		log.Panicln("创建数据库失败", err)
	} else if affected, err := result.RowsAffected(); nil != err || affected == 0 {
		log.Panicln("创建数据库失败", err)
	}
}
