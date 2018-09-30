package yml

import (
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type ExerciseGroup struct {
	Name string
	Tutor string
	Day string
	Time string
	Room string
	Participants uint32
}

type ExerciseGroupImport struct {
	CourseId string `yaml:"courseId"`
	Groups []ExerciseGroup `yaml:"groups"`
}

func ReadGroup(reader io.Reader) (ExerciseGroupImport, error)  {
	var groupImport ExerciseGroupImport

	// Read all bytes from the provided Reader
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return groupImport, err
	}

	// Unmarshal all bytes using the yaml package
	err = yaml.Unmarshal(data, &groupImport)
	if err != nil {
		return groupImport, err
	}

	return groupImport, nil
}
