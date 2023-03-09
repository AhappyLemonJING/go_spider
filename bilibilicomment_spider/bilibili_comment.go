package bilibilicommentspider

type BilibiliComment struct {
	Data struct {
		Replies []struct {
			Content struct {
				Message string `json:"message"`
			} `json:"content"`
			Replies []struct {
				Content struct {
					Message string `json:"message"`
				} `json:"content"`
				Replies interface{} `json:"replies"`
			} `json:"replies"`
		} `json:"replies"`
	} `json:"data"`
}

type BiliComment struct {
	Comments string `json:"comments"`
}

func (table *BiliComment) TableName() string {
	return "bili_comment"
}
