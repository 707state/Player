package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"musaic/db"
	"musaic/llm"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/v2"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

func TestLLMTool(t *testing.T) {
	ctx := context.Background()
	update_tool, err := utils.InferTool(
		"update_tool",
		"Update a todo item, like content, deadline...",
		llm.UpdateTodoFunc)
	if err != nil {
		t.Errorf("failed InferTool: %v", err)
		return
	}

	// 创建 DuckDuckGo 工具
	searchTool, err := duckduckgo.NewTextSearchTool(ctx, &duckduckgo.Config{})
	if err != nil {
		t.Errorf("NewTextSearchTool failed, err=%v", err)
		return
	}

	llm_provider := &llm.LLMProvider{}
	llm_provider.Init().AddTool(
		update_tool,
		llm.GetAddTodoTool(),
		&llm.ListTodoTool{},
		searchTool,
	)
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: llm_provider.OllamaUrl,
		Model:   llm_provider.ModelName,
	})
	if err != nil {
		t.Errorf("NewChatModel failed, err=%v", err)
		return
	}
	// 获取工具信息, 用于绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(llm_provider.Tools))
	var info *schema.ToolInfo
	for _, todoTool := range llm_provider.Tools {
		info, err = todoTool.Info(ctx)
		if err != nil {
			t.Errorf("get ToolInfo failed, err=%v", err)
			return
		}
		toolInfos = append(toolInfos, info)
	}

	// 将 tools 绑定到 ChatModel
	err = chatModel.BindTools(toolInfos)
	if err != nil {
		t.Errorf("BindTools failed, err=%v", err)
		return
	}

	// 创建 tools 节点
	todoToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: llm_provider.Tools,
	})
	if err != nil {
		t.Errorf("NewToolNode failed, err=%v", err)
		return
	}

	// 构建完整的处理链
	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	chain.
		AppendChatModel(chatModel, compose.WithNodeName("chat_model")).
		AppendToolsNode(todoToolsNode, compose.WithNodeName("tools"))

	// 编译并运行 chain
	agent, err := chain.Compile(ctx)
	if err != nil {
		t.Errorf("chain.Compile failed, err=%v", err)
		return
	}
	// 运行示例
	resp, err := agent.Invoke(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "添加一个学习 Eino 的 TODO，然后帮我搜索一下Eino的资料，最后帮我列举一下我的TODO项，非常感谢",
		},
	})
	if err != nil {
		t.Errorf("agent.Invoke failed, err=%v", err)
		return
	}

	// 输出结果
	for idx, msg := range resp {
		var pretty_json bytes.Buffer
		error := json.Indent(&pretty_json, []byte(msg.Content), "", "\t")
		if error != nil {
			t.Errorf("JSON Parse error: %v", error)
		}
		log.Printf("message %d: %s: %s", idx, msg.Role, pretty_json.String())
	}
}

func TestDBQuerySinglesFromMongo(t *testing.T) {
	client := db.Init()
	defer client.Disconnect(context.TODO())
	ctx := context.Background()
	searchParam := &db.QueryItem{
		Album: "Whole Lotta Red",
	}
	kv, err := db.QuerySinglesFromMongoDB(ctx, searchParam)
	if err != nil {
		t.Error(err)
	}
	if kv["count"] == 0 {
		t.Errorf("Failed to query data from MongoDB")
	}
	log.Printf("%v\n", kv)
	searchParam = &db.QueryItem{
		Duration: db.JSONDuration(7 * time.Hour),
	}
	kv, err = db.QuerySinglesFromMongoDB(ctx, searchParam)
	if err != nil {
		t.Error(err)
	}
	if kv["count"] == 0 {
		t.Errorf("Failed to query data with duration from mongodb")
	}
}

func TestLLMQuery(t *testing.T) {
	llm.Init()
	db.Init()
	question := "请你帮我查一下最近3天（从今天起）的我的音乐数据，然后提供一些总结吧"
	llm.InputChan <- question
	resp := <-llm.ResponseChan
	if resp.Err != nil {
		t.Errorf("Failed to query llm")
	}
	log.Printf("Got result: %v", resp.Resp)
}

func TestDBQueryWithJson(t *testing.T) {
	db.Init()
	jsonStr := `{
        "album": "",
        "artists": [],
        "title": "",
        "duration": "10h",
        "start_point": "2025-11-04T00:00:00Z"
    }`
	var q db.QueryItem
	if err := json.Unmarshal([]byte(jsonStr), &q); err != nil {
		t.Fatalf("Failed to unmarshal json: %v", err)
	}
	kv, err := db.QuerySinglesFromMongoDB(t.Context(), &q)
	if err != nil {
		t.Fatalf("Failed to query: %s!!!", err)
	}
	if len(kv) == 0 || kv["count"] == 0 {
		t.Fatal("There must be some result, right?")
	}
	log.Println(kv)
}
