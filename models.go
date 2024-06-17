package main

type GenRequest struct {
	UserPrompt   string `json:"user_prompt"`         // require
	AspectRatio  string `json:"aspect_ratio"`        // require
	ExpandPrompt bool   `json:"expand_prompt"`       // require
	ImageUrl     string `json:"image_url,omitempty"` //option, uploaded refer image url
}

type VideoTask struct {
	ID                  string      `json:"id"`
	Prompt              string      `json:"prompt"`
	State               string      `json:"state"` //"processing", "completed"
	CreatedAt           string      `json:"created_at"`
	Video               Video       `json:"video"`
	Liked               interface{} `json:"liked"`
	EstimateWaitSeconds interface{} `json:"estimate_wait_seconds"`
}

type Video struct {
	Url       string      `json:"url"`
	Width     int         `json:"width"`
	Height    int         `json:"height"`
	Thumbnail interface{} `json:"thumbnail"`
}
