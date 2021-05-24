// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/lordronz/ronz-vulp/nhentai"
	"github.com/lordronz/ronz-vulp/amogus"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		go func (event *linebot.Event) {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					log.Println(message.Text)
					msgArr := strings.Fields(strings.ToLower(message.Text))
					if len(msgArr) == 1 && msgArr[0] == "nhentai" {
						nhentai.NhentaiRandom(bot, event)
					} else if msgArr[0] == "nhentai" {
						queryArr := msgArr[1:]
						query := strings.Join(queryArr, "%20")
						nhentai.NhentaiSearch(bot, event, query)
					} else if strings.HasPrefix(strings.ToLower(message.Text), "g/") {
						nhCode := strings.Split(message.Text, "/")[1]
						if nhCodeInt, err := strconv.Atoi(nhCode); err != nil {
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("Invalid nhentai code")).Do(); err != nil {
								log.Print(err)
							}
						} else {
							nhentai.NhentaiSearchById(bot, event, strconv.Itoa(nhCodeInt))
						}
					} else {
						cleanSus := amogus.RemoveDups(strings.ToLower(message.Text))
						if strings.Contains(cleanSus, "sus") {
							if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("AMOGUS")).Do(); err != nil {
								log.Print(err)
							}
						}
					}
				}
			}
		}(event)
	}
}
