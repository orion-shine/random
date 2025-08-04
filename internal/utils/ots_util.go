package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/joho/godotenv"
)

var client *tablestore.TableStoreClient
var tableName string

func init() {
	_ = godotenv.Load()
	endpoint := os.Getenv("ENDPOINT")
	instance := os.Getenv("INSTANCE_NAME")
	accessKey := os.Getenv("ACCESS_KEY_ID")
	accessSecret := os.Getenv("ACCESS_KEY_SECRET")
	tableName = os.Getenv("TABLE_NAME")
	client = tablestore.NewClient(endpoint, instance, accessKey, accessSecret)
}

func createTableIfNotExists(client *tablestore.TableStoreClient, tableName string) {
	// 先列出已有表，判断是否已存在
	listResp, err := client.ListTable()
	if err != nil {
		fmt.Println("获取表列表失败:", err)
		return
	}
	for _, existingTable := range listResp.TableNames {
		if existingTable == tableName {
			fmt.Printf("表 %s 已存在，无需重新创建\n", tableName)
			return
		}
	}
	// 构造表结构
	tableMeta := new(tablestore.TableMeta)
	tableMeta.TableName = tableName
	tableMeta.AddPrimaryKeyColumn("time", tablestore.PrimaryKeyType_STRING)
	// 设置表选项（永久存储 + 单版本）
	tableOption := &tablestore.TableOption{
		TimeToAlive: -1, // 永不过期
		MaxVersion:  1,  // 最多保留一个版本
	}
	// 设置预留吞吐量为0（建议使用按量付费）
	reservedThroughput := &tablestore.ReservedThroughput{
		Readcap:  0,
		Writecap: 0,
	}
	// 构造请求
	createReq := &tablestore.CreateTableRequest{
		TableMeta:          tableMeta,
		TableOption:        tableOption,
		ReservedThroughput: reservedThroughput,
	}

	// 创建表
	_, err = client.CreateTable(createReq)
	if err != nil {
		fmt.Println("创建表失败:", err)
		return
	}
	fmt.Printf("表 %s 创建成功\n", tableName)
}

func putFrameData(data map[string]string) {
	pk := new(tablestore.PrimaryKey)
	pk.AddPrimaryKeyColumn("time", data["time"]) // 主键必须存在，确保唯一

	putRowChange := new(tablestore.PutRowChange)
	putRowChange.TableName = tableName
	putRowChange.PrimaryKey = pk
	putRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)

	// 字段列表
	keys := []string{"c1", "c2", "c3", "c4", "hex", "img1", "img2", "img3", "img4", "tid", "rand"}

	for _, key := range keys {
		val, ok := data[key]
		if !ok || val == "" {
			continue // 跳过空字段，避免 panic
		}

		// 如果字段是 hex 或 img 开头的，可能是二进制，转为 []byte
		if strings.HasPrefix(key, "img") || key == "hex" {
			putRowChange.AddColumn(key, []byte(val))
		} else {
			putRowChange.AddColumn(key, val) // 默认当作字符串写入
		}
	}

	// 创建请求并提交
	putReq := new(tablestore.PutRowRequest)
	putReq.PutRowChange = putRowChange

	_, err := client.PutRow(putReq)
	if err != nil {
		fmt.Println("写入失败:", err)
	} else {
		fmt.Println("写入成功")
	}
}

func getLatestFrames(limit int, columns []string, filterCond map[string]string, timeRange map[string]string) []map[string]interface{} {
	startPK := new(tablestore.PrimaryKey)
	endPK := new(tablestore.PrimaryKey)

	if timeRange != nil {
		start, _ := strconv.ParseInt(timeRange["start_time"], 10, 64)
		end, _ := strconv.ParseInt(timeRange["end_time"], 10, 64)
		if start > end {
			start, end = end, start
		}
		// 设定查询时间范围
		startPK.AddPrimaryKeyColumn("time", strconv.FormatInt(start, 10))
		endPK.AddPrimaryKeyColumn("time", strconv.FormatInt(end+1, 10)) // exclusive
	} else {
		startPK.AddPrimaryKeyColumnWithMinValue("time")
		endPK.AddPrimaryKeyColumnWithMaxValue("time")
	}
	rangeCriteria := &tablestore.RangeRowQueryCriteria{
		TableName:       tableName,
		StartPrimaryKey: startPK,
		EndPrimaryKey:   endPK,
		Direction:       tablestore.BACKWARD, // 从新到旧（DESCENDING）
		MaxVersion:      1,
		Limit:           int32(limit),
		ColumnsToGet:    columns,
	}
	// 构建过滤条件
	if filterCond != nil {
		cond := tablestore.NewCompositeColumnCondition(tablestore.LO_AND)
		if hex, ok := filterCond["hex"]; ok {
			cond.AddFilter(tablestore.NewSingleColumnCondition("hex", tablestore.CT_EQUAL, hex))
		}
		if randStr, ok := filterCond["rand"]; ok {
			if randInt, err := strconv.Atoi(randStr); err == nil {
				cond.AddFilter(tablestore.NewSingleColumnCondition("rand", tablestore.CT_EQUAL, randInt))
			}
		}
		rangeCriteria.Filter = cond
	}
	getRangeReq := &tablestore.GetRangeRequest{
		RangeRowQueryCriteria: rangeCriteria,
	}
	var result []map[string]interface{}
	for {
		getRangeResp, err := client.GetRange(getRangeReq)
		if err != nil {
			fmt.Println("GetRange 查询失败:", err)
			break
		}
		// 处理返回行
		result = append(result, transformFrameData(getRangeResp.Rows)...)
		// 是否还有下一页
		if getRangeResp.NextStartPrimaryKey == nil || len(result) >= limit {
			break
		}
		// 更新起始主键用于翻页
		getRangeReq.RangeRowQueryCriteria.StartPrimaryKey = getRangeResp.NextStartPrimaryKey
	}
	// 截断到 limit 条
	if len(result) > limit {
		return result[:limit]
	}
	return result
}

func getOne(primaryKey string, filterCond map[string]string) map[string]interface{} {
	pk := new(tablestore.PrimaryKey)
	pk.AddPrimaryKeyColumn("time", primaryKey)

	criteria := new(tablestore.SingleRowQueryCriteria)
	criteria.TableName = tableName
	criteria.PrimaryKey = pk
	criteria.MaxVersion = 1

	// 构建过滤条件
	if filterCond != nil {
		cond := tablestore.NewCompositeColumnCondition(tablestore.LO_AND)
		if hex, ok := filterCond["hex"]; ok {
			cond.AddFilter(tablestore.NewSingleColumnCondition("hex", tablestore.CT_EQUAL, hex))
		}
		if randStr, ok := filterCond["rand"]; ok {
			if randInt, err := strconv.Atoi(randStr); err == nil {
				cond.AddFilter(tablestore.NewSingleColumnCondition("rand", tablestore.CT_EQUAL, randInt))
			}
		}
		criteria.Filter = cond
	}

	getReq := new(tablestore.GetRowRequest)
	getReq.SingleRowQueryCriteria = criteria

	resp, err := client.GetRow(getReq)
	if err != nil {
		fmt.Println("读取失败:", err)
		return nil
	}
	if resp.Columns[0].ColumnName != "" && resp.Columns[0].Value != nil {
		var rows []*tablestore.Row

		row := &tablestore.Row{
			PrimaryKey: &resp.PrimaryKey,
			Columns:    resp.Columns,
		}
		rows = append(rows, row)
		data := transformFrameData(rows)
		return data[0]
	}

	return nil
}

// 将原始帧数据转换为标准格式
func transformFrameData(rows []*tablestore.Row) []map[string]interface{} {
	result := []map[string]interface{}{}
	for _, row := range rows {
		item := make(map[string]interface{})
		for _, pk := range row.PrimaryKey.PrimaryKeys {
			if pk.ColumnName == "time" {
				item["time"] = pk.Value
			}
		}
		for _, col := range row.Columns {
			item[col.ColumnName] = col.Value
		}
		result = append(result, item)
	}
	return result
}

func TestGetOne(t *testing.T) {
	createTableIfNotExists(client, tableName)

	filter := map[string]string{
		"hex":  "126afd38af5dfad200c49c53deb1",
		"rand": "845565712",
	}

	result := getOne("1753353553733", filter)
	if result == nil {
		t.Error("未获取到数据")
	} else {
		t.Logf("获取到数据: %+v", result)
	}
}
