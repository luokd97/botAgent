package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

const robotApiUrl = "http://cors.redoc.ly/core/engine/dm/bot-response"
const ipApiUrl = "http://ip-api.com/json"

func GetPublicIp(c *gin.Context) {
	resp, err := http.Get(ipApiUrl)
	if err == nil {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("getPublicIp() respBody = %v\n", string(respBody))
		var data map[string]interface{}
		err := json.Unmarshal(respBody, &data)
		if err == nil {
			id := fmt.Sprint(data["zip"])
			name := fmt.Sprint(data["country"])
			go AddRecord(id, name)
			c.IndentedJSON(http.StatusOK, data)
		} else {
			c.IndentedJSON(http.StatusOK, string(respBody))
			panic(err)
		}
	} else {
		fmt.Println("getPublicIp() Request error")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// @Summary 代理机器人Api
// @Description 见原始文档https://newtestchatbot.wul.ai/core/engine/dm/openapi.json
// @Tags botAgent
// @Accept json
// @Produce json
// @Success 200 body string
// @Router /bot [post]
func GetBotResponse(c *gin.Context) {
	resp, err := http.Post(robotApiUrl, "application/json", c.Request.Body)

	if err == nil {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Printf("getBotResponse() respBody = %v\n", string(respBody))
		var data map[string]interface{}
		err := json.Unmarshal(respBody, &data)
		if err == nil {
			intent := data["intent"].(map[string]interface{})
			id := fmt.Sprint(intent["id"])
			name := fmt.Sprint(intent["name"])
			go AddRecord(id, name)
			c.IndentedJSON(http.StatusOK, data)
		} else {
			c.IndentedJSON(http.StatusOK, string(respBody))
			panic(err)
		}

	} else {
		fmt.Println("getBotResponse() Request error")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}
