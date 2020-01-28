package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Team-x-AAIT/Freelancing-website/entity"
)

var baseURL = "http://localhost:8181/v1"

// PostProject is a method that validates and adds a project to the system.
func PostProject(project *entity.Project, uid string) (string, error) {

	url := baseURL + fmt.Sprintf("/user/%s/projects", uid)
	jsonStr, err := json.MarshalIndent(project, "", "\t\t")
	if err != nil {
		return "", err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusNotFound {
		return "", errors.New("not found")
	}

	locationStr := res.Header.Get("Location")
	pidA := []string{}
	pid := ""

	for index := len(locationStr) - 1; index >= 0; index-- {
		if locationStr[index] == '/' {
			break
		}
		pidA = append(pidA, string(locationStr[index]))
	}

	for index := len(pidA) - 1; index >= 0; index-- {
		pid += pidA[index]
	}

	return pid, nil
}

// ViewProject is a method that returns a project with its owner information.
func ViewProject(pid string) (*entity.ProjectUserContainer, error) {

	url := baseURL + fmt.Sprintf("/user/uid/projects/%s", pid)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	projectUserHolder := &entity.ProjectUserContainer{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, projectUserHolder)
	if err != nil {
		return nil, err
	}
	return projectUserHolder, nil
}

// UpdateProject is a method that is used for updating a project profile.
func UpdateProject(uid string, project *entity.Project) error {

	client := &http.Client{}
	url := baseURL + fmt.Sprintf("/user/%s/projects/%s", uid, project.ID)
	jsonStr, err := json.MarshalIndent(project, "", "\t\t")
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusNotFound {
		return errors.New("not found")
	}

	return nil
}

// RemoveProjectInformation is a method that is used for removing project and its dependencies.
func RemoveProjectInformation(uid, pid string) (*entity.Project, error) {

	client := &http.Client{}
	url := baseURL + fmt.Sprintf("/user/%s/projects/%s", uid, pid)

	req, _ := http.NewRequest("DELETE", url, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	projectHolder := &entity.Project{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, projectHolder)
	if err != nil {
		return nil, err
	}
	return projectHolder, nil

}

// FindProject is a method that is used for searching a projects using search key and other filters.
func FindProject(searchKey, searchBy, filterType, filterValue1, filterValue2, pageNumber string) ([]*entity.Project, error) {

	client := &http.Client{}
	url := baseURL + "/user/%s/projects/search"

	search := entity.SearchBag{}
	search.SearchKey = searchKey
	search.SearchBy = searchBy
	search.FilterTypeS = filterType
	search.FilterValue1S = filterValue1
	search.FilterValue2S = filterValue2
	search.PageNumS = pageNumber

	jsonStr, err := json.MarshalIndent(search, "", "\t\t")
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	projectHolder := []*entity.Project{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &projectHolder)
	if err != nil {
		return nil, err
	}
	return projectHolder, nil
}

// GetUserProjects is a method that returns all the projects linked to a user.
func GetUserProjects(uid string) ([]entity.Project, error) {
	url := baseURL + fmt.Sprintf("/user/%s/projects", uid)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	projectHolder := []entity.Project{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &projectHolder)
	if err != nil {
		return nil, err
	}
	return projectHolder, nil
}

// GetProjectApplicants is a function that sends an api request for getting all the applicants of a project.
func GetProjectApplicants(pid string) ([]*entity.Applicants, error) {

	url := baseURL + fmt.Sprintf("/user/uid/projects/%s/applications", pid)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	applicantsHolder := []*entity.Applicants{}
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println("--------------" + string(body))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &applicantsHolder)
	if err != nil {
		return nil, err
	}
	return applicantsHolder, nil
}

// GetApplication is a function that sends an api request for getting a certain application created by a user.
func GetApplication(pid, applicantUID string) (*entity.ApplicationBag, error) {

	url := baseURL + fmt.Sprintf("/user/%s/projects/%s/application", applicantUID, pid)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	application := &entity.ApplicationBag{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, application)
	if err != nil {
		return nil, err
	}
	return application, nil
}

// Apply is a method that enables project application process.
func Apply(application *entity.ApplicationBag) (string, error) {

	url := baseURL + fmt.Sprintf("/user/uid/projects/%s/applications", application.PID)
	jsonStr, err := json.MarshalIndent(application, "", "\t\t")
	if err != nil {
		return "", err
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}

	if res.StatusCode == http.StatusNotFound {
		return "", errors.New("not found")
	}

	locationStr := res.Header.Get("Location")
	return locationStr, nil
}

// HireApplicant is a function that sends an api request to the server for hiring an applicant.
func HireApplicant(pid, applicantUID string) error {

	client := &http.Client{}
	url := baseURL + fmt.Sprintf("/user/projects/%s/applications/%s", pid, applicantUID)

	req, _ := http.NewRequest("PUT", url, nil)
	_, err := client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

// RemoveApplicant is a method that removes or detaches an applicant from project.
func RemoveApplicant(applicantUID, pid string) (*entity.ApplicationBag, error) {

	client := &http.Client{}
	url := baseURL + fmt.Sprintf("/user/%s/projects/%s/applications", applicantUID, pid)

	req, _ := http.NewRequest("DELETE", url, nil)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	}

	applicationHolder := &entity.ApplicationBag{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, applicationHolder)
	if err != nil {
		return nil, err
	}
	return applicationHolder, nil

}
