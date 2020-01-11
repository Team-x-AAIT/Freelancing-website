package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"time"

	"github.com/Team-x-AAIT/Freelancing-website/api/entity"
)
// cookie struct
type cookie struct {
	Key        string
	Expiration time.Time
}
// struct contains entity job
type JobCollection struct {
	Jobs []entity.Job `json:"Content"`
}

var loggedIn = make([]cookie, 10)

const baseURL string = "http://localhost:8181/v1/"
// makes post request to the specific url inside the function
func PostUser(user *entity.User) error {
	requestbody, err := json.MarshalIndent(user, "", "\n")
	URL := fmt.Sprintf("%s%s", baseURL, "users")
	if err != nil {
		fmt.Println(err)
		return err
	}
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(requestbody))
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(body))
	return nil
}
// gets user
func GetUser(user *entity.User) (*entity.User, error) {
	URL := fmt.Sprintf("%s%s", baseURL, "user")
	formval := url.Values{}
	formval.Add("email", user.Email)
	formval.Add("password", user.Password)
	resp, err := http.PostForm(URL, formval)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	respjson := struct {
		Status  string
		Content entity.User
	}{}
	err = json.Unmarshal(body, &respjson)
	fmt.Println(respjson)
	if respjson.Status == "error" {
		return nil, errors.New("error")
	}
	return &respjson.Content, nil
}
// get searched job
func GetJobs(search string) ([]entity.Job, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%sjobs?category=%s", baseURL, search)
	req, _ := http.NewRequest("GET", URL, nil)
	fmt.Println(req)
	res, err := client.Do(req)
	//res, err := client.Get(URL)
	if err != nil {
		return nil, err
	}

	jobdata := &JobCollection{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, jobdata)
	if err != nil {
		return nil, err
	}
	return jobdata.Jobs, nil
}
// posts job
func PostMyJob(id uint, myjob string) error {
	// myjob := entity.MyJob{}
	client := &http.Client{}
	// requestbody, err := json.MarshalIndent(myjob, "", "\n")

	URL := fmt.Sprintf("%smyjob/%d/%s", baseURL, id, myjob)
	req, _ := http.NewRequest("POST", URL, nil)
	fmt.Println(req)
	_, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(URL)
	return nil

}
// get remommended jobs
func GetRecommendedJobs(id uint) ([]entity.Job, error) {
	client := &http.Client{}
	URL := fmt.Sprintf("%suser/%d", baseURL, id)
	req, _ := http.NewRequest("GET", URL, nil)
	fmt.Println(req)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	jobdata := &JobCollection{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, jobdata)
	if err != nil {
		return nil, err
	}
	return jobdata.Jobs, nil
}

// func GetMyJobs(id uint) ([]entity.Job, error) {
// 	client := &http.Client{}
// 	URL := fmt.Sprintf("%smyjob?userid=%d", baseURL, id)
// 	fmt.Println(URL)
// 	req, _ := http.NewRequest("GET", URL, nil)
// 	fmt.Println(req)
// 	res, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	jobdata := &JobCollection{}
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = json.Unmarshal(body, jobdata)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(jobdata.Jobs)
// 	return jobdata.Jobs, nil

// }
func PostJob(job *entity.Job) error {
	requestbody, err := json.MarshalIndent(job, "", "\n")
	URL := fmt.Sprintf("%s%s", baseURL, "job")
	if err != nil {
		fmt.Println(err)
		return err
	}
	resp, err := http.Post(URL, "application/json", bytes.NewBuffer(requestbody))
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println(string(body))
	return nil
}


