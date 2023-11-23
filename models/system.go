package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"os"
)

type System struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	IsCommon bool   `gorm:"default:false" json:"is_common"`
}

func init() {
	// 负责初始化
	table := db.Migrator().HasTable(&System{})
	if table {
		return
	}
	err := db.AutoMigrate(&System{})
	if err != nil {
		return
	}
}

func GetSystemCount() (c int64) {
	db.Debug().Table("systems").Count(&c)
	return c
}

func GetSystemUniversalList() []System {
	var system []System
	err := db.Where("is_common = ?", true).Find(&system).Error
	if err != nil {
		return nil
	}
	return system
}

func GetSystemSearch(search string) []System {
	var system []System
	// like 搜索
	err := db.Where("name like ?", "%"+search+"%").Find(&system).Error
	if err != nil {
		return nil
	}
	if len(system) == 0 {
		return []System{{
			Name: "没查询到 " + search,
		}}
	}
	return system
}

func AddSystem(addSystem AddSystemQuery) string {
	// 先查询一下这个系统名是否已经存在了
	var system []System
	db.Debug().Where("name = ?", addSystem.SystemName).Find(&system)
	if len(system) == 1 {
		// 说明已经存在了
		return "该系统已存在"
	} else if len(system) == 0 {
		// 说明不存在
		db.Debug().Create(&System{Name: addSystem.SystemName, IsCommon: addSystem.Type})
		pathList := []string{"Username/top10", "Username/top1k", "Passwd/top10", "Passwd/top1k"}

		for _, tempPath := range pathList {
			filePath := "data/" + tempPath + "/" + addSystem.SystemName + ".json"
			fileHandler, _ := os.Create(filePath)
			err := json.NewEncoder(fileHandler).Encode([]Data{})
			if err != nil {
				log.Fatalf("发生错误")
			}
		}
		return "添加成功"
	} else {
		// 说明出问题了
		return "发生错误"
	}
}
