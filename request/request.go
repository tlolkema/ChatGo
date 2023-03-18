package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type InputBody struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}

type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index        int `json:"index"`
		Message      Message
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func init() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("OPENAI_API_KEY is not set")
		os.Exit(1)
	}
}

func MakeRequest(input []string) string {

	url := "https://api.openai.com/v1/chat/completions"

	var messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	preMessage := "for the whole conversation, output in markdown format and have max 70 characters per line"

	messages = append(messages, struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		Role:    "user",
		Content: preMessage,
	})

	for i := 0; i < len(input); i++ {
		if i%2 == 0 {
			messages = append(messages, struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			}{
				Role:    "user",
				Content: input[i],
			})
		} else {
			messages = append(messages, struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			}{
				Role:    "assistant",
				Content: input[i],
			})
		}
	}

	body := InputBody{
		Model:    "gpt-3.5-turbo",
		Messages: messages,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding request body:", err)
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + os.Getenv("OPENAI_API_KEY"),
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var chatResponse ChatResponse
	err = json.Unmarshal(respBody, &chatResponse)
	if err != nil {
		fmt.Println("Error decoding response body:", err)
		return ""
	}

	return chatResponse.Choices[0].Message.Content
}
