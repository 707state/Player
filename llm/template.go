package llm

import (
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func createChatModelTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题，你有get_today和query_singles两个工具调用，前者可以知道今天的日期，后者可以通过API查询数据，如果存在数据的话，你也应该根据数据进行回答。你的目标是{goal}。"),
		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
}

func createChatModelMessagesFromTemplate(question string) map[string]any {
	// 使用模板生成消息
	return map[string]any{
		"role":     "具有丰富音乐知识和素养的唱片收藏家、音乐家以及音乐学院教授，对于各种音乐都有涉猎",
		"style":    "积极、温暖且专业",
		"question": question,
		"goal":     "帮助人们提供建议，同时针对一些不好的倾向提供指导，你的回答应该是一份完整且具有指导意义的表格",
	}
}

func createAlbumChatModelTemplate() []*schema.Message {
	return []*schema.Message{
		schema.SystemMessage("你是一个对于音乐独具慧眼的唱片收藏家和音乐家，同时也是我的专职顾问，我会给你一些json格式的数据，你需要根据这些数据进行总结，你需要以一个老道的收藏家的身份给出总结并提出建议。"),
	}
}
