package api

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/filkra/auas-cli/yml"
	"log"
	"net/url"
	"strconv"
)

const (
	FormParamGroupName 		string = "Groupname"
	FormParamTutor 			string = "Tutor"
	FormParamCourse 		string = "Course"
	FormParamDay 			string = "DayOfWeek"
	FormParamTime 			string = "Time"
	FormParamRoom 			string = "Room"
	FormParamParticipants 	string = "Maxparticipants"
	FormParamEditMode		string = "editmode"
)

const (
	GroupCreateURL 			string = "exercisegroup/create"
	EditModeSave 			string = "save"
)

type CourseInformation struct {
	CourseName string
	CourseFormId uint32
	Tutors Tutors
}

type Tutors map[string]uint32

func (c *Client) ImportGroups(groupImport yml.ExerciseGroupImport) error {
	// Retrieve course information
	info, err := c.GetCourseInformation(groupImport.CourseId)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the specified tutors exist
	err = validateTutors(groupImport.Groups, info.Tutors)
	if err != nil {
		return err
	}

	// Create the request URL
	u, err := c.BaseURL.Parse(GroupCreateURL + "/" + groupImport.CourseId)
	if err != nil {
		return err
	}

	// Send a POST request for each group to create it
	for _, group := range groupImport.Groups {
		_, err := c.httpClient.PostForm(u.String(), url.Values {
			FormParamGroupName : {group.Name},
			FormParamTutor : {fmt.Sprint(info.Tutors[group.Tutor])},
			FormParamCourse : {fmt.Sprint(info.CourseFormId)},
			FormParamDay : {group.Day},
			FormParamTime : {group.Time},
			FormParamRoom : {group.Room},
			FormParamParticipants : {fmt.Sprint(group.Participants)},
			FormParamEditMode : {EditModeSave},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GetCourseInformation(courseId string) (CourseInformation, error) {
	var info = CourseInformation{}
	info.Tutors = Tutors{}

	// Create the request URL
	u, err := c.BaseURL.Parse(GroupCreateURL + "/" + courseId)
	if err != nil {
		return info, err
	}

	// Retrieve the HTML source
	resp, err := soup.GetWithClient(u.String(), c.httpClient)
	if err != nil {
		return info, err
	}

	// Parse the HTML source
	doc := soup.HTMLParse(resp)

	// Find all tutors with their corresponding values and store them within a map
	options := doc.Find("select", "id", FormParamTutor).FindAll("option")
	for _, option := range options[1:] {
		id, _ := strconv.ParseUint(option.Attrs()["value"], 10, 32)
		info.Tutors[option.Text()] = uint32(id)
	}

	// Find the course with its corresponding value
	options = doc.Find("select", "id", FormParamCourse).FindAll("option")
	option, err := strconv.ParseUint(options[1].Attrs()["value"], 10, 32)
	if err != nil {
		return info, err
	}

	// Store the course's name and value
	info.CourseFormId = uint32(option)
	info.CourseName = options[1].Text()

	return info, nil
}

func validateTutors(groups []yml.ExerciseGroup, tutors Tutors) error {
	// Check if all specified tutors exist
	for _, group := range groups {
		if _, present := tutors[group.Tutor]; !present {
			return fmt.Errorf("tutor \"%s\" not found", group.Tutor)
		}
	}

	return nil
}
