// Author: George Fang
// Note that this project is not enorsed or affliated with MARTA in any way.

package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	//"github.com/joho/godotenv" // for local dev use
	//"github.com/rs/cors" // for local dev use
	"slices"
)

type trainArrival []struct {
	Destination    string `json:"DESTINATION"`
	Direction      string `json:"DIRECTION"`
	EventTime      string `json:"EVENT_TIME"`
	IsRealTime     string `json:"IS_REALTIME"`
	Line           string `json:"LINE"`
	NextArrival    string `json:"NEXT_ARR"`
	Station        string `json:"STATION"`
	TrainID        string `json:"TRAIN_ID"`
	WaitingSeconds string `json:"WAITING_SECONDS"`
	WaitingTime    string `json:"WAITING_TIME"`
	Delay          string `json:"DELAY"`
	Latitude       string `json:"LATITUDE"`
	Longitude      string `json:"LONGITUDE"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Train Arrival Service API Started!")
	// loadEnvError := godotenv.Load(".env.local")

	// if loadEnvError != nil {
	// 	log.Println("Cannot load environment variables from .env.local file")
	// }

	allowedOrgins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	origin := r.Header.Get("Origin")
	log.Println("A request has been received. The origin of the request is:", origin)

	// func Contains(s, v) reports whether v is present in s
	allowAccess := slices.Contains(allowedOrgins, origin)

	// preflight checks
	if r.Method == "OPTIONS" {
        if allowAccess {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            return
        } else {
            http.Error(w, "Access denied. Origin not accepted.", http.StatusForbidden)
            return
        }
    }

	if !allowAccess {
		http.Error(w, "Access denied. Origin not accepted.", http.StatusForbidden)
		return
	}

	apiKey := os.Getenv("API_KEY")

	// add apiKey to URL
	martaURL := fmt.Sprintf("https://developerservices.itsmarta.com:18096/itsmarta/railrealtimearrivals/developerservices/traindata?apiKey=%s", apiKey)

	// Get train data from MARTA
	response, httpGetErr := http.Get(martaURL)
	if httpGetErr != nil {
		log.Println("Something went wrong during GET.", httpGetErr)
		return
	}

	// Read response
	body, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		log.Println("Something went wrong reading reqest", readErr)
		return
	}

	var trainArrivalData trainArrival
	unmarshalErr := json.Unmarshal(body, &trainArrivalData)
	if unmarshalErr != nil {
		log.Println("Error unmarshaling train arrival data", unmarshalErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	encodeErr := json.NewEncoder(w).Encode(trainArrivalData)
	if encodeErr != nil {
		log.Println("Error encoding train arrival data", encodeErr)
		return
	}
}

// func main() {
// 	log.Println("Train Arrival Service API Started!")
// 	loadEnvError := godotenv.Load(".env.local")

// 	if loadEnvError != nil {
// 		log.Println("Cannot load environment variables from .env.local file")
// 	}

// 	mux := http.NewServeMux()
// 	http.HandleFunc("/", Handler)

// 	// cors.Default() setup the middleware with default options being
// 	// all origins accepted with simple methods (GET, POST). See
// 	// documentation below for more options.
// 	handler := cors.Default().Handler(mux)
// 	http.ListenAndServe(":8080", handler)
// }
