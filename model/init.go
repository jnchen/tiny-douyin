package model

import (
	"douyin/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDatabase(config *config.MySQLConfig) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", config.Username, config.Password, config.Host, config.Port, config.Dbname, config.Options)
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		fmt.Println("数据库连接错误：", err)
		return
	}

	fmt.Println("数据库连接成功")

	err = DB.Set("gorm:table_options", " COMMENT='用户表'").AutoMigrate(&User{})
	if err != nil {
		fmt.Println("初始化用户表失败", err)
		return
	}
	err = DB.Set("gorm:table_options", " COMMENT='视频表'").AutoMigrate(&Video{})
	if err != nil {
		fmt.Println("初始化视频表失败", err)
		return
	}
	err = DB.Set("gorm:table_options", " COMMENT='消息表'").AutoMigrate(&Message{})
	if err != nil {
		fmt.Println("初始化消息表失败", err)
		return
	}
	err = DB.Set("gorm:table_options", " COMMENT='关注关系表'").AutoMigrate(&Follow{})
	if err != nil {
		fmt.Println("初始化关注关系表失败", err)
		return
	}
	err = DB.Set("gorm:table_options", " COMMENT='点赞记录表'").AutoMigrate(&Favorite{})
	if err != nil {
		fmt.Print("初始化点赞记录表失败", err)
		return
	}
	err = DB.Set("gorm:table_options", " COMMENT='评论表'").AutoMigrate(&Comment{})
	if err != nil {
		fmt.Println("初始化评论表失败", err)
		return
	}
}
