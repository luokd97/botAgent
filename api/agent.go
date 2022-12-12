package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// const robotApiUrl = "http://cors.redoc.ly/core/engine/dm/bot-response"
const robotApiUrl = "https://demo.laiye.com:8083/chatbot/v1alpha1/agents/{agentId}/channels/{channelId}/getReply"
const ipApiUrl = "http://ip-api.com/json"

// 无业务逻辑，仅测试服务部署环境网络状态
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

// @Summary		代理机器人Api
// @Description	入参出参均与原始文档保持一致 https://demo.laiye.com:8083/chatbot-openapi/swagger-ui/#/ChannelReplyService/ChannelReplyService_GetReply
// @Tags			botAgent
// @Accept			json
// @Produce		json
// @Success		200	body	string
// @Router			/agents/{agentId}/channels/{channelId}/getReply [post]
func GetBotResponse(c *gin.Context) {
	agentId := c.Param("agentId")
	channelId := c.Param("channelId")

	var url = robotApiUrl
	url = strings.Replace(url, "{agentId}", agentId, 1)
	url = strings.Replace(url, "{channelId}", channelId, 1)

	ByteBody, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(ByteBody))

	fmt.Println("发起请求 url=", url, "body=", string(ByteBody))
	resp, err := http.Post(url, "application/json", c.Request.Body)

	if err == nil {
		defer resp.Body.Close()
		respBody, _ := io.ReadAll(resp.Body)
		fmt.Println("resp status=", resp.Status)
		if resp.StatusCode != 200 {
			panic("resp.Status=" + resp.Status)
		}
		fmt.Println("getBotResponse() respBody=", string(respBody))
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
