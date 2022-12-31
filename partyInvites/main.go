package main

import (
	"fmt"
	"html/template"
  "net/http"
)

type Rsvp struct {
  Name string
  Email string
  Phone string
  WillAttend bool
};

var responses = make([]*Rsvp, 0, 10)
var templates = make(map[string]*template.Template, 3)

func loadTemplates() {
  templateNames := [5]string { "welcome", "form", "thanks", "list", "sorry" };
  for index, name := range templateNames {
    temp, err := template.ParseFiles("layout.html", name + ".html")
    if (err == nil) {
      templates[name] = temp
      fmt.Println("Loaded template", index, name)
    } else {
      panic(err)
    }
  }
}

func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
  templates["welcome"].Execute(writer, nil)
}

func listHandler(writer http.ResponseWriter, request *http.Request) {
  templates["list"].Execute(writer, responses)
}

type formData struct {
  *Rsvp
  Errors []string
}

func formHandler(writer http.ResponseWriter, request *http.Request) {
  if (request.Method == http.MethodGet) {
    templates["form"].Execute(writer, formData {
      Rsvp: &Rsvp{}, Errors: []string{},
    })
  } else if (request.Method == http.MethodPost) {
    request.ParseForm()
    responseData := Rsvp {
      Name: request.Form["name"][0],
      Email: request.Form["email"][0],
      Phone: request.Form["phone"][0],
      WillAttend: request.Form["willattend"][0] == "true",
    }

    errors := []string{};
    if (responseData.Name == "") {
      errors = append(errors, "Please enter you name")
    }
    if (responseData.Email == "") {
      errors = append(errors, "Please enter your email address")
    }
    if (responseData.Phone == "") {
      errors = append(errors, "Please enter your phone number")
    }

    if (len(errors) > 0) {
      templates["form"].Execute(writer, formData{
        Rsvp: &responseData, Errors: errors,
      })
    } else {
      responses = append(responses, &responseData)
      if (responseData.WillAttend) {
        templates["thanks"].Execute(writer, responseData.Name)
      } else {
        templates["sorry"].Execute(writer, responseData.Name)
      }
    }
  }
}

func main() {
  loadTemplates();

  http.HandleFunc("/", welcomeHandler)
  http.HandleFunc("/list", listHandler)
  http.HandleFunc("/form", formHandler)

  fmt.Println("Server Started on port 5000")
  err := http.ListenAndServe(":5000", nil)
  if (err != nil) {
    fmt.Println("Server error::", err)
  }
}
