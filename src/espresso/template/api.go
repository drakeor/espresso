package template

import "fmt"
import "espresso/core"
import "io"
import "encoding/hex"
import "strings"
import "regexp"
import lua "github.com/aarzilli/golua/lua"
import "crypto/md5"
import "crypto/sha1"
import "crypto/sha256"
import "crypto/sha512"
import "database/sql"
import "net/http"
import "log"

type StringArray struct {
	Str []string
}

var WebRoot string

var DBRoot *core.DBO

var SessRoot *core.SESSOBJ

// Workaround for sessions.
var CurrentWrite http.ResponseWriter
var CurrentReq http.Request
var CurrentGetURI string
var CurrentPostValues map[string][]string
var CurrentSessionID string

// Workaround for sessions
func SetCurrentClient(w http.ResponseWriter, req http.Request, getURI string) {
	CurrentReq = req
	CurrentWrite = w
	CurrentGetURI = getURI
	CurrentSessionID = ""
	CurrentReq.ParseForm()
	//log.Println("New request from: " + req.RemoteAddr + " Accessing Resource: " + req.RemoteURL);
	CurrentPostValues = CurrentReq.Form
}

func InitAPI(L *lua.State, webroot string, DB *core.DBO, SESS *core.SESSOBJ) {
	// Form Data
	//L.Register("POST", POSTfunction) REQUEST()
	//L.Register("GET", GETfunction) REQUEST()
	L.Register("REQUEST", REQUESTfunction)

	// Sessions
	L.Register("session_start", StartSession)
	L.Register("session_set", SetSession)
	L.Register("session_get", GetSession)

	// Non-Documented functions. These already exist in lua in some way.
	// L.Register("strlen", StrLen) string.length()
	// L.Register("strrev", StrRev) string.reverse()
	// L.Register("substr", SubStr) string.gsub()
	// L.Register("strpos", StrPos) 
	
	// String
	L.Register("explode", Explode)
	L.Register("implode", Implode)
	L.Register("trim", Trim)

	// Hashing
	L.Register("md5", Md5)
	L.Register("sha1", Sha1)
	L.Register("sha256", Sha256)
	L.Register("sha512", Sha512)

	// Register Database Calls.
	L.Register("mysql_addrow", AddRow)
	L.Register("mysql_queryid", QueryId)
	L.Register("mysql_queryfind", QueryFind)
	L.Register("mysql_querycustom", QueryCustom)
	L.Register("mysql_queryall", QueryAll)

	DBRoot = DB
	SessRoot = SESS
	WebRoot = webroot
}

// TODO: We eventually want to be able to pass GET("name") and have it return the name
func GETfunction(L *lua.State) int {
	str := L.ToString(1)	
	L.PushString(CurrentGetURI)
	L.PushString(CurrentReq.FormValue(str))
	return 1
}

func POSTfunction(L *lua.State) int {
	str := L.ToString(1)
	L.PushString(CurrentReq.FormValue(str))
	//log.Fatal(CurrentReq.FormValue("username"))
	return 1
}

func REQUESTfunction(L *lua.State) int {
	str := L.ToString(1)
	L.PushString(CurrentReq.FormValue(str))
	//log.Fatal(CurrentReq.FormValue("username"))
	return 1
}

func StartSession(L *lua.State) int {
	CurrentSessionID = SessRoot.StartSession(CurrentWrite, &CurrentReq)
	return 0
}

func GetSession(L *lua.State) int {
	if CurrentSessionID == "" {
		log.Print("StartSession needs to be called before GetSession!")
	}
	CurrentSession := SessRoot.RetrieveSession(CurrentSessionID)
	L.PushString(CurrentSession)
	return 1
}

// Sets the session
func SetSession(L *lua.State) int {
	if CurrentSessionID == "" {
		log.Print("StartSession needs to be called before SetSession!")
	}
	str := L.ToString(1)
	SessRoot.SetSession(CurrentSessionID, str)
	return 0
}


// Gets the string's length
// Returns the length of the given string.
func StrLen(L *lua.State) int {
	str := L.ToString(1)
	L.PushInteger(int64(len(str)))
	return 1
}

// Reverses a string
// Returns string, reversed
func StrRev(L *lua.State) int {
	str := L.ToString(1)
	strRev := ""
	for _, c := range str {
		strRev = string(c) + strRev
	}
	L.PushString(strRev)
	return 1
}

// Returns a section of the string.
// Takes 2 or 3 arguments.
// String, From, Length
func SubStr(L *lua.State) int {
	str := L.ToString(1)
	from := L.ToInteger(2)
	to := L.ToInteger(3)
	if to == 0 {
		// To not declared. Make it the length of the string - from.
		to = len(str) - from
	}

	L.PushString(str[from : from+to])
	return 1
}

// Takes a lua string and breaks it into a lua table.
func Explode(L *lua.State) int {
	haystack := L.ToString(1)
	needle := L.ToString(2)
	arr := strings.Split(haystack, needle)
	L.CreateTable(len(arr), 0)
	i := 1
	for _, str := range arr {
		L.PushString(str)
		L.RawSeti(-2, i)
		i++
	}

	return 1
}

// Inputs a lua table of strings and outputs a string with the dividers
// inside it.
func Implode(L *lua.State) int {
	div := L.ToString(1)
	L.PushNil()
	str := ""
	for L.Next(2) != 0 {
		str = str + L.ToString(-1) + div
		L.Pop(1)
	}
	str = str[0 : len(str)-len(div)]
	L.PushString(str)
	return 1
}

// Trims a string.
func Trim(L *lua.State) int {
	str := L.ToString(1)
	L.PushString(strings.Trim(str, " "))
	return 1
}

// First position in a string
func StrPos(L *lua.State) int {
	str := L.ToString(1)
	needle := L.ToString(2)

	reg := regexp.MustCompile(needle)
	d := reg.FindAllStringIndex(str, -1)

	if len(d) > 0 {
		L.PushInteger(int64(d[0][0]))
	} else {
		L.PushBoolean(false)
	}
	return 1
}

// Hashing

func Md5(L *lua.State) int {
	str := L.ToString(1)
	h := md5.New()
	io.WriteString(h, str)
	b := h.Sum(nil)
	L.PushString(hex.EncodeToString(b))
	return 1
}

func Sha1(L *lua.State) int {
	str := L.ToString(1)
	h := sha1.New()
	io.WriteString(h, str)
	b := h.Sum(nil)
	L.PushString(hex.EncodeToString(b))
	return 1
}

func Sha256(L *lua.State) int {
	str := L.ToString(1)
	h := sha256.New()
	io.WriteString(h, str)
	b := h.Sum(nil)
	L.PushString(hex.EncodeToString(b))
	return 1
}

func Sha512(L *lua.State) int {
	str := L.ToString(1)
	h := sha512.New()
	io.WriteString(h, str)
	b := h.Sum(nil)
	L.PushString(hex.EncodeToString(b))
	return 1
}

// Queries all rows in a table.
func QueryAll(L *lua.State) int {
	table := L.ToString(1)

	rows := DBRoot.QueryAll(table)

	ExportRows(L, rows)

	return 1
}

// Takes a variadic amount of values for the prepared statement. Custom query.
func QueryCustom(L *lua.State) int {
	query := L.ToString(1)
	if L.GetTop()-1 != 0 {
		values := make([]interface{}, L.GetTop()-1, L.GetTop()-1)

		for i, _ := range values {
			values[i] = L.ToString(i + 3)
		}

		rows := DBRoot.QueryCustom(query, values...)

		ExportRows(L, rows)
	} else {
		rows := DBRoot.QueryCustom(query)

		ExportRows(L, rows)
	}

	return 1
}

// Finds a specified value in a table that matches a given field.
func QueryFind(L *lua.State) int {
	table := L.ToString(1)
	field := L.ToString(2)
	value := L.ToString(3)

	rows := DBRoot.QueryFind(table, field, value)

	ExportRows(L, rows)

	return 1
}

// Queries a specific id from a given table.
func QueryId(L *lua.State) int {
	table := L.ToString(1)
	value := L.ToInteger(2)

	rows := DBRoot.QueryId(table, value)

	ExportRows(L, rows)

	return 1
}

// Queries a specific id from a given table.
func AddRow(L *lua.State) int {
	table := L.ToString(1)
	data := make(map[string]string, L.ObjLen(2))
	L.PushNil()
	for L.Next(2) != 0 {
		data[L.ToString(-2)] = L.ToString(-1)
		L.Pop(1)
	}

	DBRoot.AddRow(table, data)

	return 0
}

// Exports Sql Rows as a Lua Table
func ExportRows(L *lua.State, rows *sql.Rows) {
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	z := 1
	L.NewTable() // 1

	// Fetch rows
	for rows.Next() {
		L.NewTable()

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// Print data
		for i, value := range values {
			switch value.(type) {
			case nil:
				L.PushString("NULL")
			case []byte:
				L.PushString(string(value.([]byte)))
			case int:
				L.PushString(fmt.Sprintf("%d", value))
			case int16:
				L.PushString(fmt.Sprintf("%d", value))
			case int32:
				L.PushString(fmt.Sprintf("%d", value))
			case int64:
				L.PushString(fmt.Sprintf("%d", value))
			case int8:
				L.PushString(fmt.Sprintf("%d", value))
			case bool:
				if value == true {
					L.PushString("1")
				} else {
					L.PushString("0")
				}
			default:
				L.PushString("NULL")
			}
			L.SetField(-2, string(columns[i]))
		}

		L.RawSeti(-2, z)
		z++
	}
}
