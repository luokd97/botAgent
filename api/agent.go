package api

import (
	"bytes"
	"compress/gzip"
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
	fmt.Println("resp=", resp)
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
// @Router			/chatbot/v1alpha1/agents/{agentId}/channels/{channelId}/getReply [post]
func GetBotResponse(c *gin.Context) {
	agentId := c.Param("agentId")
	channelId := c.Param("channelId")

	var url = robotApiUrl
	url = strings.Replace(url, "{agentId}", agentId, 1)
	url = strings.Replace(url, "{channelId}", channelId, 1)

	byteBody, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(byteBody))

	fmt.Println("发起请求 url=", url, "body=", string(byteBody))
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, c.Request.Body)
	req.Header = c.Request.Header
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "请求机器人api失败")
		panic(err)
	}
	if resp.StatusCode != 200 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "请求机器人api状态码异常")
		panic("resp.Status=" + resp.Status)
	}

	fmt.Println("resp=", resp)
	defer resp.Body.Close()
	respBodyBytes, _ := io.ReadAll(resp.Body)
	gzipReader, _ := gzip.NewReader(bytes.NewReader(respBodyBytes))
	defer gzipReader.Close()
	deBuffer := new(bytes.Buffer)
	_, err = io.Copy(deBuffer, gzipReader)
	deGzipBytes, _ := io.ReadAll(deBuffer)

	fmt.Println("getBotResponse() respBodyBytes string(deGzipBytes)=", string(deGzipBytes))
	var data map[string]interface{}
	err = json.Unmarshal(deGzipBytes, &data)
	if err == nil {
		intent := data["intent"].(map[string]interface{})
		id := fmt.Sprint(intent["id"])
		name := fmt.Sprint(intent["name"])
		go AddRecord(id, name)
		c.IndentedJSON(http.StatusOK, data)
	} else {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Intent解析失败")
		panic(err)
	}

}
