package models

type Action struct {
	Name      string                 `json:"name"`
	Groups    []string               `json:"groups"`
	Script    string                 `json:"script"`
	Tags      []string               `json:"tags"`
	Triggers  []string               `json:"triggers"`
	Hosts     []Service              `json:"hosts"`
	Variables map[string]interface{} `json:"variables"`
}

type ActionExtendedv3 struct {
	Name      string
	Groups    []string `json:"groups"`
	Script    string
	Triggers  []ActionExtendedv3
	Tags      []string
	Hosts     []Connection
	Variables map[string]interface{}
}
