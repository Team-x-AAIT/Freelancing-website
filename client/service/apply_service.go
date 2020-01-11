package service

import (
	"errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
)
// reaponse struct
type response struct {
	Status  string
	Content interface{}
}
// send get request to the  url 
func GetJobTo(id uint) (*entity.Job, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%sjbyid?id=%d", baseURL, id)
	fmt.Println(URL)
	req, _ := http.NewRequest("GET", URL, nil)
	fmt.Println(req)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	jobdata := entity.Job{}
	fmt.Println(jobdata)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	respjson := struct {
		Status  string
		Content entity.Job
	}{}
	err = json.Unmarshal(body, &respjson)
	fmt.Println(respjson)
	if respjson.Status == "error" {
		return nil, errors.New("error")
	}
	return &respjson.Content, nil
}
