/**
* The file handles all the database operations internally
*/


package core

import "fmt"
import "log"
import "database/sql"
import _ "github.com/go-sql-driver/mysql"

type DBO struct {
	con    *sql.DB 	// SQL Connection Pool
	err    error	// Error holder
	inited bool   	// Whether the DBO has been created.
	Prefix string 	// Table Prefix
}

// Creates the DBO and connects to the Database.
func InitDBO(host string, dbUser string, dbPassword string, dbName string) *DBO {
	return InitDBOWithPort(host, dbUser, dbPassword, dbName, 3306)
}

// Creates the DBO with a given port.
func InitDBOWithPort(host string, dbUser string, dbPassword string, dbName string, port int) *DBO {

	// Connect to the MYSQL Database here
	D := &DBO{}
	q := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", dbUser, dbPassword, host, port, dbName)
	D.inited = true
	D.con, D.err = sql.Open("mysql", q)
	if D.err != nil {
		log.Println("SQL database connection failure: ", D.err)
	}

	return D
}

// Queries all rows in a table.
func (D DBO) QueryAll(table string) *sql.Rows {
	if D.inited {
		q := fmt.Sprintf("SELECT * FROM %s", table)
		rows, err := D.con.Query(q)
		if err != nil {
			log.Println("SQL Error: ", err)
			return &sql.Rows{}
		} else {
			return rows
		}
	} else {
		return &sql.Rows{}
	}
}

// Takes a variadic amount of values for the prepared statement. Custom query.
func (D DBO) QueryCustom(query string, values ...interface{}) *sql.Rows {
	if D.inited {
		stmt, err := D.con.Prepare(query)
		if err != nil {
			log.Println("SQL Error: ", err)
			return &sql.Rows{}
		} else {
			rows, err2 := stmt.Query(values...)
			if err2 != nil {
				log.Println("SQL Error: ", err2)
				stmt.Close()
				return &sql.Rows{}
			} else {
				stmt.Close()
				return rows
			}
		}
	} else {
		return &sql.Rows{}
	}
}

// Finds a specified value in a table that matches a given field.
func (D DBO) QueryFind(table string, field string, value string) *sql.Rows {
	if D.inited {
		q := fmt.Sprintf("SELECT * FROM %s WHERE `%s` = ?", table, field)
		return D.QueryCustom(q, value)
	} else {
		return &sql.Rows{}
	}
}

// Queries a specific id from a given table.
func (D DBO) QueryId(table string, id int) *sql.Rows {
	if D.inited {
		q := fmt.Sprintf("SELECT * FROM %s WHERE `uid` = ?", table)
		return D.QueryCustom(q, id)
	} else {
		return &sql.Rows{}
	}
}

// Queries a specific id from a given table.
func (D DBO) AddRow(table string, data map[string]string) int64 {
	if D.inited {
		pre := fmt.Sprintf("INSERT INTO %s (", table)
		post := ") VALUES ("

		values := make([]interface{}, len(data), len(data))
		i := 0
		for key, val := range data {
			pre = pre + fmt.Sprintf("`%s`, ", key)
			post = post + "?, "
			values[i] = val
			i++
		}

		pre = pre[0 : len(pre)-2]
		post = post[0:len(post)-2] + ")"
		q := pre + post

		stmt, err := D.con.Prepare(q)

		if err != nil {
			log.Println("SQL Error: ", err)
			return -1
		} else {
			res, err2 := stmt.Exec(values...)

			if err2 != nil {
				log.Println("SQL Error: ", err2)
				stmt.Close()
				return -1
			} else {
				stmt.Close()
				last, _ := res.LastInsertId()
				return last
			}
		}

	} else {
		return -1
	}
}

func (D DBO) Close() {
	D.con.Close()
}
