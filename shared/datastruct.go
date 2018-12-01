package shared

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

// Query defines the cloudant query structure
type Query struct {
	Selector map[string]interface{} `json:"selector"`
	Fields   []string               `json:"fields,omitempty"`
	Sort     []interface{}          `json:"sort,omitempty"`
	Limit    int                    `json:"limit,omitempty"`
	Skip     int                    `json:"skip,omitempty"`
}

type UsersDocument struct {
	Rev string `json:"rev"`
	ID  string `json:"id"`
}

// Todo is the structure of the data containing in each task
type Todo struct {
	ID     string `json:"id"`
	Rev    string `json:"rev"`
	Title  string `json:"title"`
	Shared bool   `json:"shared"`
	State  string `json:"state"` // in-progress or done
	Email  string `json:"email"`
}

// TasksDocument is the structure of a row within a document
type TasksDocument struct {
	Rev string `json:"rev"`
	ID  string `json:"id"`
}
