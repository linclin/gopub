package column

import (
	"library/common"
	"library/paging"
	"strings"
)

func init() {
	paging.PagingColumns["Dictionary"] = Dictionary{}
}

type Dictionary struct{}

func (dictionary Dictionary) GetColumnValue(column common.Info, properties []string) string {
	ret := ""
	formats := strings.Split(column.Info("parms").String("format"), ";")
	dict := make(map[string]string, 0)
	for _, format := range formats {
		items := strings.Split(format, ":")
		if len(items) == 2 {
			dict[items[0]] = items[1]
		}
	}
	if properties[0] != "" {
		ret = dict[properties[0]]
	}
	return ret
}
