package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type fieldDetails struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type multipleChoiceAnswer struct {
	Type   string `json:"type"`
	Choice struct {
		Label string `json:"label"`
	} `json:"choice"`
	Field fieldDetails `json:"field"`
}

type textAnswer struct {
	Type  string       `json:"type"`
	Text  string       `json:"text"`
	Field fieldDetails `json:"field"`
}

type formField struct {
	Id                        string `json:"id"`
	Title                     string `json:"title"`
	Type                      string `json:"type"`
	Ref                       string `json:"ref"`
	Allow_multiple_selections bool   `json:"allow_multiple_selections"`
	Allow_other_choice        bool   `json:"allow_other_choice"`
}

type fromResponse struct {
	Form_id      string `json:"form_id"`
	Token        string `json:"token"`
	Submitted_at string `json:"submitted_at"`
	Landed_at    string `json:"landed_at"`
	Hidden       *struct {
		Target string `json:"target"`
	} `json:"hidden,omitempty"`
	Definition struct {
		Id      string        `json:"id"`
		Title   string        `json:"title"`
		Fields  []formField   `json:"fields"`
		Answers []interface{} `json:"answers"`
	} `json:"definition"`
}

type typeformResp struct {
	Event_id      string       `json:"event_id"`
	Event_type    string       `json:"event_type"`
	Form_response fromResponse `json:"form_response"`
}

var NewIncidenceReportMap map[string]string
var InternalIncidenceReportMap map[string]string

func init() {
	//forms fields id for new incidences
	NewIncidenceReportMap["qJrspeePNrLf"] = "method"
	NewIncidenceReportMap["qxvyhBkq13zP"] = "store_name"
	NewIncidenceReportMap["eYu8WCWbzE6I"] = "client_name"
	NewIncidenceReportMap["tu0PjF5FHsVz"] = "email"
	NewIncidenceReportMap["qybB8lcFUqI6"] = "phone"
	NewIncidenceReportMap["fpBdfsWHIMpG"] = "location"
	NewIncidenceReportMap["NKEuGeR5gN4w"] = "description"
	NewIncidenceReportMap["DuFU6b0JpvGs"] = "availability"

	//form fields id for internal incidences
	InternalIncidenceReportMap["cI0N9BlDnNUl"] = "installer"
	InternalIncidenceReportMap["mesBgqpX5Vo7"] = "order_number"
	InternalIncidenceReportMap["vXPVBa5dwOUm"] = "description"
	InternalIncidenceReportMap["lOEMisiljIby"] = "visual_ref1"
}

func digestTypeformAnswer(fieldAnswer []interface{}, new bool) (map[string]any, error) {

	valuesMap := make(map[string]any)
	var fieldsMap map[string]string

	if new {
		fieldsMap = NewIncidenceReportMap
	} else {
		fieldsMap = InternalIncidenceReportMap
	}

	for _, value := range fieldAnswer {
		v, ok := value.(multipleChoiceAnswer)
		if !ok {
			v, ok := value.(textAnswer)
			if !ok {
				return nil, errors.New("Not supported answer type")
			}
			valuesMap[fieldsMap[v.Field.Id]] = v
		}
		valuesMap[fieldsMap[v.Field.Id]] = v
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
