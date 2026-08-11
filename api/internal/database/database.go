package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/doopush/doopush/api/internal/config"
	"github.com/doopush/doopush/api/internal/models"
	"github.com/doopush/doopush/api/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	connectMaxAttempts = 10
	connectRetryDelay  = 3 * time.Second
)

// Connect 连接数据库；MySQL 容器冷启动期间端口已开但服务未就绪，最多重试 10 次。
func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetString("DB_USERNAME", "root"),
		config.GetString("DB_PASSWORD", "password"),
		config.GetString("DB_HOST", "mysql"),
		config.GetString("DB_PORT", "3306"),
		config.GetString("DB_DATABASE", "doopush"),
	)

	var (
		db    *gorm.DB
		sqlDB *sql.DB
		err   error
	)
	for attempt := 1; attempt <= connectMaxAttempts; attempt++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			if sqlDB, err = db.DB(); err == nil {
				if err = sqlDB.Ping(); err == nil {
					break
				}
			}
		}
		if attempt == connectMaxAttempts {
			log.Fatalf("数据库连接失败（已重试 %d 次）: %v", connectMaxAttempts, err)
		}
		log.Printf("数据库连接失败（第 %d/%d 次），%s 后重试: %v",
			attempt, connectMaxAttempts, connectRetryDelay, err)
		time.Sleep(connectRetryDelay)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db
	log.Println("数据库连接成功")
}

// AutoMigrate 自动迁移数据表
func AutoMigrate() {
	// 执行自动迁移
	if err := DB.AutoMigrate(models.AllModels()...); err != nil {
		log.Fatal("数据库迁移失败:", err)
	}

	// 一次性清理 TCP 时代的死字段（GORM AutoMigrate 不会删列）
	_ = DB.Migrator().DropColumn(&models.Device{}, "gateway_node")
	_ = DB.Migrator().DropColumn(&models.Device{}, "connection_id")

	// 为升级前创建的应用补齐 App Key，并替换认证改造期间生成的旧格式 Key。
	var apps []models.App
	if err := DB.Where("app_key = '' OR app_key IS NULL OR app_key NOT REGEXP '^dp_ak_[A-Za-z0-9]{32}$'").Find(&apps).Error; err == nil {
		for i := range apps {
			appKey := utils.GenerateSecureToken(models.AppKeyPrefix)
			if err := DB.Model(&apps[i]).Update("app_key", appKey).Error; err != nil {
				log.Fatal("App Key回填失败:", err)
			}
		}
	} else {
		log.Fatal("查询待回填App Key失败:", err)
	}

	// 必须在旧数据回填之后创建唯一索引，否则多个空值会使升级迁移失败。
	if !DB.Migrator().HasIndex(&models.App{}, "idx_apps_app_key") {
		if err := DB.Exec("CREATE UNIQUE INDEX idx_apps_app_key ON apps (app_key)").Error; err != nil {
			log.Fatal("创建App Key唯一索引失败:", err)
		}
	}

	// 为升级前的用户审计记录补齐显式主体信息。
	if err := DB.Model(&models.AuditLog{}).
		Where("principal_id = 0 AND user_id > 0").
		Updates(map[string]interface{}{"principal_type": "user", "principal_id": gorm.Expr("user_id")}).Error; err != nil {
		log.Fatal("审计主体回填失败:", err)
	}

	log.Println("数据库迁移完成")
}
