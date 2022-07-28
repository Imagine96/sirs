package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type fieldDetails struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type multipleChoiceAnswer struct {
	Type   string `json:"type,omitempty"`
	Choice struct {
		Label string `json:"label,omitempty"`
	} `json:"choice,omitempty"`
	Field fieldDetails `json:"field,omitempty"`
}

type textAnswer struct {
	Type  string       `json:"type,omitempty"`
	Text  string       `json:"text,omitempty"`
	Field fieldDetails `json:"field,omitempty"`
}

type formField struct {
	Id                        string `json:"id,omitempty"`
	Title                     string `json:"title,omitempty"`
	Type                      string `json:"type,omitempty"`
	Ref                       string `json:"ref,omitempty"`
	Allow_multiple_selections bool   `json:"allow_multiple_selections,omitempty"`
	Allow_other_choice        bool   `json:"allow_other_choice,omitempty"`
}

type fromResponse struct {
	Form_id      string `json:"form_id,omitempty"`
	Token        string `json:"token,omitempty"`
	Submitted_at string `json:"submitted_at,omitempty"`
	Landed_at    string `json:"landed_at,omitempty"`
	Hidden       *struct {
		Target string `json:"target,omitempty"`
	} `json:"hidden,omitempty"`
	Definition struct {
		Id      string        `json:"id,omitempty"`
		Title   string        `json:"title,omitempty"`
		Fields  []formField   `json:"fields,omitempty"`
		Answers []interface{} `json:"answers,omitempty"`
	} `json:"definition,omitempty"`
}

type typeformResp struct {
	Event_id      string       `json:"event_id,omitempty"`
	Event_type    string       `json:"event_type,omitempty"`
	Form_response fromResponse `json:"form_response,omitempty"`
}

//TF -> typeform form
var tFNewIncidenceReportMap map[string]string
var tFInternalIncidenceReportMap map[string]string

//TB -> trello board
var TBIncidenceReportMap map[string]string

func init() {
	//forms fields id for new incidences - target custom field id
	tFNewIncidenceReportMap["qJrspeePNrLf"] = "62e21c03e837e24019703ea7" //client_type
	tFNewIncidenceReportMap["qxvyhBkq13zP"] = "store_name"
	tFNewIncidenceReportMap["eYu8WCWbzE6I"] = "62e21cbca9f76f7245154417" //client_name
	tFNewIncidenceReportMap["tu0PjF5FHsVz"] = "62e21cc8332faa093067d1ab" //client_email
	tFNewIncidenceReportMap["qybB8lcFUqI6"] = "62e21cd86edb5f1fb90674aa" //client_phone
	tFNewIncidenceReportMap["fpBdfsWHIMpG"] = "62e21ce608abc7502f0df2e3" //client_direction
	tFNewIncidenceReportMap["NKEuGeR5gN4w"] = "description"              //desc
	tFNewIncidenceReportMap["DuFU6b0JpvGs"] = "62e21cfa8d1e7f229bae09df" //client_availability

	//form fields id for internal incidences - target custom field id
	tFInternalIncidenceReportMap["cI0N9BlDnNUl"] = "62e21fc62c53508ebc065a02" //installer
	tFInternalIncidenceReportMap["mesBgqpX5Vo7"] = "62cfd6a1ca81118c3d84caac" //order_number
	tFInternalIncidenceReportMap["vXPVBa5dwOUm"] = "description"              //desc
}

func digestTypeformAnswers(fieldAnswer []interface{}, new bool) (map[string]string, error) {

	valuesMap := make(map[string]string)
	var fieldsMap map[string]string

	if new {
		fieldsMap = tFNewIncidenceReportMap
	} else {
		fieldsMap = tFInternalIncidenceReportMap
	}

	for _, value := range fieldAnswer {
		v, ok := value.(multipleChoiceAnswer)
		if !ok {
			v, ok := value.(textAnswer)
			if !ok {
				return nil, errors.New("Not supported answer type")
			}
			valuesMap[fieldsMap[v.Field.Id]] = v.Text
		}
		valuesMap[fieldsMap[v.Field.Id]] = v.Choice.Label
	}
	return valuesMap, nil
}

func HandleRequest(ctx context.Context, event typeformResp) (string, error) {
	err := publishTrelloCard(event)
	if err != nil {
		return "400", err
	}
	return "200", nil
}

func getEnvVar(key string) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		return "", errors.New("failed to load env")
	}

	if envValue, exist := os.LookupEnv(key); exist == false {
		return "", errors.New(fmt.Sprintf("var %v not found", key))
	} else {
		return envValue, nil
	}

}

func main() {
	/* lambda.Start(HandleRequest) */

}
