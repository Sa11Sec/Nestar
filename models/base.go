package models

import (
	"Nestar/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB
var err error

func init() {
	// 只复制建立数据库连接
	db, err = gorm.Open(sqlite.Open(config.DatabaseSetting.Name), &gorm.Config{})
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	log.Println("AllinCount 数据表初始化中...")
	table := db.Migrator().HasTable(&AllinCount{})
	if !table {
		err := db.AutoMigrate(&AllinCount{})
		if err != nil {
			log.Fatalf("初始化数据库错误：%v", err)
		}
	}

	// 去数据库里查询一下，AllinCount 是不是空的，没有就填充 0
	nameList := []string{"username", "passwd"}
	for _, name := range nameList {
		var allinCount AllinCount
		err := db.Where("name = ?", name).First(&allinCount).Error
		if err != nil {
			log.Println("在 AllinCount 数据表中插入数据")
			// 没有就直接写 0
			allinCount = AllinCount{Name: name, Count: 0}
			err := db.Create(&allinCount).Error
			if err != nil {
				log.Fatalf("AllinCount 表数据插入失败：%v", name)
			}
		}
	}

	log.Println("AllinData 数据表初始化中...")
	table = db.Migrator().HasTable(&AllinDataUsername{})
	if !table {
		err := db.AutoMigrate(&AllinDataUsername{})
		if err != nil {
			log.Fatalf("初始化数据库错误：%v", err)
		}
	}

	log.Println("AllinData 数据表初始化中...")
	table = db.Migrator().HasTable(&AllinDataPasswd{})
	if !table {
		err := db.AutoMigrate(&AllinDataPasswd{})
		if err != nil {
			log.Fatalf("初始化数据库错误：%v", err)
		}
	}

}
