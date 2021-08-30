package preconditions

import (
	"database/sql"
	"fmt"
	"reflect"
)

func NonNilDB(db *sql.DB, msg string) *sql.DB {
	if db == nil {
		panic(fmt.Sprintf("NonNilDB: %s", msg))
	}
	return db
}

func NonNilStmt(stmt *sql.Stmt, msg string) *sql.Stmt {
	if stmt == nil {
		panic(fmt.Sprintf("NonNilStmt: %s", msg))
	}
	return stmt
}

func NonEmptyString(val string, msg string) string {
	if val == "" {
		panic(fmt.Sprintf("NonEmptyString: %s", msg))
	}
	return val

}

func AssertNonNil(val interface{}, msg string) {
	if val == nil || reflect.ValueOf(val).IsNil() {
		panic(fmt.Sprintf("AssertNonNil: %s", msg))
	}
}

func AssertNil(val interface{}, msg string) {
	if !(val == nil || reflect.ValueOf(val).IsNil()) {
		panic(fmt.Sprintf("AssertNil: %s", msg))
	}
}

func AssertUniqueStrings(s []string) []string {
	for i := 0; i < len(s); i++ {
		for j := i + 1; j < len(s); j++ {
			if s[i] == s[j] {
				panic(fmt.Errorf("You used duplicated elements in sequence: '%+v'.", s))
			}
		}
	}
	return s
}
