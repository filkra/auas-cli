package api

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"github.com/filkra/auas-cli/yml"
	"log"
	"net/url"
	"strconv"
	"strings"
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
	CourseShowURL			string = "course/show"
	GroupCreateURL 			string = "exercisegroup/create"
	GroupEditURL			string = "exercisegroup/edit"
	GroupDeleteURL			string = "exercisegroup/delete"
	EditModeSave 			string = "save"
	EditModeDelete			string = "delete"
)

const (
	GroupTableTitle			string = "Übungsgruppen"
)

var weekValues = map[string]string{
	"Montag": 		"Mo",
	"Dienstag": 	"Tu",
	"Mittwoch": 	"We",
	"Donnerstag": 	"Th",
	"Freitag": 		"Fr",
	"Samstag": 		"Sa",
	"Sonntag": 		"Su",
}

type CourseInformation struct {
	Name   string
	FormId uint32
	Tutors Tutors
}

type GroupInformation struct {
	Name         string
	Room         string
	Day          string
	Time         string
	Tutor        string
	Participants string
	Space        string
	Id           string
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
			FormParamCourse : {fmt.Sprint(info.FormId)},
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

func (c *Client) UpdateGroups(groupImport yml.ExerciseGroupImport) error {
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

	groupInfo, err := c.GetGroups(groupImport.CourseId)
	if err != nil {
		return err
	}

	groupMapping := createGroupMapping(groupInfo)

	// Send a POST request for each group to create it
	for _, group := range groupImport.Groups {
		// Create the request URL
		u, err := c.BaseURL.Parse(GroupEditURL + "/" + groupMapping[group.Name])
		if err != nil {
			return err
		}

		_, err = c.httpClient.PostForm(u.String(), url.Values {
			FormParamGroupName : {group.Name},
			FormParamTutor : {fmt.Sprint(info.Tutors[group.Tutor])},
			FormParamCourse : {fmt.Sprint(info.FormId)},
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

func (c *Client) WriteGroups(courseId string, groupInfo []GroupInformation) error {
	// Retrieve course information
	info, err := c.GetCourseInformation(courseId)
	if err != nil {
		log.Fatal(err)
	}

	groupMapping := createGroupMapping(groupInfo)

	// Send a POST request for each group to create it
	for _, group := range groupInfo {
		// Create the request URL
		u, err := c.BaseURL.Parse(GroupEditURL + "/" + groupMapping[group.Name])
		if err != nil {
			return err
		}

		_, err = c.httpClient.PostForm(u.String(), url.Values {
			FormParamGroupName : {group.Name},
			FormParamTutor : {fmt.Sprint(info.Tutors[group.Tutor])},
			FormParamCourse : {fmt.Sprint(info.FormId)},
			FormParamDay : {weekValues[group.Day]},
			FormParamTime : {group.Time},
			FormParamRoom : {group.Room},
			FormParamParticipants : {group.Space},
			FormParamEditMode : {EditModeSave},
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GetGroups(courseId string) ([]GroupInformation, error) {
	// Create the request URL
	u, err := c.BaseURL.Parse(CourseShowURL + "/" + courseId)
	if err != nil {
		return nil, err
	}

	// Retrieve the HTML source
	resp, err := soup.GetWithClient(u.String(), c.httpClient)
	if err != nil {
		return nil, err
	}

	// Parse the HTML source
	doc := soup.HTMLParse(resp)

	table := findGroupTable(&doc)
	if table == nil {
		return nil, nil
	}

	groups := readGroups(table)

	return groups, nil
}

func (c *Client) DeleteGroups(groupIds []GroupInformation) error {
	for _, group := range groupIds {
		// Create the request URL
		u, err := c.BaseURL.Parse(GroupDeleteURL + "/" + group.Id)
		if err != nil {
			return err
		}

		// Send a POST request to delete the specified group
		_, err = c.httpClient.PostForm(u.String(), url.Values {
			FormParamEditMode : {EditModeDelete},
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
		info.Tutors[strings.TrimSpace(option.Text())] = uint32(id)
	}

	// Find the course with its corresponding value
	options = doc.Find("select", "id", FormParamCourse).FindAll("option")
	option, err := strconv.ParseUint(options[1].Attrs()["value"], 10, 32)
	if err != nil {
		return info, err
	}

	// Store the course's name and value
	info.FormId = uint32(option)
	info.Name = strings.TrimSpace(options[1].Text())

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

func findGroupTable(root *soup.Root) *soup.Root {
	tables := root.FindAll("table", "class", "tablesorter")
	for _, table := range tables {
		if table.FindPrevSibling().Text() == GroupTableTitle {
			return &table
		}
	}
	return nil
}

func readGroups(table *soup.Root) []GroupInformation {
	var groups []GroupInformation
	rows := table.Find("tbody").FindAll("tr")
	for _, row := range rows {
		group := GroupInformation{
			Name: strings.TrimSpace(row.Children()[0].Text()),
			Room: strings.TrimSpace(row.Children()[1].Text()),
			Day: strings.TrimSpace(row.Children()[2].Text()),
			Time: strings.TrimSpace(row.Children()[3].Text()),
			Tutor: strings.TrimSpace(row.Children()[4].Text()),
			Participants: strings.TrimSpace(row.Children()[5].Text()),
			Space: strings.TrimSpace(row.Children()[6].Text()),
			Id: strings.TrimSpace(readGroupId(&row))}
		groups = append(groups, group)
	}
	return groups
}

func readGroupId(row *soup.Root) string {
	link := row.Find("a", "class", "icon-cup")
	href := link.Attrs()["href"]
	split := strings.Split(href, "/")
	return split[len(split) - 1]
}

func createGroupMapping(groupInfo []GroupInformation) map[string]string {
	mapping := map[string]string{}
	for _, info := range groupInfo {
		mapping[info.Name] = info.Id
	}
	return mapping
}
