package vab

type VAB struct {
	Id       int64    `json:"id"`
	Name     string   `json:"name"`
	ScriptId int64    `json:"scriptId"`
	URL      *string  `json:"url"`
	Music    []string `json:"music"`
}

type Audio struct {
	Id         int64   `json:"id"`
	VideoId    int64   `json:"videoId"`
	PageId     *int64  `json:"pageId"`
	Voice      *string `json:"voice"`
	VoiceStyle *string `json:"voiceStyle"`
	URL        *string `json:"url"`
}

type SlideShow struct {
	Id             int64    `json:"id"`
	VideoId        int64    `json:"videoId"`
	ScriptSplitsId int64    `json:"scriptSplitsId"`
	AudioDuration  *float64 `json:"audioDuration"`
}

type Video struct {
	Id      int64   `json:"id"`
	VideoId int64   `json:"videoId"`
	PageId  *int64  `json:"pageId"`
	URL     *string `json:"url"`
	AudioId *int64  `json:"audioId"`
}
