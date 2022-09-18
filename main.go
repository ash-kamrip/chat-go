package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.Write([]byte(`
		<head>
		<Title> Chat
		</Title>
		</head>
		<body>
		Let's Chat !
		</body>
		</html>
		`)) // <html>

	})

	//starting the web server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		// print error
		log.Fatal("ListenAndServe:", err)
	}
}
