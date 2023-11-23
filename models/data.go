package models

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"strings"
)

type Data struct {
	Cnt string `json:"cnt"`
	Lv  uint   `json:"lv"`
}

// DataCollection 普通分类的相关函数
type DataCollection struct {
	Name  string
	Top10 []Data
	Top1k []Data
}

func GetDataAll(systemName string) []DataCollection {
	UserNameDC := DataCollection{
		Name:  "Username",
		Top10: GetData("Username/top10/" + systemName),
		Top1k: GetData("Username/top1k/" + systemName),
	}
	PasswdDC := DataCollection{
		Name:  "Passwd",
		Top10: GetData("Passwd/top10/" + systemName),
		Top1k: GetData("Passwd/top1k/" + systemName),
	}
	return []DataCollection{UserNameDC, PasswdDC}
}

func GetData(path string) []Data {
	file, _ := os.Open("data/" + path + ".json")
	var info []Data
	err := json.NewDecoder(file).Decode(&info)
	if err != nil {
		return nil
	}
	return info
}

type AddDataQuery struct {
	SystemName string `json:"systemName"`
	Where      string `json:"where"`
	Clazz      string `json:"clazz"`
	Data       string `json:"data"`
}

type DataSliceDecrement []Data

func (d DataSliceDecrement) Len() int { return len(d) }

func (d DataSliceDecrement) Swap(i, j int) { d[i], d[j] = d[j], d[i] }

func (d DataSliceDecrement) Less(i, j int) bool { return d[i].Lv > d[j].Lv }

func isAlready(s string, info DataSliceDecrement) (int, bool) {
	for num, data := range info {
		if data.Cnt == s {
			return num, true
		}
	}
	return -1, false
}

// judgmentMerge 将 d1 的数据 ==排序合并==> 到 d2 中
func judgmentMerge(d1 DataSliceDecrement, d2 DataSliceDecrement) DataSliceDecrement {
	// 判断合并, 判断 top10 中的 在不在 top1k 中
	for numTop10, data := range d1 {
		numTop1k, validity := isAlready(data.Cnt, d2)
		if validity {
			// 如果 info 中有了，那么就 + 优先级
			d2[numTop1k].Lv += d1[numTop10].Lv
		} else {
			// 如果没有，那么就给加后面去，优先级 +1
			d2 = append(d2, data)
		}
	}
	return d2
}

func addData(data DataSliceDecrement, path string) {
	file, _ := os.Open(path)
	var info DataSliceDecrement
	err := json.NewDecoder(file).Decode(&info)
	if err != nil {
		log.Fatalf("data.go addData Json 绑定错误: %v", err)
	}
	err = file.Close()
	if err != nil {
		log.Fatalf("err %v", err)
	}

	info = judgmentMerge(data, info)
	sort.Sort(info)

	dataFile, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	err = json.NewEncoder(dataFile).Encode(&info)
	if err != nil {
		log.Fatalf("err %v", err)
	}
}

func AddData(query AddDataQuery) {
	path := "data/" + query.Where + "/" + query.Clazz + "/" + query.SystemName + ".json"
	// info 是 一个 Data 列表，然后 query 是一个字符串
	var tempInfo DataSliceDecrement
	query.Data = strings.TrimRight(query.Data, "\n")
	dataList := strings.Split(query.Data, "\n")
	for _, data := range dataList {
		tempInfo = append(tempInfo, Data{Cnt: data, Lv: 1})
	}

	addData(tempInfo, path)
}

func SortData(sq SortQuery) {
	pathTop10 := "data/" + sq.Where + "/" + "top10" + "/" + sq.SystemName + ".json"
	pathTop1k := "data/" + sq.Where + "/" + "top1k" + "/" + sq.SystemName + ".json"

	// 读取这两个文件
	fileTop10, _ := os.Open(pathTop10)
	fileTop1k, _ := os.Open(pathTop1k)

	var infoTop10 DataSliceDecrement
	var infoTop1k DataSliceDecrement

	err := json.NewDecoder(fileTop10).Decode(&infoTop10)
	if err != nil {
		return
	}
	err = json.NewDecoder(fileTop1k).Decode(&infoTop1k)
	if err != nil {
		return
	}

	infoTop1k = judgmentMerge(infoTop10, infoTop1k)
	// 排序
	sort.Sort(infoTop1k)
	// 分割
	var tempTop10 []Data
	var tempTop1k []Data
	var tempAllin []Data

	length := infoTop1k.Len()
	if length <= 0 {
		log.Fatalf("error")
	} else if length <= 10 {
		tempTop10 = infoTop1k[:length]
	} else if length <= 1000 {
		tempTop10 = infoTop1k[:10]
		tempTop1k = infoTop1k[10:length]
	} else { // 多余的，存入 temp.json
		tempTop10 = infoTop1k[:10]
		tempTop1k = infoTop1k[10:1000]
		tempAllin = infoTop1k[1000:length]
	}

	dataFileTop10, _ := os.OpenFile(pathTop10, os.O_WRONLY|os.O_TRUNC, 0644) // 只写模式|打开文件时清空文件
	dataFileTop1k, _ := os.OpenFile(pathTop1k, os.O_WRONLY|os.O_TRUNC, 0644) // 只写模式|打开文件时清空文件

	err = json.NewEncoder(dataFileTop10).Encode(&tempTop10)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	// 如果 top1k 为 nil ，那么删除内容
	if tempTop1k != nil {
		err = json.NewEncoder(dataFileTop1k).Encode(&tempTop1k)
		if err != nil {
			log.Fatalf("err %v", err)
		}
	} else {
		err = json.NewEncoder(dataFileTop1k).Encode([]Data{})
		if err != nil {
			log.Fatalf("err %v", err)
		}
	}

	// 如果 temp 不为空，则加载文件，合并文件 ！ 暂时未实现
	if tempAllin != nil {
		path := "data/" + sq.Where + "/" + "temp.json"
		addData(tempAllin, path)
	}
}
