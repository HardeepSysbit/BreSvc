package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	bre "github.com/HardeepSysbit/bre"
	"github.com/gorilla/mux"
)

var (
	apiVer        = "/api/v1/"
	Port          string
	JwtExpiryTime time.Duration = 0
)

// Http Handlers to handle incoming API Calls
func HandleReq() {

	// Instantiate Router
	muxRouter := mux.NewRouter()

	// Set Routes
	muxRouter.HandleFunc(apiVer+"brePkg", setBrePkg).Methods("PUT")
	muxRouter.HandleFunc(apiVer+"brePkg", exeBrePkg).Methods("POST")

	fmt.Println("Waiting on Port :" + Port)

	// err := http.ListenAndServeTLS(":"+Port, "cert/certificate.crt", "cert/privateKey.key", muxRouter)
	err := http.ListenAndServe(":"+Port, muxRouter)
	if err != nil {
		log.Fatal(err)

	}
}

// Intitializes the BRE package with rules and corresponding actions
func setBrePkg(w http.ResponseWriter, r *http.Request) {

	// Read Message Body
	reqBody, err := ioutil.ReadAll(r.Body)

	// Error in Body so return error response
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Please supply BRE package information in JSON format"))
		return
	}

	// Send Body to BRE
	success, err := bre.SetBrePkg(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	if success {
		// Return Respose
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("BRE Package Accepted"))
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("Unable to compile"))
	}
}

// Sends current facts to the rules engines for processing and returns the results
func exeBrePkg(w http.ResponseWriter, r *http.Request) {

	// Read Message Body
	reqBody, err := ioutil.ReadAll(r.Body)

	// Error in Body so return error response
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 -Please supply BRE facts information in JSON format"))
		return
	}

	// Send Body to BRE to Excute
	facts, err := bre.ExeBrePkg(reqBody)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	} else {

		factsJson, err := json.Marshal(facts)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(factsJson))

	}
}
