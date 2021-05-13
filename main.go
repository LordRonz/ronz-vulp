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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/lordronz/ronz-vulp/nhentai"
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
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if strings.Contains(strings.ToLower(message.Text), "sus") {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("AMOGUS")).Do(); err != nil {
						log.Print(err)
					}
				} else if strings.ToLower(message.Text) == "nhentai" {
					client := &http.Client{}
					req, err := http.NewRequest("GET", "https://nhentai.net/random", nil)
					if err != nil {
						log.Print(err)
						return
					}
					req.Header.Add("User-Agent", "ronz-vulp")
					res, err := client.Do(req)
					if err != nil {
						log.Print(err)
						return
					}
					splittedUrl := strings.Split(res.Request.URL.String(), "/")
					randomUrl := fmt.Sprintf("https://nhentai.net/api/gallery/%s", splittedUrl[len(splittedUrl) - 2])

					req, err = http.NewRequest("GET", randomUrl, nil)
					if err != nil {
						log.Print(err)
						return
					}
					req.Header.Add("User-Agent", "ronz-vulp")
					res, err = client.Do(req)
					if err != nil {
						log.Print(err)
						return
					}
					if res.Body != nil {
						defer res.Body.Close()
					}
					body, err := ioutil.ReadAll(res.Body)
					if err != nil {
						log.Print(err)
						return
					}
					nhentaiRes := nhentai.Nhentai{}
					err = json.Unmarshal(body, &nhentaiRes)
					if err != nil {
						log.Printf("unable to parse value: %q, error: %s", string(body), err.Error())
						return
					}
					imgUrl := "https://i.nhentai.net/galleries/" + nhentaiRes.Media_id + "/" + "1." + nhentai.NhentaiExtension[nhentaiRes.Images.Pages[0].T]
					template := linebot.NewImageCarouselTemplate(
						linebot.NewImageCarouselColumn(
							imgUrl,
							linebot.NewURIAction("Go to LINE", "https://nhentai.net/g/" + strconv.Itoa(nhentaiRes.Id)),
						),
					)
					if _, err := bot.ReplyMessage(
						event.ReplyToken,
						linebot.NewTemplateMessage("Hentai for you", template),
					).Do(); err != nil {
						log.Print(err)
						return
					}
				}
			}
		}
	}
}
