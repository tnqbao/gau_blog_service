package ai

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type TagResponse struct {
	Tag string `json:"tag"`
}

// CleanAIResponse xử lý kết quả trả về từ AI để đảm bảo JSON hợp lệ.
func CleanAIResponse(response string) string {
	response = strings.TrimSpace(response)
	response = strings.Trim(response, "```json")
	response = strings.Trim(response, "```")
	response = strings.Trim(response, "`")
	response = strings.Trim(response, `\`)
	return response
}

func DetectTag(body string) string {
	if body == "" {
		return ""
	}

	prompt := fmt.Sprintf(`Defend paragraph, generate JSON following this format:
	{
		"tag": "tag_name"
	}
	example:
	{
		"tag": "18+"
	}

	paragraph: %s

	Note: Only return JSON, no other text!
	Prioritize tag selection by age rating (18+, 13+, if no sensitive words then move to lower priority tags) -> story type (spiritual/ detective/ horror)
    * Only 5 tags are allowed: 18+, 13+, Tâm Linh, Kỳ Án, Kinh Dị    
`, body)

	aiClient := NewAIClient()
	if aiClient == nil {
		log.Println("AI client not initialized")
		return ""
	}

	resp, err := aiClient.GetAIResponse(prompt)
	if err != nil {
		log.Println("Error getting AI response:", err)
		return ""
	}

	cleanResp := CleanAIResponse(resp)

	var tagResponse TagResponse
	if err := json.Unmarshal([]byte(cleanResp), &tagResponse); err != nil {
		log.Println("Error parsing AI response:", err)
		return ""
	}

	return tagResponse.Tag
}
