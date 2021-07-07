package webmq

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"testing"
)

func TestWebuser_Webuserxml(t *testing.T) {

	wu := Webuser{
		Webroles: []Approle{
			{
				Name:   "role1",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Apiroles: []Approle{
			{
				Name:   "role2",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Variables: []Variable{
			{
				Name:  "httpsPort",
				Value: "9443",
			},
			{
				Name:  "httpHost",
				Value: "*",
			},
			{
				Name:  "mqRestCorsAllowedOrigings",
				Value: "*",
			},
		},
	}

	xml := wu.Webuserxml()
	fmt.Printf("%s\n", xml)
}

func TestWebuser_Webuserxml2(t *testing.T) {
	wu := Webuser{
		Webroles: []Approle{
			{
				Name:   "role1",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Apiroles: []Approle{
			{
				Name:   "role2",
				Users:  []string{"user1", "user2"},
				Groups: []string{"group1", "group2"},
			},
		},
		Variables: []Variable{
			{
				Name:  "httpsPort",
				Value: "9443",
			},
			{
				Name:  "httpHost",
				Value: "*",
			},
			{
				Name:  "mqRestCorsAllowedOrigins",
				Value: "*",
			},
		},
	}

	d, err := yaml.Marshal(&wu)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	fmt.Printf("--- m dump:\n%s\n\n", string(d))
}