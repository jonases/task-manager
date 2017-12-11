package models

// Context defines the structure used to send data to the HTML templates
type Context struct {
	Title   string
	Section string
	UserDB
}

type UserDB struct {
	FirstName     string `json:"fname,omitempty"`
	LastName      string `json:"lname,omitempty"`
	Email         string `json:"email,omitempty"`
	Password      string `json:"password,omitempty"`
	AccountActive bool   `json:"account_active,omitempty"`
}

type AlldocsResult struct {
	TotalRows int `json:"total_rows"`
	Offset    int
	Rows      []map[string]interface{}
}

type MsgData struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Message string `json:"msg"`
}

// Query defines the cloudant query structure
type Query struct {
	Selector map[string]interface{} `json:"selector"`
	Fields   []string               `json:"fields,omitempty"`
	Sort     []interface{}          `json:"sort,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	Skip     int                    `json:"skip,omitempty"`
}

type MsgDocument struct {
	Rev string `json:"rev"`
	ID  string `json:"id"`
}

type UsersDocument struct {
	Rev string `json:"rev"`
	ID  string `json:"id"`
}
