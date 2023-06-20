package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey     string `yaml:"authtoken,omitempty"`
	ScheduleID string `yaml:"schedule_id,omitempty"`
}

func getConfig() Config {
	config := Config{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Unable to read ", homedir+"/.pd.yml")
	}
	if file, err := os.ReadFile(homedir + "/.pd.yml"); err == nil {
		if err := yaml.Unmarshal(file, &config); err != nil {
			log.Fatalln("Unable to parse ~/.pd.yml")
		}
	}

	if env := os.Getenv("PD_API_KEY"); env != "" {
		config.APIKey = env
	}
	if env := os.Getenv("PD_SCHEDULE_ID"); env != "" {
		config.ScheduleID = env
	}

	if config.APIKey == "" {
		log.Fatalln("authtoken is not present in ~/.pd.yml file. Set PD_API_KEY environment variable to override.")
	}
	if config.ScheduleID == "" {
		log.Fatalln("schedule_id is not present in ~/.pd.yml file. Set PD_SCHEDULE_ID environment variable to override.")
	}

	return config
}
