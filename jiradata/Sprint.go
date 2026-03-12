package jiradata

type Sprint struct {
	ID            int    `json:"id,omitempty" yaml:"id,omitempty"`
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	State         string `json:"state,omitempty" yaml:"state,omitempty"`
	StartDate     string `json:"startDate,omitempty" yaml:"startDate,omitempty"`
	EndDate       string `json:"endDate,omitempty" yaml:"endDate,omitempty"`
	OriginBoardID int    `json:"originBoardId,omitempty" yaml:"originBoardId,omitempty"`
	Goal          string `json:"goal,omitempty" yaml:"goal,omitempty"`
}

type SprintList struct {
	Values []*Sprint `json:"values,omitempty" yaml:"values,omitempty"`
	IsLast bool      `json:"isLast,omitempty" yaml:"isLast,omitempty"`
}

type Board struct {
	ID   int    `json:"id,omitempty" yaml:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

type BoardList struct {
	Values []*Board `json:"values,omitempty" yaml:"values,omitempty"`
	IsLast bool     `json:"isLast,omitempty" yaml:"isLast,omitempty"`
}
