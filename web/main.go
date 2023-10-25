package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/marcelstanley/millionerc20/client"
	dapp_client "github.com/marcelstanley/millionerc20/client"
	dapp_cli_image "github.com/marcelstanley/millionerc20/image"
)

const (
	PORT            = ":3000"
	MAX_UPLOAD_SIZE = 1024 * 1024 // 1MB
)

type GlobalState struct {
	DappImage *image.RGBA
}

var global GlobalState
var sessionManager *scs.SessionManager

func getHandler(w http.ResponseWriter, r *http.Request) {
	// Display submit response, when available
	response := sessionManager.GetString(r.Context(), "response")
	component := page(response)
	component.Render(r.Context(), w)
	//log.Printf("response: %v\n", response)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// Update state
	// TODO Should we limit the form size? I don't think so
	//	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

	// TODO Make sure we repond with an error if the image is too big
	// XXX This does not seemt to be working
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		updateStatus(w, r, "Max upload size of 1MB exceeded")
		return
	}

	// Get file from request
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		updateStatus(w, r, err.Error())
		return
	}
	defer file.Close()

	f, err := fileHeader.Open()
	if err != nil {
		updateStatus(w, r, err.Error())
		return
	}

	img, err := dapp_cli_image.Decode(f)
	if err != nil {
		updateStatus(w, r, err.Error())
		return
	}
	//log.Printf("img.Rect: %v\n", img.Rect)

	x, err := strconv.Atoi(r.FormValue("posX"))
	if err != nil {
		x = 0
	}
	y, err := strconv.Atoi(r.FormValue("posY"))
	if err != nil {
		y = 0
	}

	// updateStatus(w, r, "Sending image with bounds %v", img.Rect.Bounds())
	log.Printf("sending image with bounds %v", img.Rect.Bounds())

	// submit image to dapp and capture result
	_, err = dapp_client.SendImageAndWait(image.Pt(x, y), img)
	if err != nil {
		log.Printf("err: %v\n", err)
		updateStatus(w, r, err.Error())
		return
	}

	getHandler(w, r)
}

func updateDappImage(w http.ResponseWriter) {
	global.DappImage = client.GetDappImage()
	png.Encode(w, global.DappImage)
}

// TODO Understand why a call to updateStatus triggers an additional call to getHandler under the hood
// TODO use goroutine for this?
func updateStatus(w http.ResponseWriter, r *http.Request, status string, values ...any) {
	if values != nil {
		status = fmt.Sprintf(status, values)
	}
	sessionManager.Put(r.Context(), "response", status)
	getHandler(w, r)
}

func main() {
	// Initialize session
	sessionManager = scs.New()
	sessionManager.Lifetime = 24 * time.Hour

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			//TODO refactor postHnadler to return and error, which would updateStatus and let getHandler happen always
			postHandler(w, r)
			return
		}
		getHandler(w, r)
	})

	mux.HandleFunc("/dapp_image", func(w http.ResponseWriter, r *http.Request) {
		updateDappImage(w)
	})

	// Include static content
	//mux.Handle("/state/", http.StripPrefix("/state/", http.FileServer(http.Dir("state"))))

	// Add middleware
	muxWithSessionMiddleware := sessionManager.LoadAndSave(mux)

	// Start server
	log.Printf("listening on http://localhost%v\n", PORT)
	if err := http.ListenAndServe(PORT, muxWithSessionMiddleware); err != nil {
		log.Printf("error listening: %v", err)
	}
}
