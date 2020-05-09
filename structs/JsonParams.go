package structs

type JsonParams struct {
	TplId  string      `json:"tpl"`
	Phone  string      `json:"phone"`
	Params interface{} `json:"params"`
}
