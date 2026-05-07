// package pg mirrors the generated-app postgres helper surface needed by copied source files.
package pg

import (
	"database/sql"
)

// pg is assigned by generated applications before helper functions issue queries.
var Pg *sql.DB
