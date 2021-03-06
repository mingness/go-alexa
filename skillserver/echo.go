package skillserver

import (
	"encoding/json"
	"errors"
	"log"
	"time"
)

// Request Functions
func (this *EchoRequest) VerifyTimestamp() bool {
	reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", this.Request.Timestamp)
	if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
		return true
	}

	return false
}

func (this *EchoRequest) VerifyAppID(myAppID string) bool {
	if this.Session.Application.ApplicationID == myAppID {
		return true
	}

	return false
}

func (this *EchoRequest) GetSessionID() string {
	return this.Session.SessionID
}

func (this *EchoRequest) GetUserID() string {
	return this.Session.User.UserID
}

func (this *EchoRequest) GetRequestType() string {
	return this.Request.Type
}

func (this *EchoRequest) GetIntentName() string {
	if this.GetRequestType() == "IntentRequest" {
		return this.Request.Intent.Name
	}

	return this.GetRequestType()
}

func (this *EchoRequest) GetSlotValue(slotName string) (string, error) {
	if _, ok := this.Request.Intent.Slots[slotName]; ok {
		return this.Request.Intent.Slots[slotName].Value, nil
	}

	return "", errors.New("Slot name not found.")
}

func (this *EchoRequest) AllSlots() map[string]EchoSlot {
	return this.Request.Intent.Slots
}

// Response Functions
func NewEchoResponse() *EchoResponse {
	er := &EchoResponse{
		Version: "1.0",
		Response: EchoRespBody{
			ShouldEndSession: true,
		},
	}

	return er
}

func (this *EchoResponse) OutputSpeech(text string) *EchoResponse {
	this.Response.OutputSpeech = &EchoRespPayload{
		Type: "PlainText",
		Text: text,
	}

	return this
}

func (this *EchoResponse) OutputSpeechSSML(text string) *EchoResponse {
	this.Response.OutputSpeech = &EchoRespPayload{
		Type: "SSML",
		SSML: text,
	}

	return this
}

func (this *EchoResponse) Card(title string, content string) *EchoResponse {
	return this.SimpleCard(title, content)
}

func (this *EchoResponse) SimpleCard(title string, content string) *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type:    "Simple",
		Title:   title,
		Content: content,
	}

	return this
}

func (this *EchoResponse) StandardCard(title string, content string, smallImg string, largeImg string) *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type:    "Standard",
		Title:   title,
		Content: content,
	}

	if smallImg != "" {
		this.Response.Card.Image.SmallImageURL = smallImg
	}

	if largeImg != "" {
		this.Response.Card.Image.LargeImageURL = largeImg
	}

	return this
}

func (this *EchoResponse) LinkAccountCard() *EchoResponse {
	this.Response.Card = &EchoRespPayload{
		Type: "LinkAccount",
	}

	return this
}

func (this *EchoResponse) Reprompt(text string) *EchoResponse {
	this.Response.Reprompt = &EchoReprompt{
		OutputSpeech: EchoRespPayload{
			Type: "PlainText",
			Text: text,
		},
	}

	return this
}

// AudioPlayerPlay Sends Alexa a command to stream the audio file identified by the specified audioItem.
// For more information see https://developer.amazon.com/blogs/post/Tx1DSINBM8LUNHY/new-alexa-skills-kit-ask-feature-audio-streaming-in-alexa-skills
func (this *EchoResponse) AudioPlayerPlay(url, token string) *EchoResponse {
	this.Response.Directives = &[]EchoRespAudioDir{
		EchoRespAudioDir{
			Type:         "AudioPlayer.Play",
			PlayBehavior: "REPLACE_ALL",
			AudioItem: EchoAudioItem{
				Stream: EchoAudioStream{
					Token:         token,
					URL:           url,
					OffsetInMilli: 0,
				},
			},
		},
	}

	mrs, err := json.Marshal(this.Response)
	log.Println(err)
	log.Println(string(mrs))

	return this
}

// AudioPlayerStop Stops any currently playing audio stream.
// For more information see https://developer.amazon.com/blogs/post/Tx1DSINBM8LUNHY/new-alexa-skills-kit-ask-feature-audio-streaming-in-alexa-skills
func (this *EchoResponse) AudioPlayerStop(url, token string) *EchoResponse {
	this.Response.Directives = &[]EchoRespAudioDir{
		EchoRespAudioDir{
			Type: "AudioPlayer.Stop",
		},
	}

	mrs, err := json.Marshal(this.Response)
	log.Println(err)
	log.Println(string(mrs))

	return this
}

// AudioPlayerClearQueue Clears the queue of all audio streams.
// For more information see https://developer.amazon.com/blogs/post/Tx1DSINBM8LUNHY/new-alexa-skills-kit-ask-feature-audio-streaming-in-alexa-skills
func (this *EchoResponse) AudioPlayClearQueue(url, token string) *EchoResponse {
	this.Response.Directives = &[]EchoRespAudioDir{
		EchoRespAudioDir{
			Type: "AudioPlayer.ClearQueue",
		},
	}

	mrs, err := json.Marshal(this.Response)
	log.Println(err)
	log.Println(string(mrs))

	return this
}

func (this *EchoResponse) EndSession(flag bool) *EchoResponse {
	this.Response.ShouldEndSession = flag

	return this
}

func (this *EchoResponse) String() ([]byte, error) {
	jsonStr, err := json.Marshal(this)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

// Request Types

type EchoRequest struct {
	Version string      `json:"version"`
	Session EchoSession `json:"session"`
	Request EchoReqBody `json:"request"`
}

type EchoSession struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes struct {
		String map[string]interface{} `json:"string"`
	} `json:"attributes"`
	User struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

type EchoReqBody struct {
	Type      string     `json:"type"`
	RequestID string     `json:"requestId"`
	Timestamp string     `json:"timestamp"`
	Intent    EchoIntent `json:"intent,omitempty"`
	Reason    string     `json:"reason,omitempty"`
}

type EchoIntent struct {
	Name  string              `json:"name"`
	Slots map[string]EchoSlot `json:"slots"`
}

type EchoSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Response Types

type EchoResponse struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Response          EchoRespBody           `json:"response"`
}

type EchoRespBody struct {
	OutputSpeech     *EchoRespPayload    `json:"outputSpeech,omitempty"`
	Card             *EchoRespPayload    `json:"card,omitempty"`
	Reprompt         *EchoReprompt       `json:"reprompt,omitempty"` // Pointer so it's dropped if empty in JSON response.
	Directives       *[]EchoRespAudioDir `json:"directives,omitempty"`
	ShouldEndSession bool                `json:"shouldEndSession"`
}

type EchoReprompt struct {
	OutputSpeech EchoRespPayload `json:"outputSpeech,omitempty"`
}

type EchoRespImage struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type EchoRespPayload struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Text    string        `json:"text,omitempty"`
	SSML    string        `json:"ssml,omitempty"`
	Content string        `json:"content,omitempty"`
	Image   EchoRespImage `json:"image,omitempty"`
}

type EchoRespAudioDir struct {
	Type         string        `json:"type"`
	PlayBehavior string        `json:"playBehavior,omitempty"`
	AudioItem    EchoAudioItem `json:"audioItem,omitempty"`
}

type EchoAudioItem struct {
	Stream EchoAudioStream `json:"stream"`
}

type EchoAudioStream struct {
	Token         string `json:"token"`
	URL           string `json:"url"`
	OffsetInMilli int    `json:"offsetInMilliseconds"`
}
