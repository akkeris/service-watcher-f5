package structs

type Spacespec struct {
	Name     string `json:"name"}`
	Internal bool   `json:"internal"}`
}

type Rulespec struct {
	Name         string `json:"name"`
	Partition    string `json:"partition"`
	ApiAnonymous string `json:"apiAnonymous"`
}

type Memberspec struct {
	Name string `json:"name"`
}

type Poolspec struct {
	Name      string       `json:"name"`
	Partition string       `json:"partition"`
	Monitor   string       `json:"monitor"`
	Members   []Memberspec `json:"members"`
}

type Virtualspec struct {
	Rules []string `json:"rules"`
}
