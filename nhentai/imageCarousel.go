package nhentai

type ImageCarousel struct {
	Type string `json:"type"`
	Alttext string `json:"altType"`
	Template struct {
		Type string `json:"type"`
		Columns []struct {
			ImageUrl string `json:"imageUrl"`
			Action struct {
				Type string `json:"type"`
				Label string `json:"label"`
			} `json:"action"`
		} `json:"columns"`
	} `json:"template"`
}
