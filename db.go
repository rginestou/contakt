package main

import (
	"log"
	"time"

	"gopkg.in/mgo.v2/bson"

	mgo "gopkg.in/mgo.v2"
)

// Contact model
type Contact struct {
	ID        bson.ObjectId `bson:"_id" json:"_id"`
	Name      string        `bson:"name" json:"name"`
	Job       string        `bson:"job" json:"job"`
	Address   string        `bson:"address" json:"address"`
	Phone     string        `bson:"phone" json:"phone"`
	Email     string        `bson:"email" json:"email"`
	Comment   string        `bson:"comment" json:"comment"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
}

// ToID ...
func (c Contact) ToID() string {
	return c.ID.Hex()
}

var db *mgo.Database

// Connect to the database
func connect(server, database string) {
	session, err := mgo.Dial(server)
	if err != nil {
		log.Fatal(err)
	}

	db = session.DB(database)
}

func getAllContacts() ([]Contact, error) {
	var resp []Contact
	err := db.C("contacts").Find(bson.M{}).All(&resp)

	return resp, err
}

func getContactByID(ID bson.ObjectId) (Contact, error) {
	var resp Contact
	err := db.C("contacts").FindId(ID).One(&resp)

	return resp, err
}

func updateContact(contact Contact) error {
	_, err := db.C("contacts").Upsert(bson.M{"_id": contact.ID}, contact)

	return err
}

func deleteContactByID(ID bson.ObjectId) error {
	return db.C("contacts").RemoveId(ID)
}
