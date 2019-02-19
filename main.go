package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Version of the service
const version = "1.0.0"

// Config info; APTrust host and API key
var solrURL string
var solrCore string

// favHandler is a dummy handler to silence browser API requests that look for /favicon.ico
func favHandler(c *gin.Context) {
}

// versionHandler reports the version of the serivce
func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Aries Virgo version %s", version)
}

// healthCheckHandler reports the health of the serivce
func healthCheckHandler(c *gin.Context) {
	hcMap := make(map[string]string)
	hcMap["AriesVirgo"] = "true"
	// ping the api with a minimal request to see if it is alive
	url := fmt.Sprintf("%s/%s/select?q=*:*&wt=json&rows=0", solrURL, solrCore)
	_, err := getAPIResponse(url)
	if err != nil {
		hcMap["Virgo"] = "false"
	} else {
		hcMap["Virgo"] = "true"
	}
	c.JSON(http.StatusOK, hcMap)
}

/// ariesPing handles requests to the aries endpoint with no params.
// Just returns and alive message
func ariesPing(c *gin.Context) {
	c.String(http.StatusOK, "APTrust Virgo API")
}

// ariesLookup will query APTrust for information on the supplied identifer
func ariesLookup(c *gin.Context) {
	ID := c.Param("id")
	var qps []string
	qps = append(qps, url.QueryEscape(fmt.Sprintf("id:\"%s\"", ID)))
	qps = append(qps, url.QueryEscape(fmt.Sprintf("alternate_id_facet:\"%s\"", ID)))
	qps = append(qps, url.QueryEscape(fmt.Sprintf("barcode_facet:\"%s\"", ID)))
	fl := "&fl=id,shadowed_location_facet,marc_display,alternate_id_facet,barcode_facet,title_display"
	url := fmt.Sprintf("%s/%s/select?q=%s&wt=json&indent=true%s", solrURL, solrCore, strings.Join(qps, "+"), fl)
	respStr, err := getAPIResponse(url)
	if err != nil {
		log.Printf("Query for %s FAILED: %s", ID, err.Error())
		c.String(http.StatusNotFound, err.Error())
		return
	}
	c.String(http.StatusOK, respStr)
}

// getAPIResponse is a helper used to call a JSON endpoint and return the resoponse as a string
func getAPIResponse(url string) (string, error) {
	log.Printf("Get resonse for: %s", url)
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Unable to GET %s: %s", url, err.Error())
		return "", err
	}

	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	respString := string(bodyBytes)
	if resp.StatusCode != 200 {
		return "", errors.New(respString)
	}
	return respString, nil
}

/**
 * MAIN
 */
func main() {
	log.Printf("===> Aries Virgo service staring up <===")

	// Get config params
	log.Printf("Read configuration...")
	var port int
	flag.IntVar(&port, "port", 8080, "Aries Virgo port (default 8080)")
	flag.StringVar(&solrURL, "solrurl", "http://solr.lib.virginia.edu:8082/solr", "Solr base URL")
	flag.StringVar(&solrCore, "solrcore", "core", "Solr core")
	flag.Parse()

	log.Printf("Setup routes...")
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.GET("/favicon.ico", favHandler)
	router.GET("/version", versionHandler)
	router.GET("/healthcheck", healthCheckHandler)
	api := router.Group("/api")
	{
		api.GET("/aries", ariesPing)
		api.GET("/aries/:id", ariesLookup)
	}

	portStr := fmt.Sprintf(":%d", port)
	log.Printf("Start Aries Virgo v%s on port %s", version, portStr)
	log.Fatal(router.Run(portStr))
}
