package wuid

import (
	"database/sql"
	"fmt"
	"github.com/edwingeng/wuid/mysql/wuid"
	"sort"
	"strconv"
)

var w *wuid.WUID

func Init(dsn string) {
	newDB := func() (*sql.DB, bool, error) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, false, err
		}
		return db, true, nil
	}

	w = wuid.NewWUID("default", nil)
	_ = w.LoadH28FromMysql(newDB, "wuid")
}

func GenUid(dsn string) string {
	if w == nil {
		Init(dsn)
	}
	return fmt.Sprintf("%#016x", w.Next())
}

func CombineId(aId, bId string) string {
	ids := []string{aId, bId}
	// 对ids切片进行排序，确保其元素按升序排列。
	// 该排序基于将切片中的字符串元素转换为uint64类型数值进行比较，
	// 以确保排序结果是根据数字大小而非字符串顺序。
	sort.Slice(ids, func(i, j int) bool {
		// 将ids[i]转换为uint64类型，忽略错误，因为假设所有元素都是有效的数字字符串。
		a, _ := strconv.ParseUint(ids[i], 10, 64)
		// 将ids[j]转换为uint64类型，同样忽略错误。
		b, _ := strconv.ParseUint(ids[j], 10, 64)
		// 返回a < b，用于sort.Slice判断元素顺序。
		return a < b
	})
	return fmt.Sprintf("%s_%s", ids[0], ids[1])
}
