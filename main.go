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
	Destination string `json:"destination"`
	Events      []struct {
		ReplyToken string `json:"replyToken"`
		Type       string `json:"type"`
		Timestamp  int64  `json:"timestamp"`
		Source     struct {
			Type   string `json:"type"`
			UserID string `json:"userId"`
		} `json:"source"`
		Message struct {
			ID   string `json:"id"`
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"events"`
}

type ReplyMessage struct {
	ReplyToken string `json:"replyToken"`
	Messages   []Text `json:"messages"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Profile struct {
	UserID        string `json:"userId"`
	DisplayName   string `่json:"displayName"`
	PictureURL    string `่json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
}

var ChannelToken = "l96rN/xy/3oDAZazD3E0xfv6VkJBjJugL9UEmIosZ24e1BhpCEwyLVX0R/O1QEsR+A6jbYx7DPrC66/BB5Ue/JKVzRwukjuAA0v+XirzOtwViD7CMIktROrk8Pa/2oVkjtabU3J5uyWJZUh3NSZ54gdB04t89/1O/w1cDnyilFU="

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

		fullname := getProfile(Line.Events[0].Source.UserID)
		log.Printf("%+v \n" , Line.Events)

		var text Text
		var test string = Line.Events[0].Message.Text

		if test == "Hi" {

			text = Text{
				Type: "text",
				Text: "ยินดีต้อนรับคุณ " + fullname,
			}
			log.Println(text)

		} else if test == "สวัสดี" {
			text = Text{
				Type: "text",

				
				Text: "สวัสดีคุณ " + fullname,
			}
			log.Println(text)
		} else {
			text = Text{
				Type: "text",
				Text: Line.Events[0].Message.Text + "? ขอโทษค่ะ ไม่มีคำสั่งนี้ในระบบ กรุณาตรวจสอบคำสั่งที่ถูกต้องอีกครั้งค่ะ",
			}
			log.Println(text)
		}

		message := ReplyMessage{
			ReplyToken: Line.Events[0].ReplyToken,
			Messages: []Text{
				text,
			},
		}

		log.Println(text)

		replyMessageLine(message)
		log.Println("%% message success")
		return c.String(http.StatusOK, "ok")

	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", appPort)))

	// e.Logger.Fatal(e.Start(":8080"))
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

func getProfile(userId string) string {

	url := "https://api.line.me/v2/bot/profile/" + userId

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannelToken)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var profile Profile
	if err := json.Unmarshal(body, &profile); err != nil {
		log.Println("%% err \n")

	}

	log.Println(profile.DisplayName)
	return profile.DisplayName

}
