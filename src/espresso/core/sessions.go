/**
* This function handles all the user session stuff
* so developers can track their progress.
*/

package core

import (
    "net/http"
    "strings"
    "io"
    "io/ioutil"
    "log"
    "os"
    "time"
    "crypto/md5"
    "encoding/hex"
)

type SESSOBJ struct {
	Storepath	string
	inited		bool
}

func InitSessions(sessionStore string) *SESSOBJ {
	SOBJ := &SESSOBJ{}
	SOBJ.Storepath = sessionStore
	SOBJ.inited = true
	return SOBJ
}

// This starts the session, returning the unique session ID back to Espresso
func (D SESSOBJ) StartSession(w http.ResponseWriter, r *http.Request) string {

	// Grab the session ID from the cookie
	session_cookie, err_sess := r.Cookie("gote_session");
	needNewSession := false
	session_name := ""

	// If the user doesn't have a gote_session, set this to error out later and create a new session.
	if err_sess != nil {
		needNewSession = true	
		session_name = "//"
	} else {
		session_name = session_cookie.Value;
	}

	if(strings.Contains(session_name, "/")) {
		needNewSession = true
		session_name = "invalid_session"
	}
	
	// Attempt to read the session
	_, err := ioutil.ReadFile(D.Storepath + "/" + session_name);
	if err != nil {
		log.Println("Error loading usersession. Must create new session. Session: " + session_name)
		needNewSession = true	
	}

	// We need a new session.
	if needNewSession {
		// TODO: Check if the session expired
		time_seed := time.Now().String()
		h := md5.New()
		io.WriteString(h, time_seed)
		io.WriteString(h, r.RemoteAddr)
		b := h.Sum(nil)
		session_id := hex.EncodeToString(b)
		log.Println("Creating new session with Session ID of " + session_id);
	    	cookie := &http.Cookie{Name:"gote_session", Value:session_id, Expires:time.Now().Add(356*24*time.Hour), HttpOnly:true}
		http.SetCookie(w, cookie)
		_, err := os.Create(D.Storepath + "/" + session_id)
		if err != nil {
			log.Fatal("CRITICAL ERROR: Cannot write to session path!")		
		}
		return string(session_id)
	} else {
		//log.Println("Session found with ID: " + string(session_name) + " and data: " + string(cookie_data))
		return string(session_name)
	}
}

func (D SESSOBJ) RetrieveSession(currentSession string) string {
	SessionFile := D.Storepath + "/" + currentSession
	// See if the session still exists on the server
	if _, err := os.Stat(SessionFile); err != nil {
		// TODO: Check if the session expired
		log.Print(SessionFile + " cannot be loaded! Session must not exist!")
	} else {
		// Load the session
		file_data, err := ioutil.ReadFile(SessionFile)
		if err != nil {
			log.Print(SessionFile + "cannot be read! ")
			return ""
		}
		return string(file_data)
	}

	return ""
}

func (D SESSOBJ) SetSession(currentSession string, data string) {
	SessionFile := D.Storepath + "/" + currentSession
	if _, err := os.Stat(SessionFile); err != nil {
		// TODO: Check if the session expired
		log.Print("SetSession: " + SessionFile + " cannot be loaded! Session must not exist!")
	} else {
		fdir, err := os.Create(SessionFile)
		if(err != nil) {
			log.Print("Failed to open session!")
		} else {
			_, err := fdir.WriteString(data)
			if err != nil {
				log.Print("Failed to write to session! Session: " + SessionFile + " Contents: " + data)
				log.Fatal(err)
			}
		}
	}
}

