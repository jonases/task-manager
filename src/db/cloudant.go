package db

import (
	"log"
	"models"
	"os"

	cloudant "github.com/IBM-Bluemix/go-cloudant"
	couchdb "github.com/timjacobi/go-couchdb"
)

var (
	Alldocs  models.AlldocsResult
	DB       *cloudant.DB
	Client   *cloudant.Client
	UsersDoc models.UsersDocument
	MsgsDoc  models.MsgDocument
)

// CloudantInit creates a client connection with the cloudant database
func CloudantInit() {

	user := os.Getenv("CLOUDANT_USER_NAME")
	password := os.Getenv("CLOUDANT_PASSWORD")

	var err error
	// create the client
	Client, err = cloudant.NewClient(user, password)
	if err != nil {
		log.Fatal(err)
	}

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
	if err := DB.Get(id, doc, nil); err != nil {
		log.Println(err)
	}
}

// Query queries the database, based on the field and value provided
func Query(field string, value string) []interface{} {
	query := cloudant.Query{}
	query.Selector = make(map[string]interface{})
	query.Selector[field] = value
	result, err := DB.SearchDocument(query)
	if err != nil {
		log.Println(err)
	}
	return result
}

// CreateDocument creates a document in the cloudant database
func CreateDocument(data interface{}) error {
	id, rev, err := DB.CreateDocument(data)
	if err != nil {
		log.Println(err)
		return err
	}

	if DB.Name() == "users" || DB.Name() == "users_test" {
		log.Println("Its users DB: ", DB.Name())
		UsersDoc.Rev = rev
		UsersDoc.ID = id
	} else {
		log.Println("Its NOT users DB: ", DB.Name())
		MsgsDoc.Rev = rev
		MsgsDoc.ID = id
	}

	return nil
}

// DeleteDocument deletes a doc from the db
func DeleteDocument(id, rev string) {
	_, err := DB.Delete(id, rev)
	if err != nil {
		log.Println(err)
	}
}
