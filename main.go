package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/labstack/echo"
)
type LineMessage struct {
	Destination string 		`json:"destination"`
	Events      []struct {
		ReplyToken string 	`json:"replyToken"`
		Type       string	`json:"type"`
		Timestamp  int64  	`json:"timestamp"`
		Source     struct {
			Type   string 	`json:"type"`
			UserID string 	`json:"userId"`
		}`json:"source"`
		Message struct {
			ID   string 	`json:"id"`
			Type string 	`json:"type"`
			Text string 	`json:"text"`
		} `json:"message"`
	} `json:"events"`
}


type ReplyMessage struct {
	ReplyToken 	string `json:"replyToken"`
	Messages   	[]Text `json:"messages"`
}

type Text struct {
	Type 		string `json:"type"`
	Text 		string `json:"text"`
} 

var ChannelToken = "cjH2DqkZ3GJnIhdaJ1yviiUMLyd5oaslIVuoz8CwZumNDUStw+zzXKpClradFc7eox/zT7imst2SNOQ0krDWq58XnSd0vrAf1QHCYfS0KJ0HGUrY3bBhmFKhNBG4FcMM4fIqAzCQTZ+xgcBklhka6wdB04t89/1O/w1cDnyilFU="


func main() {
	appPort := os.Getenv("PORT")
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	e.POST("/webhook", func(c echo.Context) error {
		
		Line := new(LineMessage)
		if err := c.Bind(Line); err != nil {
			log.Println("err")
			return c.String(http.StatusOK, "error")
		}

		text := Text{
			Type : "text",
			Text : "ข้อความเข้ามา : " + Line.Events[0].Message.Text  + " ยินดีต้อนรับ : ",
		}
		
		message := ReplyMessage{
			ReplyToken : Line.Events[0].ReplyToken ,
			Messages : []Text{
				text,
			},
		}
		
		replyMessageLine(message)
		
		log.Println("%% message success")
		return c.String(http.StatusOK, "ok")
		
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", appPort)))
}


func replyMessageLine(Message ReplyMessage) error {
	value, _ := json.Marshal(Message)

	url := "https://api.line.me/v2/bot/message/reply"

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return err
}