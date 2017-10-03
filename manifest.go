package main

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

/*
type URL struct {
	Host   string
	Domain string
}

func ParseURL(s, domain string) URL {
	p := strings.SplitN(s, ".", 2)
	if len(p) == 1 {
		p = append(p, domain)
	}
	return URL{
		Host:   p[0],
		Domain: p[1],
	}
}

func (u URL) String() string {
	if u.Domain == "" {
		return u.Host
	}
	return fmt.Sprintf("%s.%s", u.Host, u.Domain)
}

type User struct {
	Name     string `yaml:"username"`
	Password string `yaml:"password"`
}

type Organization struct {
	Users             map[string][]string `yaml:"users"`
	Domains           []string            `yaml:"domains"`
	Environment       map[string]string   `yaml:"env"`
	Spaces            map[string]*Space   `yaml:"spaces"`
	Quota             string              `yaml:"quota"`
    Quotas            map[string]*Quota   `yaml:"quotas"`
    SecurityGroupSets *SecurityGroupSet   `yaml:"security_group_sets"`
}

type Space struct {
	SSH                  string                 `yaml:"ssh"`
	Domain               string                 `yaml:"domain"`
	Users                map[string][]string    `yaml:"users"`
	Environment          map[string]string      `yaml:"env"`
	SharedServices       map[string]string      `yaml:"services"`
	Quota                string                 `yaml:"quota"`
	Applications         []*Application         `yaml:"apps"`
	UserProvidedServices []*UserProvidedService `yaml:"user-provided-services"`
    SecurityGroupSets    *SecurityGroupSet      `yaml:"security_group_sets"`
}

type Application struct {
	Name     string   `yaml:"name"`
	Hostname string   `yaml:"hostname"`
	Domain   string   `yaml:"domain"`
	URLs     []string `yaml:"urls"`

	Repository string `yaml:"repo"`
	Path       string `yaml:"path"`
	Image      string `yaml:"image"`
	Buildpack  string `yaml:"buildpack"`

	Memory      string            `yaml:"memory"`
	Disk        string            `yaml:"disk"`
	Instances   int               `yaml:"instances"`
	Environment map[string]string `yaml:"env"`

	BoundServices  map[string]string `yaml:"bind"`
	SharedServices []string          `yaml:"shared"`
}

type Quota struct {
	Memory                map[string]string `yaml:"memory"`
	TotalAppInstances     string            `yaml:"app-instances"`
	ServiceInstances      string            `yaml:"service-instances"`
	Routes                string            `yaml:"routes"`
	PaidPlans             bool              `yaml:"allow-paid-plans"`
	NumRoutesWithResPorts string            `yaml:"reserve-route-ports"`
}

type Manifest struct {
	Domains           []string                  `yaml:"domains"`
	Users             []User                    `yaml:"users"`
	Quotas            map[string]*Quota         `yaml:"quotas"`
	Organizations     map[string]*Organization  `yaml:"organizations"`
    SecurityGroups    map[string]*SecurityGroup `yaml:"security_groups"`
    SecurityGroupSets *SecurityGroupSet         `yaml:"security_group_sets"`
}

type UserProvidedService struct {
	Name            string      `yaml:"name"`
	Credentials     interface{} `yaml:"credentials"`
	RouteServiceUrl string      `yaml:"route_service_url"`
	SyslogDrainUrl  string      `yaml:"syslog_drain_url"`
}

type SecurityGroup struct {
	Rules             []interface{} `yaml:"rules"`
	SecurityGroupFile string        `yaml:"security_group_file"`
}

type SecurityGroupSet struct {
    Running []string `yaml:"running"`
    Staging []string `yaml:"staging"`
}
*/

/*
services:
  - name: my-cool-service
	plan: redis
	JSONParameters: "{Memory : "4Gb" "}""
*/

type Service struct {
	ServiceName    string `yaml:"instance_name"`
	Broker         string `yaml:"brokerName"`
	PlanName       string `yaml:"planName"`
	JSONParameters string `yaml:"JSONParameters"`
}
type Manifest struct {
	Services []Service `yaml:"create-services"`
}

func ParseManifest(src io.Reader) (Manifest, error) {
	var m Manifest
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return m, err
	}

	err = yaml.Unmarshal(b, &m)
	if err != nil {
		return m, err
	}

	return m, nil
}
