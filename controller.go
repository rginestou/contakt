package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type indexData struct {
	Contact
	Title string
	ID    string
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	id := bson.NewObjectId()

	ids, ok := r.URL.Query()["id"]
	if ok && bson.IsObjectIdHex(ids[0]) {
		id = bson.ObjectIdHex(ids[0])
	}

	if r.Method == "POST" {
		r.ParseForm()

		var contact Contact
		contact.ID = id
		contact.Name = r.PostForm["name"][0]
		contact.Job = r.PostForm["job"][0]
		contact.Address = r.PostForm["address"][0]
		contact.Email = r.PostForm["email"][0]
		contact.Phone = r.PostForm["phone"][0]
		contact.Comment = r.PostForm["comment"][0]
		contact.CreatedAt = time.Now()

		updateContact(contact)

		http.Redirect(w, r, "/contact?id="+contact.ID.Hex(), http.StatusSeeOther)
		return
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

	tmpl.ExecuteTemplate(w, "contact", indexData{
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

	tmpl.ExecuteTemplate(w, "contact", indexData{
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

	tmpl.ExecuteTemplate(w, "contact", indexData{
		Title: "Contakt • Nouveau",
	})
}

func indexGET(w http.ResponseWriter, r *http.Request) {
	// Load templates
	tmpl, err := template.ParseFiles(
		"view/home.html",
		"view/head.html",
	)
	if err != nil {
		log.Println(err)
		// pageNotFoundHandler(w, r)
		return
	}

	tmpl.ExecuteTemplate(w, "home", indexData{
		Title: "Contakt",
	})
}
