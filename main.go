package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Plugin struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Type        string `json:"type"`
	Owner       string `json:"owner"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Forks       int    `json:"fork_count"`
	Stars       int    `json:"star_count"`
	Watchers    int    `json:"watch_count"`
	Issues      int    `json:"issues_count"`
}

//Global Variables
var plugins []Plugin     
var w http.ResponseWriter 


func Filter(vs []Plugin, f func(Plugin) bool) []Plugin {
	vsf := make([]Plugin, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, `Snap Plugin API Server:

/plugins
/plugins/collector
/plugins/processor
/plugins/publisher
/plugin/:name`)
}

func ListPlugin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func ListPlugins(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	
	plugin_type := strings.ToLower(ps.ByName("type"))
	if plugin_type != "" {
		for _, value := range(plugins) {
			if(plugin_type == value.Type){
			     output, _ := json.MarshalIndent(value, "", "    ")
		         fmt.Fprint(w, string(output))
			}
		}
	} else {
		output, _ := json.MarshalIndent(plugins, "", "    ")
		fmt.Fprint(w, string(output))
	}
}	


func ListOwners(w http.ResponseWriter, r *http.Request, ps httprouter.Params){

	plugin_type := strings.ToLower(ps.ByName("name"))
		for _, value := range(plugins) {
			if plugin_type != "" && plugin_type == value.Owner {
		         output, _ := json.MarshalIndent(value, "", "    ")
		         fmt.Fprint(w, string(output))
			} else {
				if plugin_type == "" {
		     		 output, _ := json.MarshalIndent(value.Owner, "," , value.Name)
					 fmt.Fprintln(w, string(output))
			}
	  	}
	}	
}	



func main() {
	//Load the File
	file, e := ioutil.ReadFile("./plugins.json")
	if e != nil {
		fmt.Fprintf(w, "File error: %v\n", e)
	}
	
	err := json.Unmarshal(file, &plugins)
	if err != nil {
		fmt.Fprint(w, err)
	} 

    //Load Plugins
	plugins = Filter(plugins, func(v Plugin) bool {
				return strings.Contains(v.Type, "")
			})


	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/plugins", ListPlugins)
	router.GET("/plugins/:type", ListPlugins)
	router.GET("/plugin/:name", ListPlugin)
	router.GET("/authors", ListOwners) 
	router.GET("/author/:name", ListOwners) 
	

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}
