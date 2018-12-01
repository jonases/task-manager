package shared

import (
	"log"
	"os"

	cloudant "github.com/IBM-Bluemix/go-cloudant"
	couchdb "github.com/timjacobi/go-couchdb"
)

var (
	Alldocs  AlldocsResult
	DB       *cloudant.DB
	Client   *cloudant.Client
	UsersDoc UsersDocument
)

// CloudantInit creates a client connection with the cloudant database
func CloudantInit() {

	user := os.Getenv("CLOUDANT_USER_NAME")
	password := os.Getenv("CLOUDANT_PASSWORD")
	var err error

	log.Println("connecting to cloudant")

	// create the client
	Client, err = cloudant.NewClient(user, password)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("successfully connected")

}

// CreateDBConnection creates a connection to a specific database
func CreateDBConnection(dbName string) {
	var err error
	DB, err = Client.EnsureDB(dbName)
	if err != nil {
		log.Println(err)
	}
}

// GetAllDocs retrieves all documents from the database
func GetAllDocs() error {

	options := make(couchdb.Options)
	// this options returns all data in the database document
	options["include_docs"] = true
	err := DB.AllDocs(&Alldocs, options)
	if err != nil {
		return err
	}
	return nil

}

// GetByID gets a document by its ID
func GetByID(id string, doc interface{}) {
	if err := DB.GetDocument(id, doc, nil); err != nil {
		log.Println(err)
	}
}

// QueryByFieldAndValue queries the database, based on the field and value provided
func QueryByFieldAndValue(field string, value string) ([]interface{}, error) {
	query := cloudant.Query{}
	query.Selector = make(map[string]interface{})
	query.Selector[field] = value
	result, err := DB.SearchDocument(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return result, nil
}

// CreateDocument creates a document in the cloudant database
func CreateDocument(data interface{}) error {
	id, rev, err := DB.CreateDocument(data)
	if err != nil {
		log.Println(err)
		return err
	}

	UsersDoc.Rev = rev
	UsersDoc.ID = id

	return nil
}

// InsertTodo inserts the todo in the cloudant database
func InsertTodo(data interface{}) (err error) {
	_, _, err = DB.CreateDocument(data)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

// DeleteDocument deletes a doc from the db
func DeleteDocument(id, rev string) error {
	_, err := DB.Delete(id, rev)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// UpdateDocument updates a doc in the db
func UpdateDocument(id, rev string, doc interface{}) error {
	_, err := DB.UpdateDocument(id, rev, doc)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
