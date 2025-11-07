package llm

import (
	"context"
	"musaic/db"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

func GetSinglesInfoTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "query_singles",
		Desc: "Query Singles from MongoDB with given params",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"album": {
				Desc:     "The album to search for, not a must",
				Type:     schema.String,
				Required: false,
			},
			"artists": {
				Desc:     "The artists to search for, not a must",
				Type:     schema.Array,
				Required: false,
			},
			"title": {
				Desc:     "The title of single to search for, not a must",
				Type:     schema.String,
				Required: false,
			},
			"duration": {
				Desc:     "The duration of singles updated. The format is like 10d(10 days), 72h(72hours)...",
				Type:     schema.String,
				Required: false,
			},
			"start_point": {
				Desc:     "The time point to start searching. The format is RFC3339, like yyyy-mm-ddThh:mm:ssZ. If this field is not needed, please pass null instead of empty string.",
				Required: false,
			},
		}),
	}
	return utils.NewTool(info, db.QuerySinglesFromMongoDBFunc)
}

func GetTodayTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "get_today",
		Desc: "Get today tool, this tool doesn't need parameter, it returns immediately/获取今日日期，这个工具不需要传递参数，直接返回",
	}
	return utils.NewTool(info, getTodayInfo)
}
func getTodayInfo(ctx context.Context, params any) (string, error) {
	return time.Now().In(time.UTC).String(), nil
}
