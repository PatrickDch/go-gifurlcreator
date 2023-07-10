package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var APIKey = ""

type Responseer struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Url    string `json:"url"`
}

type Poster struct {
	RegID string `json:"regid"`
}

type Target struct {
	Results []struct {
		Media_formats struct {
			Tinygif struct {
				Url string `json:"url"`
			} `json:"tinygif"`
		} `json:"media_formats"`
	} `json:"results"`
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func main() {

	secrouter := httprouter.New()
	secrouter.POST("/bs", BS)
	fmt.Println("Server started on :8099 port")
	log.Fatal(http.ListenAndServe(":8099", secrouter))

}

func BS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var po Poster
	err := json.NewDecoder(r.Body).Decode(&po)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	urlStr := string(po.RegID)
	fmt.Println(urlStr)
	/* search for excited top 8 GIFs */

	urlflag := "https://tenor.googleapis.com/v2/search?q=fixit&random=true&media_filter=tinygif&limit=1&key=" + APIKey
	fmt.Println(urlflag)
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlflag, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64; rv:114.0) Gecko/20100101 Firefox/114.0`)
	respo, err := client.Do(req)

	if respo.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", respo.StatusCode, respo.Status)
		fmt.Println(respo.StatusCode)
	}
	fmt.Println(respo.StatusCode, respo.Status)
	defer respo.Body.Close()
	decoder := json.NewDecoder(respo.Body)
	var t Target
	decoder.Decode(&t)

	gdata := t.Results[0].Media_formats.Tinygif.Url

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(gdata)

	var resp *Responseer

	if gdata != "" {
		resp = &Responseer{Status: http.StatusOK, Msg: "Gif Found", Url: gdata}

	} else {
		resp = &Responseer{Status: http.StatusNotFound, Msg: "No gif found", Url: ""}
	}

	respJson, err := json.Marshal(resp)

	if err != nil {
		fmt.Fprint(w, "Error occurred while creating json response")
		return
	}

	fmt.Fprint(w, string(respJson))
}
