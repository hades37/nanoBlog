package dao

import (
	"fmt"
	"sync"

	"nanoBlog/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DBConfig 数据库配置

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return db
}

// InitDB 初始化数据库连接
func InitDB(conf *config.DBConfig) error {
	var err error
	once.Do(func() {
		err = connectDB(conf)
	})

	return err
}

// connectDB 连接数据库
func connectDB(conf *config.DBConfig) error {
	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Charset,
	)

	// GORM 配置
	gormConfig := &gorm.Config{
		// 命名策略
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
		// 日志配置
		Logger: logger.Default.LogMode(logger.Info),
	}

	// 建立连接
	gormDB, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	// 获取通用数据库对象
	sqlDB, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(conf.MaxLifetime)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %v", err)
	}

	db = gormDB
	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
