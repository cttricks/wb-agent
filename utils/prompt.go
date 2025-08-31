package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Structs mapping the YAML schema
type Instructions struct {
	Rules      []string    `yaml:"rules"`
	Tables     []Table     `yaml:"tables"`
	Procedures []Procedure `yaml:"procedures"`
}

type Table struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Columns     []Column  `yaml:"columns"`
	Usage       []string  `yaml:"usage"`
	Examples    []Example `yaml:"examples"`
}

type Column struct {
	Name  string `yaml:"name"`
	Type  string `yaml:"type"`
	Notes string `yaml:"notes"`
}

type Example struct {
	Question string `yaml:"question"`
	SQL      string `yaml:"sql"`
}

type Procedure struct {
	Name        string    `yaml:"name"`
	Params      []Column  `yaml:"params"`
	Description string    `yaml:"description"`
	Examples    []Example `yaml:"examples"`
}

func SystemPrompt() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Check if cached Instruction.md already exists
	cacheFile := filepath.Join(dir, "Instruction.md") // singular
	if _, err := os.Stat(cacheFile); err == nil {
		data, err := os.ReadFile(cacheFile)
		if err != nil {
			return "", err
		}
		fmt.Println("[AI] Using cached Instruction.md")
		return string(data), nil
	}

	// Compose new from YAML
	fmt.Println("[AI] Composing instruction from Instructions.yaml")
	yamlFile := filepath.Join(dir, "Instruction.yaml")
	instr, err := loadInstructions(yamlFile)
	if err != nil {
		return "", err
	}

	systemPrompt := buildPrompt(instr)

	// Save into Instruction.md
	err = os.WriteFile(cacheFile, []byte(systemPrompt), 0644)
	if err != nil {
		fmt.Println("[ERROR] Failed to save composed Instruction.md")
	}

	return systemPrompt, nil
}

// load YAML into struct
func loadInstructions(filename string) (*Instructions, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var instr Instructions
	err = yaml.Unmarshal(data, &instr)
	if err != nil {
		return nil, err
	}
	return &instr, nil
}

// build system instruction text for OpenAI
func buildPrompt(instr *Instructions) string {
	output := "You are an assistant that converts natural language into MySQL queries.\n\n"

	// rules
	output += "## Rules:\n"
	for _, rule := range instr.Rules {
		output += "- " + rule + "\n"
	}
	output += "\n"

	// tables
	for _, t := range instr.Tables {
		output += fmt.Sprintf("### Table: %s\n", t.Name)
		output += t.Description + "\n"
		output += "Columns:\n"
		for _, c := range t.Columns {
			output += fmt.Sprintf("- %s (%s): %s\n", c.Name, c.Type, c.Notes)
		}
		output += "Usage notes:\n"
		for _, u := range t.Usage {
			output += "- " + u + "\n"
		}
		output += "Examples:\n"
		for _, ex := range t.Examples {
			output += fmt.Sprintf("- Q: %s\n  SQL:\n%s\n", ex.Question, ex.SQL)
		}
		output += "\n"
	}

	// procedures
	if len(instr.Procedures) > 0 {
		output += "## Stored Procedures\n"
		for _, p := range instr.Procedures {
			output += fmt.Sprintf("### %s\n", p.Name)
			output += p.Description + "\n"
			if len(p.Params) > 0 {
				output += "Parameters:\n"
				for _, param := range p.Params {
					output += fmt.Sprintf("- %s (%s): %s\n", param.Name, param.Type, param.Notes)
				}
			}
			for _, ex := range p.Examples {
				output += fmt.Sprintf("- Q: %s\n  SQL:\n%s\n", ex.Question, ex.SQL)
			}
			output += "\n"
		}
	}

	return output
}
