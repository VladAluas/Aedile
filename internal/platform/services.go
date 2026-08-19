package platform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/VladAluas/Aedile/internal/config"
	"gopkg.in/yaml.v3"
)

const templateInfo = "internal/generator/validators/"

func (p *Platform) DescribeService(ctx context.Context, service string) error {
	templatePath := filepath.Join(templateInfo, service+".yaml")

	info, err := GetTemplInfo(templatePath)
	if err != nil {
		return err
	}

	ShowTemplate(info)

	return nil
}

func GetTemplInfo(path string) (*config.TemplateInfo, error) {
	var tmplInfo config.TemplateInfo

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, &tmplInfo); err != nil {
		return nil, err
	}

	return &tmplInfo, nil
}

func ShowTemplate(info *config.TemplateInfo) {
	fmt.Printf("Service:      %s\n", info.Name)
	fmt.Printf("Description:  %s\n", info.Description)
	fmt.Println("Configuration: ")
	for name, prop := range info.Configuration {

		// Configuration Attribute
		fmt.Printf("  %s\n", name)

		// Required: bool
		fmt.Printf("    - required:       %v\n", prop.Required)

		// Type: string
		if prop.Type != "" {
			fmt.Printf("    - type:           %v\n", prop.Type)
		}

		// Default Values: any
		if prop.Default != nil {
			fmt.Printf("    - default values: %v\n", prop.Default)
		}

		// Description: string
		if prop.Description != "" {
			fmt.Printf("    - description:    %v\n", prop.Description)
		}

		// Accepted Values: list
		if prop.Values != nil {
			fmt.Printf("    - accepted values: %v\n", prop.Values)
		}
	}
}

func (p *Platform) ListServices(ctx context.Context) error {
	files, err := os.ReadDir(templateInfo)
	if err != nil {
		return err
	}
	fmt.Print("Available Services:\n")

	for _, file := range files {
		clFile := strings.Replace(file.Name(), ".yaml", "", 1)

		fmt.Printf("  - %s\n", clFile)
	}

	return nil
}

// Need to implement at a later date
//
// func (p *Platform) ValidateServices(ctx context.Context) error {
// 	cfg, err := config.Load()
// 	if err != nil {
// 		return err
// 	}
//
// 	files, err := os.ReadDir(templateInfo)
// 	if err != nil {
// 		return err
// 	}
//
// 	for _, file := range files {
// 		fmt.Print(cfg[file])
// 	}
//
// 	return nil
// }
