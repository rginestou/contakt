package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
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
	Filter   string
	Contacts []Contact
}

func contactGET(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")
	if ids == "" || !bson.IsObjectIdHex(ids) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(ids)

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

func contactPOST(w http.ResponseWriter, r *http.Request) {
	id := bson.NewObjectId()
	ids := r.URL.Query().Get("id")
	if ids != "" && bson.IsObjectIdHex(ids) {
		id = bson.ObjectIdHex(ids)
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

	// Delete contact picture
	if r.FormValue("delete-avatar") == "on" {
		os.Remove("data/" + id.Hex() + ".png")
	}

	// Contact picture
	file, _, err := r.FormFile("avatar")
	if err == nil {
		f, _ := os.OpenFile("data/"+id.Hex()+".png", os.O_WRONLY|os.O_CREATE, 0666)
		io.Copy(f, file)

		f.Close()
		file.Close()
	}

	http.Redirect(w, r, "/contact?id="+contact.ID.Hex(), http.StatusSeeOther)
}

func editContactGET(w http.ResponseWriter, r *http.Request) {
	ids := r.URL.Query().Get("id")
	if ids == "" || !bson.IsObjectIdHex(ids) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(ids)

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
	ids := r.URL.Query().Get("id")
	if ids == "" || !bson.IsObjectIdHex(ids) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	id := bson.ObjectIdHex(ids)

	if err := deleteContactByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func indexGET(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")

	// Load templates
	tmpl, _ := template.ParseFiles(
		"view/home.html",
		"view/head.html",
	)

	// Get contacts
	contacts, _ := getAllContacts(filter)

	tmpl.ExecuteTemplate(w, "home", indexData{
		Title:    "Contakt",
		Filter:   filter,
		Contacts: contacts,
	})
}
