package nhentai

type Nhentai struct {
	Id int `json:"id"`
	Media_id string `json:"media_id"`
	Title struct {
		English string `json:"english"`
		Japanese string `json:"japanese"`
		Pretty string `json:"pretty"`
	} `json:"title"`
	Images struct {
		Pages []struct {
			T string `json:"t"`
			W int `json:"w"`
			H int `json:"h"`
		} `json:"pages"`
		Cover struct {
			T string `json:"t"`
			W int `json:"w"`
			H int `json:"h"`
		} `json:"cover"`
		Thumbnail struct {
			T string `json:"t"`
			W int `json:"w"`
			H int `json:"h"`
		} `json:"thumbnail"`
	} `json:"images"`
	Scanlator string `json:"scanlator"`
	Upload_date int `json:"upload_date"`
	Num_pages int `json:"num_pages"`
	Num_favorites int `json:"num_favorites"`
}

type NhentaiGalleries struct {
	Result []Nhentai `json:"result"`
	Num_pages int `json:"num_pages"`
	Per_page int `json:"per_page"`
}

var NhentaiExtension map[string]string = map[string]string {
	"j": "jpg",
	"p": "png",
	"g": "gif",
}