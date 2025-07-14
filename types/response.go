package types

type CreateResponse struct {
	Data    []string `json:"data"`
	Created int      `json:"created"`
}

type ReadResponse struct {
	Data  []map[string]interface{} `json:"data"`
	Total int64                    `json:"total"`
}

type UpdateResponse struct {
	Data    []string `json:"data"`
	Updated int64    `json:"updated"`
}

type DeleteResponse struct {
	Data    []string `json:"data"`
	Deleted int64    `json:"deleted"`
}
