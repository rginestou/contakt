package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type contactData struct {
	Contact
	Title string
	ID    string
}

type indexData struct {
	Title    string
	Contacts []Contact
}

func contactPOST(w http.ResponseWriter, r *http.Request) {
	id := bson.NewObjectId()
	ids, ok := r.URL.Query()["id"]
	if ok && bson.IsObjectIdHex(ids[0]) {
		id = bson.ObjectIdHex(ids[0])
	}

	r.ParseMultipartForm(32 << 20)

	var contact Contact
	contact.ID = id
	contact.Name = r.FormValue("name")
	contact.Job = r.FormValue("job")
	contact.Address = r.FormValue("address")
	contact.Email = r.FormValue("email")
	contact.Phone = r.FormValue("phone")
	contact.Comment = r.FormValue("comment")
	contact.CreatedAt = time.Now()

	updateContact(contact)

	// Contact picture
	file, handler, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, _ := os.OpenFile("data/"+id.Hex()+filepath.Ext(handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)

	http.Redirect(w, r, "/contact?id="+contact.ID.Hex(), http.StatusSeeOther)
}

func contactGET(w http.ResponseWriter, r *http.Request) {
	id := bson.NewObjectId()

	ids, ok := r.URL.Query()["id"]
	if ok && bson.IsObjectIdHex(ids[0]) {
		id = bson.ObjectIdHex(ids[0])
	}

	// Get contact
	contact, err := getContactByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/contact.html",
		"view/head.html",
	)

	tmpl.ExecuteTemplate(w, "contact", contactData{
		Contact: contact,
		Title:   "Contakt",
		ID:      id.Hex(),
	})
}

func editContactGET(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "No contact id provided", http.StatusBadRequest)
		return
	}
	id := bson.ObjectIdHex(ids[0])

	// Get contact
	contact, err := getContactByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/contact_edit.html",
		"view/head.html",
	)

	tmpl.ExecuteTemplate(w, "contact", contactData{
		Contact: contact,
		Title:   "Contakt • Editer",
		ID:      id.Hex(),
	})
}

func newContactGET(w http.ResponseWriter, r *http.Request) {
	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/contact_edit.html",
		"view/head.html",
	)

	tmpl.ExecuteTemplate(w, "contact", contactData{
		Title: "Contakt • Nouveau",
	})
}

func deleteContactPOST(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "No contact id provided", http.StatusBadRequest)
		return
	}
	id := bson.ObjectIdHex(ids[0])

	if err := deleteContactByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func indexGET(w http.ResponseWriter, r *http.Request) {
	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/home.html",
		"view/head.html",
	)

	// Get contacts
	contacts, _ := getAllContacts()

	tmpl.ExecuteTemplate(w, "home", indexData{
		Title:    "Contakt",
		Contacts: contacts,
	})
}
