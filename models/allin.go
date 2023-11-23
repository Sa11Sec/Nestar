package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

// AllinCount 汇总的相关函数
type AllinCount struct {
	gorm.Model
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"not null" json:"name"`
	Count uint   `gorm:"not null" json:"count"`
}

type AllinDataUsername struct {
	gorm.Model
	ID  uint   `gorm:"primaryKey" json:"id"`
	Cnt string `gorm:"not null" json:"cnt"`
	Lv  uint   `gorm:"not null" json:"lv"`
}

type AllinDataPasswd struct {
	gorm.Model
	ID  uint   `gorm:"primaryKey" json:"id"`
	Cnt string `gorm:"not null" json:"cnt"`
	Lv  uint   `gorm:"not null" json:"lv"`
}

type AllinDataDataCollection struct {
	Name   string
	TopAll []Data
	Count  uint
}

// GetAllinCount 获取 AllinData 的数量
func GetAllinCount(where string) uint {
	var allinCount AllinCount
	err = db.Where("name = ?", strings.ToLower(where)).First(&allinCount).Error
	if err != nil {
		log.Fatalf("查询 top all 数量错误：%v", err)
	}
	return allinCount.Count
}

//func getOneTopAllData(reader *bufio.Reader) (bool, Data) {
//	line, _, err := reader.ReadLine()
//	if err == io.EOF {
//		return true, Data{Cnt: "", Lv: 0}
//	}
//	if err != nil {
//		log.Fatalf("GetAllinData 文件读取错误 %v", err)
//	}
//	tempData := strings.Split(string(line), "🍊")
//	if len(tempData) == 2 {
//		cnt := tempData[0]
//		lv, _ := strconv.ParseUint(tempData[1], 10, 64)
//		return false, Data{Cnt: cnt, Lv: uint(lv)}
//	}
//	// 不等于 2 说明是个空行，这种情况是要避免的，没办法，只能当成文件结束了。
//	return true, Data{Cnt: "", Lv: 0}
//}

// GetAllin 获取 Allin 的信息
func GetAllin() []AllinDataDataCollection {
	UsernameDC := AllinDataDataCollection{
		Name:  "Username",
		Count: GetAllinCount("username"),
	}
	PasswdDC := AllinDataDataCollection{
		Name:  "Passwd",
		Count: GetAllinCount("passwd"),
	}
	return []AllinDataDataCollection{UsernameDC, PasswdDC}
}

// AddAllinData 向 AllinData 中整理数据
func AddAllinData(s SortQuery) {
	// 读取对应的 temp.json
	tempPath := "data/" + s.Where + "/temp.json"
	tempFileHandle, _ := os.Open(tempPath)
	var tempJson DataSliceDecrement
	err = json.NewDecoder(tempFileHandle).Decode(&tempJson)
	if err != nil {
		log.Fatalf("SortAllin() json 解析错误：%v", err)
	}
	err := tempFileHandle.Close()
	if err != nil {
		return
	}

	if s.Where == "Username" {
		AddUsername(tempJson)
	} else if s.Where == "Passwd" {
		AddPasswd(tempJson)
	} else {
		log.Fatalf("AddAllinData() error : 不正确的 Where 参数")
	}
	// 将 tempPath 清空 还没写
	tempFileHandle, _ = os.OpenFile(tempPath, os.O_WRONLY|os.O_TRUNC, 0644) // 只写模式|打开文件时清空文件
	err = json.NewEncoder(tempFileHandle).Encode([]Data{})
	if err != nil {
		return
	}
}

func AddUsername(tempJson DataSliceDecrement) {

	for _, data := range tempJson {
		var allinDataUsername []AllinDataUsername
		db.Debug().Where("cnt = ?", data.Cnt).Find(&allinDataUsername)
		if len(allinDataUsername) == 1 {
			// 如果存在呢，就合并
			allinDataUsername[0].Lv += data.Lv
			db.Where("id = ?", allinDataUsername[0].ID).Updates(&allinDataUsername[0])
		} else if len(allinDataUsername) == 0 {
			// 如果不存在，则添加
			db.Create(&AllinDataUsername{Cnt: data.Cnt, Lv: data.Lv})
		} else {
			log.Fatalf("发生了莫名奇妙的错误")
		}
	}
	var count int64
	db.Model(&AllinDataUsername{}).Count(&count)

	var allinCount AllinCount
	db.Debug().Where("name = ?", "username").First(&allinCount)
	allinCount.Count = uint(count)
	db.Debug().Where("id = ?", allinCount.ID).Updates(&allinCount)
}

func AddPasswd(tempJson DataSliceDecrement) {
	for _, data := range tempJson {
		var allinDataPasswd []AllinDataPasswd
		db.Debug().Where("cnt = ?", data.Cnt).Find(&allinDataPasswd)
		if len(allinDataPasswd) == 1 {
			// 如果存在呢，就合并
			allinDataPasswd[0].Lv += data.Lv
			db.Debug().Where("id = ?", allinDataPasswd[0].ID).Updates(&allinDataPasswd[0])
		} else if len(allinDataPasswd) == 0 {
			// 如果不存在，则添加
			db.Create(&AllinDataPasswd{Cnt: data.Cnt, Lv: data.Lv})
		} else {
			log.Fatalf("11111111 error")
		}
	}

	var count int64
	db.Model(&AllinDataPasswd{}).Count(&count)

	var allinCount AllinCount
	db.Where("name = ?", "passwd").First(&allinCount)
	allinCount.Count = uint(count)
	db.Debug().Where("id = ?", allinCount.ID).Updates(&allinCount)
}

// GetAllinData 获取多少条数据
func GetAllinData(copyQuery CopyQuery) string {
	if copyQuery.Where == "Username" {
		var allinDataUsername []AllinDataUsername
		return GetAllinDataHelp(allinDataUsername, copyQuery.Count)
	} else if copyQuery.Where == "Passwd" {
		var allinDataPasswd []AllinDataPasswd
		return GetAllinDataHelp(allinDataPasswd, copyQuery.Count)
	} else {
		log.Fatalf("AddAllinData() error : 不正确的 Where 参数")
		return ""
	}
}

func GetAllinDataHelp(oneType interface{}, count uint) string {
	db.Debug().Order("lv desc").Limit(int(count)).Find(&oneType)
	result := ""
	switch targetType := oneType.(type) {
	case []AllinDataUsername:
		for _, data := range targetType {
			result += data.Cnt + "\n"
		}
	case []AllinDataPasswd:
		for _, data := range targetType {
			result += data.Cnt + "\n"
		}
	}
	return result
}

//func GetUsername(count uint) string {
//	var allinDataUsername []AllinDataUsername
//	db.Debug().Order("lv desc").Limit(int(count)).Find(&allinDataUsername)
//	result := ""
//	for _, data := range allinDataUsername {
//		result += data.Cnt + "\n"
//	}
//	return result
//}
//
//func GetPasswd(count uint) string {
//	var allinDataPasswd []AllinDataPasswd
//	db.Debug().Order("lv desc").Limit(int(count)).Find(&allinDataPasswd)
//	result := ""
//	for _, data := range allinDataPasswd {
//		result += data.Cnt + "\n"
//	}
//	return result
//}
