package types

type CreateRequest struct {
	Prefix string        `json:"prefix"`
	Data   []interface{} `json:"data"`
}

type ReadRequest struct {
	Prefix string                 `json:"prefix"`
	Filter map[string]interface{} `json:"filter,omitempty"`
	Sort   map[string]int         `json:"sort,omitempty"`
	Limit  int                    `json:"limit,omitempty"`
	Skip   int                    `json:"skip,omitempty"`
}

type UpdateRequest struct {
	Prefix string                 `json:"prefix"`
	Filter map[string]interface{} `json:"filter"`
	Data   interface{}            `json:"data"`
	Upsert bool                   `json:"upsert,omitempty"`
}

type DeleteRequest struct {
	Prefix string                 `json:"prefix"`
	Filter map[string]interface{} `json:"filter"`
}
