package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	dapp_image "github.com/marcelstanley/millionerc20/image"
)

const (
	PORT            = ":3000"
	MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
)

type GlobalState struct {
	Count int
}

var global GlobalState
var sessionManager *scs.SessionManager

func getHandler(w http.ResponseWriter, r *http.Request) {
	// Display submit response, when available
	response := sessionManager.GetString(r.Context(), "response")
	component := page(response)
	component.Render(r.Context(), w)

	//TODO Display dapp-image retieved from dapp endpoint
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Update state
	// TODO Should we limit the form size? I don't think so
	//	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	// TODO Make sure we repond with an error if the image is too big
	// XXX This does not seemt to be working
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		updateStatus(r, "Max upload size of 1MB exceeded (%v)", http.StatusBadRequest)
		return
	}

	// Get file from request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		updateStatus(r, "Max upload size of 1MB exceeded (%v)", http.StatusBadRequest)
		return
	}
	defer file.Close()

	f, err := fileHeader.Open()
	if err != nil {
		status := fmt.Sprintf(err.Error() + " (%v)")
		updateStatus(r, status, http.StatusBadRequest)
		return
	}

	img, err := dapp_image.Decode(f)
	if err != nil {
		status := fmt.Sprintf(err.Error() + " (%v)")
		updateStatus(r, status, http.StatusBadRequest)
		return
	}
	//log.Printf("img.Rect: %v", img.Rect)

	// Translate image
	x, err := strconv.Atoi(r.FormValue("posX"))
	y, err := strconv.Atoi(r.FormValue("posY"))
	if err == nil {
		//log.Printf("pt: (%v, %v)", x, y)
		img.Rect = img.Rect.Add(image.Pt(x, y))
		//log.Printf("NEW img.Rect: %v", img.Rect)
	}

	updateStatus(r, "Upload successful")

	// Display form
	getHandler(w, r)
}

func updateStatus(r *http.Request, status string, values ...any) {
	if values != nil {
		status = fmt.Sprintf(status, values)
	}
	sessionManager.Put(r.Context(), "response", status)
}

func main() {
	// Initialize session
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	// Handle POST and GET requests
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postHandler(w, r)
			return
		}
		getHandler(w, r)
	})

	// Include static content
	mux.Handle("/state/", http.StripPrefix("/state/", http.FileServer(http.Dir("state"))))

	// Add middleware
	muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	// Start server
	fmt.Printf("listening on http://localhost%v\n", PORT)
	if err := http.ListenAndServe(PORT, muxWithSessionMiddleware); err != nil {
		log.Printf("error listening: %v", err)
	}
}
