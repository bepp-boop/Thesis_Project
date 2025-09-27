package parser

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Workflow is a (simplified) representation of a GitHub Actions workflow.
type Workflow struct {
	Jobs map[string]WorkflowJob `yaml:"jobs"`
}

// WorkflowJob is a (simplified) representation of a workflow job.
type WorkflowJob struct {
	Name  string    `yaml:"name"`
	Steps []JobStep `yaml:"steps"`
}

// JobStep is a (simplified) representation of a workflow job step object.
type JobStep struct {
	With        map[string]string `yaml:"with"`
	Env         map[string]string `yaml:"env"`
	Name        string            `yaml:"name"`
	Run         string            `yaml:"run"`
	Uses        string            `yaml:"uses"`
	UsesComment string            `yaml:"-"`
}

func (step *JobStep) UnmarshalYAML(node *yaml.Node) error {
	for i := range node.Content {
		if i%2 == 1 {
			continue
		}

		key := node.Content[i].Value
		value := node.Content[i+1]

		var err error
		switch key {
		case "env":
			err = value.Decode(&step.Env)
		case "name":
			step.Name = value.Value
		case "run":
			step.Run = value.Value
		case "uses":
			step.Uses = value.Value
			step.UsesComment = strings.TrimLeft(value.LineComment, "# ")
		case "with":
			err = value.Decode(&step.With)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// ParseWorkflow parses a GitHub Actions workflow file into a Workflow struct.
func ParseWorkflow(data []byte) (Workflow, error) {
	var workflow Workflow
	if err := yaml.Unmarshal(data, &workflow); err != nil {
		return workflow, fmt.Errorf("could not parse workflow: %v", err)
	}

	return workflow, nil
}
