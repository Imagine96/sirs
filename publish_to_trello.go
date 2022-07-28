package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/guregu/null.v3"
)

type tAttachments struct {
	Id        string   `json:"id"`
	Bytes     string   `json:"bytes"`
	Date      string   `json:"date"`
	EdgeColor string   `json:"edgeColor"`
	IdMember  string   `json:"idMember"`
	IsUpload  bool     `json:"isUpload"`
	MimeType  string   `json:"mimeType"`
	Name      string   `json:"name"`
	Previews  []string `json:"previews"`
	Url       string   `json:"url"`
	Pos       string   `json:"pos"`
}

type tSticker struct {
	Id          string   `json:"id"`
	Top         int      `json:"top"`
	Left        float32  `json:"left"`
	ZIndex      int      `json:"zIndex"`
	Rotate      int      `json:"rotate"`
	Image       string   `json:"image"`
	ImageUrl    string   `json:"imageUrl"`
	ImageScaled []string `json:"imageScaled"`
}

type tCardCover struct {
	IdAttachment         null.String `json:"idAttachment"`
	Color                null.String `json:"color"`
	IdUploadedBackground null.String `json:"idUploadedBackground"`
	Size                 string      `json:"size"`
	Brightness           string      `json:"brightness"`
	IdPlugin             null.String `json:"idPlugin"`
}

type tCard struct {
	Attachments           []tAttachments `json:"attachments"`
	Stickers              []tSticker     `json:"stickers"`
	CheckItemStates       interface{}    `json:"checkItemStates"`
	Id                    string         `json:"id"`
	Badges                interface{}    `json:"badges"`
	Closed                bool           `json:"closed"`
	DueComplete           bool           `json:"dueComplete"`
	DateLastActivity      string         `json:"dateLastActivity"`
	Desc                  string         `json:"desc"`
	Due                   null.String    `json:"due"`
	DueReminder           null.String    `json:"dueReminder"`
	Email                 null.String    `json:"email"`
	IdBoard               string         `json:"idBoard"`
	IdChecklists          []string       `json:"idChecklists"`
	IdList                string         `json:"idList"`
	IdMembers             []string       `json:"idMembers"`
	IdMembersVoted        []string       `json:"idMembersVoted"`
	IdShort               int            `json:"idShort"`
	IdAttachmentCover     null.String    `json:"idAttachmentCover"`
	Labels                []interface{}  `json:"labels"`
	IdLabels              []string       `json:"idLabels"`
	ManualCoverAttachment bool           `json:"manualCoverAttachment"`
	Name                  string         `json:"name"`
	Pos                   int            `json:"pos"`
	ShortLink             string         `json:"shortLink"`
	ShortUrl              string         `json:"shortUrl"`
	Start                 interface{}    `json:"start"`
	Subscribed            bool           `json:"subscribed"`
	Url                   string         `json:"url"`
	IsTemplate            bool           `json:"isTemplate"`
	CardRole              null.String    `json:"cardRole"`
	limits                interface{}
	DescData              struct {
		Emoji interface{} `json:"emoji"`
	} `json:"descData"`
}

const (
	createCardBaseUrl     = "https://api.trello.com/1/cards?idList=%v&key=%v&token=%v&desc=%v"
	setOrderNumberBaseUrl = "https://api.trello.com/1/card/%v/customField/%v/item?&key=%v&token=%v" //cardId, customFieldId
)

func getCardStruct(r *http.Response) (*tCard, error) {
	if bodyBytes, err := ioutil.ReadAll(r.Body); err != nil {
		return nil, err
	} else {
		var body tCard
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			return nil, err
		}
		return &body, nil
	}
}

func publishTrelloCard(t typeformResp) error {

	/* client := &http.Client{Timeout: 10 * time.Second} */

	token, err := getEnvVar("TOKEN")
	if err != nil {
		return err
	}
	apiKey, err := getEnvVar("KEY")
	if err != nil {
		return err
	}
	newIncFormId, err := getEnvVar("NEW_INC_FORM_ID")
	if err != nil {
		return err
	}
	internalIncFormId, err := getEnvVar("INTERNAL_INC_FORM_ID")
	if err != nil {
		return err
	}

	switch t.Form_response.Form_id {
	case newIncFormId:
		{
			if listId, err := getEnvVar("NEW_INC_LIST_ID"); err != nil {
				return err
			} else {

				data, err := digestTypeformAnswers(t.Form_response.Definition.Answers, true)
				if err != nil {
					return err
				}

				resp, err := http.Post(fmt.Sprintf(createCardBaseUrl, listId, apiKey, token, data["description"]), "application/json", nil)
				if err != nil {
					return err
				}

				defer resp.Body.Close()

				card, err := getCardStruct(resp)
				if err != nil {
					return err
				}

				return NIRPopulateCustomFields(card.ShortUrl, t)
			}
		}
	case internalIncFormId:
		{
			if listId, err := getEnvVar("INTERNAL_INC_LIST_ID"); err != nil {
				return err
			} else {

				data, err := digestTypeformAnswers(t.Form_response.Definition.Answers, true)
				if err != nil {
					return err
				}

				resp, err := http.Post(fmt.Sprintf(createCardBaseUrl, listId, apiKey, token, data["description"]), "application/json", nil)
				if err != nil {
					return err
				}
				defer resp.Body.Close()

				card, err := getCardStruct(resp)

				if err != nil {
					return err
				}

				return IIRPopulateCustomFields(card.ShortUrl, t)
			}
		}
	default:
		{
			return errors.New("unregistered form id")
		}
	}
}

//IIR -> internal incidence report
func IIRPopulateCustomFields(cardShortcut string, t typeformResp) error {
	return nil
}

//NIR -> new incidence report
func NIRPopulateCustomFields(cardShortcut string, t typeformResp) error {
	return nil
}

//pcf -> Populate Custom Field
func pcfOrderNumber(cardId string, value string, token string, apiKey string, client *http.Client) error {

	orderNumberFieldId, err := getEnvVar("CFID_ORDER_NUMBER")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(setOrderNumberBaseUrl, cardId, orderNumberFieldId, apiKey, token), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if _, err := client.Do(req); err != nil {
		return err
	}

	return nil
}
