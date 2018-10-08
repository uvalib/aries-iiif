package main

import (
	"fmt"
	"net/http"
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

/* Handles a request for information about a single PID */
func iiifHandlePid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	logger.Printf("%s %s", r.Method, r.RequestURI)

	pid := params.ByName("pid")

	var derivativeFile, serviceUrl string

	if strings.HasPrefix(pid, "shanti-image") {
		// pid format: mandala-image{,-dev}-id
		pcs := strings.Split(pid, "-")
		pidId := pcs[len(pcs)-1]

		derivativeFile = config.mandalaDirPrefix.value + buildSubPath(pidId) + "/" + pid + ".jp2"
		serviceUrl = config.mandalaUrlTemplate.value
	} else {
		// pid format: {uva-lib,tsm}:id
		pcs := strings.Split(pid, ":")
		pidType := strings.Join(pcs[:len(pcs)-1], ":")
		pidId := pcs[len(pcs)-1]

		derivativeFile = config.iiifDirPrefix.value + "/" + pidType + buildSubPath(pidId) + "/" + pidId + ".jp2"
		serviceUrl = config.iiifUrlTemplate.value
	}

	serviceUrl = strings.Replace(serviceUrl, "{PID}", pid, 1)

	json := fmt.Sprintf(`{ "identifier": [ "%s" ], "service_url": [ { "url": "%s", "protocol": "iiif" } ], "derivative_file": [ "%s" ] }`, pid, serviceUrl, derivativeFile)

	fmt.Fprintf(w, json)
}
