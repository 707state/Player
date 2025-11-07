package llm

import (
	"context"
	"log"
	"musaic/util"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

type LLMProvider struct {
	OllamaUrl string
	ModelName string
	Tools     []tool.BaseTool
}
type LLMResp struct {
	Resp []*schema.Message
	Err  error
}

func (provider *LLMProvider) Init() *LLMProvider {
	provider.ModelName = util.GetEnv("MODEL_NAME", "qwen3:1.7b")
	provider.OllamaUrl = util.GetEnv("OLLAMA_URL", "http://localhost:11434")
	return provider
}
func (provider *LLMProvider) AddTool(tool ...tool.BaseTool) *LLMProvider {
	provider.Tools = append(provider.Tools, tool...)
	return provider
}

func (provider *LLMProvider) ToolInfos(ctx context.Context) ([]*schema.ToolInfo, error) {
	toolInfos := make([]*schema.ToolInfo, 0, len(provider.Tools))
	for _, tool := range provider.Tools {
		info, err := tool.Info(ctx)
		if err != nil {
			return nil, err
		}
		toolInfos = append(toolInfos, info)
	}
	return toolInfos, nil
}

func RunAgent(inputCh <-chan string, responseCh chan<- LLMResp) {

	ctx := context.Background()
	llm_provider := &LLMProvider{}
	llm_provider.Init()
	// first, get tools
	chatTemplate := createChatModelTemplate()
	llm_provider.AddTool(GetTodayTool(), GetSinglesInfoTool())
	//query model is used to query data from MongoDB
	queryModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: llm_provider.OllamaUrl,
		Model:   llm_provider.ModelName,
	})
	if err != nil {
		log.Fatalf("Failed to create query model: %v", err)
	}
	// used to chat with data from query model tool calling
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: llm_provider.OllamaUrl,
		Model:   llm_provider.ModelName,
	})
	if err != nil {
		log.Fatalf("Failed to create chat model: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to create search model: %v", err)
	}
	// get tools info
	toolsInfo, err := llm_provider.ToolInfos(ctx)
	if err != nil {
		log.Fatalf("Failed to create tools info: %v", err)
	}

	err = queryModel.BindTools(toolsInfo)
	if err != nil {
		log.Fatalf("Failed to bind tools: %v", err)
	}

	musicToolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: llm_provider.Tools,
	})
	if err != nil {
		log.Fatalf("NewToolNode failed, err=%v", err)
	}
	const (
		nodeKeyOfTemplate      = "query_template"
		nodeKeyOfAlbumTemplate = "album_template"
		nodeKeyOfChatModel     = "chat_model"
		nodeKeyOfQueryModel    = "query_model"
		nodeKeyOfTools         = "tools"
		toList                 = "toList"
		toList2                = "toList2"
	)
	type state struct {
		currentRound int
		msgs         []*schema.Message
	}
	g := compose.NewGraph[map[string]any, []*schema.Message](compose.WithGenLocalState(func(ctx context.Context) *state { return &state{} }))
	_ = g.AddChatTemplateNode(nodeKeyOfTemplate, chatTemplate)

	_ = g.AddChatModelNode(nodeKeyOfChatModel, chatModel,
		compose.WithStatePreHandler(func(ctx context.Context, input []*schema.Message, state *state) ([]*schema.Message, error) {
			state.msgs = append(state.msgs, input...)
			input = append(createAlbumChatModelTemplate(), input[len(input)-1])
			return input, nil
		}),
		compose.WithStatePostHandler(func(ctx context.Context, input *schema.Message, state *state) (*schema.Message, error) {
			return input, nil
		}))
	_ = g.AddChatModelNode(nodeKeyOfQueryModel, queryModel)

	_ = g.AddToolsNode(nodeKeyOfTools, musicToolsNode,
		compose.WithStatePreHandler(func(ctx context.Context, input *schema.Message, state *state) (*schema.Message, error) {
			state.msgs = append(state.msgs, input)
			return input, nil
		}))

	_ = g.AddLambdaNode(toList, compose.ToList[*schema.Message]())

	_ = g.AddEdge(compose.START, nodeKeyOfTemplate)
	_ = g.AddEdge(nodeKeyOfTemplate, nodeKeyOfQueryModel)
	_ = g.AddEdge(nodeKeyOfQueryModel, nodeKeyOfTools)
	_ = g.AddEdge(nodeKeyOfTools, nodeKeyOfChatModel)
	_ = g.AddEdge(nodeKeyOfChatModel, toList)
	_ = g.AddEdge(toList, compose.END)
	agent, err := g.Compile(ctx)
	if err != nil {
		log.Fatalf("chain.Compile failed, err=%v", err)
	}
	for input := range inputCh {
		// This will be input again
		resp, err := agent.Invoke(ctx, createChatModelMessagesFromTemplate(input))
		responseCh <- LLMResp{
			Resp: resp,
			Err:  err,
		}
	}
}
