package chat

import (
	"context"
	"fmt"
	"iter"
	"os"
	"strings"

	"google.golang.org/genai"
)

var ctx context.Context = context.Background()
var client *genai.Client

const PROMPT_MENTAL_HEALTH = `
You are a mental health chatbot. You can answer any question related to mental health. Please answer the question as best as you can.

Use makrdown format to make the answer more readable. Use bullet points if needed.
`

func ClientGenAI() *genai.Client {
	var err error
	client, err = genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  os.Getenv("GEMINI_KEY"),
		Backend: genai.BackendGeminiAPI,
	})

	if err != nil {
		panic(err)
	}

	return client
}

func BuatChat(sejarah []*genai.Content) *genai.Chat {
	history := []*genai.Content{
		genai.NewContentFromText("Hello! I'm your MindfulAI assistant. How are you feeling today?", genai.RoleModel),
	}

	history = append(history, sejarah...)

	var config = &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(PROMPT_MENTAL_HEALTH, genai.RoleUser),
	}

	// There is a bug for Chatbot in adding history. It suddenly forget everything what the user asked after 2 or more user's chat. [23-05-20025 22:11]
	fmt.Println("SEJARAH: ", len(history))
	chat, err := client.Chats.Create(ctx, "gemini-2.0-flash", config, history)
	if err != nil {
		panic(err)
	}

	return chat
}

func KirimPesan(cs *genai.Chat, pesan string) (*genai.GenerateContentResponse, error) {
	return cs.SendMessage(ctx, genai.Part{Text: pesan})
}

func KirimPesanStream(cs *genai.Chat, pesan string) iter.Seq2[*genai.GenerateContentResponse, error] {
	return cs.GenerateContentStream(ctx, "gemini-2.0-flash", genai.Text(pesan), nil)
}

func DapatinJudulPercakapan(pesan string) string {
	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text("Write ONLY ONE title (maximum 30 words) for this conversation based on the following text: "+pesan+"\nNo need to add any other text. Just write the title with no dramatic."),
		nil,
	)

	return strings.Replace(result.Text(), "\n", "", -1)
}

func DapatinUrgencyLevel(pesan string) int {
	result, _ := client.Models.GenerateContent(
		ctx,
		"gemini-2.0-flash",
		genai.Text("Write ONLY ONE urgency level (1-5) for this conversation based on the following text: "+pesan+"\nNo need to add any other text. Just write the urgency level."),
		nil,
	)

	urgencyLevel := strings.Replace(result.Text(), "\n", "", -1)
	if urgencyLevel == "" {
		return 1
	}

	return int(urgencyLevel[0] - '0')
}

// func DapatinHistoryKarakter(karakterChat models.KarakterChat) ([]*genai.Content, []models.DataHistoryChat) {
// 	genAIHistoryChat := []*genai.Content{}
// 	dataHistoryChat := []models.DataHistoryChat{}

// 	sort.Slice(karakterChat.History, func(i, j int) bool {
// 		return karakterChat.History[i].Posisi < karakterChat.History[j].Posisi
// 	})

// 	for _, v := range karakterChat.History {
// 		genAIHistoryChat = append(genAIHistoryChat, &genai.Content{
// 			Role:  v.Role,
// 			Parts: []genai.Part{genai.Text(v.Chat)},
// 		})

// 		dataHistoryChat = append(dataHistoryChat, models.DataHistoryChat{
// 			ID:    v.ID,
// 			Chat:  v.Chat,
// 			Role:  v.Role,
// 			Waktu: v.CreatedAt,
// 		})
// 	}

// 	return genAIHistoryChat, dataHistoryChat
// }
