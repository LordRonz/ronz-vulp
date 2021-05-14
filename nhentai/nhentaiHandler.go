package nhentai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"strconv"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

func NhentaiRandom(bot *linebot.Client, event *linebot.Event) {
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
	nhentaiRes := Nhentai{}
	err = json.Unmarshal(body, &nhentaiRes)
	if err != nil {
		log.Printf("unable to parse value: %q, error: %s", string(body), err.Error())
		return
	}
	var columns []*linebot.ImageCarouselColumn
	for i := 1; i < 6; i++ {
		imgUrl := "https://i.nhentai.net/galleries/" + nhentaiRes.Media_id + "/" + strconv.Itoa(i) + "." + NhentaiExtension[nhentaiRes.Images.Pages[i - 1].T]
		columns = append(
			columns,
			linebot.NewImageCarouselColumn(
				imgUrl,
				linebot.NewURIAction("g/" + strconv.Itoa(nhentaiRes.Id), "https://nhentai.net/g/" + strconv.Itoa(nhentaiRes.Id)),
			),
		)
	}
	template := linebot.NewImageCarouselTemplate(
		columns...,
	)
	if _, err := bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewTemplateMessage("Hentai for you", template),
	).Do(); err != nil {
		log.Print(err)
		return
	}
}
