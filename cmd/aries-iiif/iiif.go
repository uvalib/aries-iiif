package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
)

/* Builds the subpath to the source image for a given ID, where subpath
 * is the ID split into subdirectories in groups of two.
 * e.g. ID "123456789" => subpath "/12/34/56/78/9"
 */
func buildSubPath(id string) (subPath string) {

	chars := strings.Split(id, "")

	subPath = ""

	i := 0

	for _, c := range chars {
		if i%2 == 0 {
			subPath = subPath + "/"
		}
		i++

		subPath = subPath + c
	}

	logger.Printf("[%s] => [%s]", id, subPath)

	return
}

// if this is an IIIF PID, returns the service URL and derivative file
func processIiifPid(pid string) (string, string, error) {
	// valid pid forms:
	// tsm:1234567
	// uva-lib:1234567
	re := regexp.MustCompile(`^(tsm|uva-lib):\d+$`)

	if !re.MatchString(pid) {
		logger.Printf("%s is NOT IIIF", pid)
		return "", "", errors.New("PID is not IIIF")
	}

	pos := strings.LastIndex(pid, ":")
	pidType := pid[:pos]
	pidId := pid[pos+1:]

	derivativeFile := config.iiifDirPrefix.value + "/" + pidType + buildSubPath(pidId) + "/" + pidId + ".jp2"
	serviceUrl := config.iiifUrlTemplate.value

	return derivativeFile, serviceUrl, nil
}

// if this is a Mandala PID, returns the service URL and derivative file
func processMandalaPid(pid string) (string, string, error) {
	// valid forms:
	// shanti-image-1234567
	// shanti-image-dev-1234567
	re := regexp.MustCompile(`^shanti-image-(|dev-)\d+$`)

	if !re.MatchString(pid) {
		logger.Printf("%s is NOT Mandala", pid)
		return "", "", errors.New("PID is not Mandala")
	}

	pos := strings.LastIndex(pid, "-")
	pidId := pid[pos+1:]

	derivativeFile := config.mandalaDirPrefix.value + buildSubPath(pidId) + "/" + pid + ".jp2"
	serviceUrl := config.mandalaUrlTemplate.value

	return derivativeFile, serviceUrl, nil
}

// if this is a supported PID, returns the service URL and derivative file
func processPid(pid string) (string, string, error) {
	if x, y, err := processIiifPid(pid); err == nil {
		return x, y, err
	}

	if x, y, err := processMandalaPid(pid); err == nil {
		return x, y, err
	}

	return "", "", errors.New("PID is neither IIIF nor Mandala")
}

/* Handles a request for information about a single PID */
func iiifHandlePid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)

	pid := params.ByName("pid")

	var derivativeFile, serviceUrl string
	var err error

	// get file and url info, if this is a known PID type
	if derivativeFile, serviceUrl, err = processPid(pid); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Printf("processing failed for PID: [%s]", pid)
		fmt.Fprintf(w, "Invalid PID: %s", pid)
		return
	}

	// ensure the derivative file exists?
	if config.ensureExists.value {
		if _, err = os.Stat(derivativeFile); os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			logger.Printf("derivative file does not exist: %s", derivativeFile)
			fmt.Fprintf(w, "Derivative file not found: %s", derivativeFile)
			return
		}
	}

	// reformat url with actual PID
	serviceUrl = strings.Replace(serviceUrl, "{PID}", pid, 1)

	// build Aries API response object
	var iiifResponse AriesAPI

	iiifResponse.addIdentifier(pid)
	iiifResponse.addDerivativeFile(derivativeFile)
	iiifResponse.addServiceUrl(ServiceUrl{Url: serviceUrl, Protocol: "iiif"})

	w.Header().Set("Content-Type", "application/json")

	j, jerr := json.Marshal(iiifResponse)
	if jerr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Printf("JSON marshal failed: [%s]", jerr.Error())
		fmt.Fprintf(w, "JSON marshal failed")
		return
	}

	fmt.Fprintf(w, string(j))
}
