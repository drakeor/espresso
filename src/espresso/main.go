package main

import "espresso/core"
import "espresso/template"
import "espresso/utils/rand"
import "net/http"
import "io"
import "os"
import "log"
import "sync"
import "time"
import "fmt"
import "strings"
import lua "github.com/aarzilli/golua/lua"

type Espresso struct {
	Name   string
	Core   *core.Core
	Config *core.Config
}

type Handler struct {
	espresso *Espresso
	lua  *lua.State
	C    *core.Core
	Conf *core.Config
	mu   sync.Mutex
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	
	// Setup the GoTE object for this session.
	espresso := &Espresso{}
	espresso.Name = h.Conf.Params.Name
	espresso.Core = h.C
	espresso.Config = h.Conf
	

	// Hack off the GET request at the end of the URL
	fullURI := strings.Split(req.RequestURI, "?")
	getURI := ""
	if len(fullURI) > 1 {
		getURI = fullURI[1]
	}
	fileURI := fullURI[0]
	requestedFile := espresso.Config.Params.Webroot + fileURI

	// If they used / instead of specifying an index, search for an index file
	if fileURI[len(fileURI)-1:len(fileURI)] == "/" {
		if _, err := os.Stat(requestedFile + "index.lua"); err == nil {
			requestedFile = requestedFile + "index.lua"
		} else if _, err := os.Stat(requestedFile + "index.htm"); err == nil {
			requestedFile = requestedFile + "index.htm"
		} else if _, err := os.Stat(requestedFile + "index.html"); err == nil {
			requestedFile = requestedFile + "index.html"
		} else {
			requestedFile = requestedFile + "none.404"
		}
	}

	// If the requested file is a lua file, we want to grab it.
	if requestedFile[len(requestedFile)-3:len(requestedFile)] == "lua" {

		// We want to time how long it takes to generate the page
		var startTime = time.Now()

		// Parse the template for Lua Syntax		
		template.SetCurrentClient(w, *req, getURI)		
		tmpl := template.ParseTemplate(requestedFile)

		// Set up lua session and register initial functions
		L := lua.NewState()
		L.OpenLibs()
		defer L.Close()
		template.InitAPI(L, h.Conf.Params.Webroot, h.C.DB, h.C.SESS)
		L.Register("print", func(L *lua.State) int {
			txt := L.ToString(1)
			io.WriteString(w, txt)
			return 0
		})

		// Now parse the requested file
		err := L.DoString(tmpl)
		if err != nil {
			log.Print("ERROR: Caught Lua Error: ", err)
		}
		
		// Grab the page render time and display it in the server console
		var endTime = time.Now()
		var timeDif = (endTime.Nanosecond() - startTime.Nanosecond()) / 1000000
		log.Printf("Time to render page (ms): %d", timeDif)

	} else {

		// Send the file contents over
		_, err2 := os.Stat(requestedFile)
		// TODO: Change this to ServeFile
		if err2 == nil {
			http.ServeFile(w, req, requestedFile)
		} else {
			io.WriteString(w, "<html><head></head><body><strong>404: "+requestedFile+"</strong></body></html>")
		}
	}
	
}

func main() {
	
	// Initialize Lua
	L := lua.NewState()
	L.OpenLibs()
	defer L.Close()

	// Initialize core components
	Config := core.InitConfig()
	Core := core.InitCore(Config)
	MTTEST := rand.CreateRandom()
	MTTEST.Seed(82812722)
	defer Core.Close()

	// Register our custom lua functions and redirect "print" to the console
	template.InitAPI(L, Config.Params.Webroot, Core.DB, Core.SESS)
	L.Register("print", func(L *lua.State) int {
		txt := L.ToString(1)
		log.Print(txt)
		return 0
	})

	// Our handling function that parses the requests
	var h Handler
	h.lua = L
	h.C = Core
	h.Conf = Config

	// DB Test
	// print(Core.DB.QueryCustom("SELECT * FROM sandbox_accounts.account"))

	// Register handler
	http.Handle("/", h)
	listenUrl := fmt.Sprintf("%s:%d", Config.Params.IP, Config.Params.Port)
	log.Print("Initialized Espresso. Attempting to listen on " + listenUrl)
	err := http.ListenAndServe(listenUrl, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
