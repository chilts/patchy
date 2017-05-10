// Copyright 2017 Andrew Chilton
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package patchy

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var patchKey = "patch"

type Options struct {
	Dir           string
	PropertyTable string
}

type patch struct {
	Forward string
	Reverse string
}

func oneRowOneColBool(rows *sql.Rows) (bool, error) {
	var val bool
	var rowCount int

	for rows.Next() {
		rowCount++

		err := rows.Scan(&val)
		if err != nil {
			return val, err
		}
	}

	return val, rows.Err()
}

func insertPatchLevel(db *sql.DB, tableName string) error {
	sqlIns := `INSERT INTO ` + tableName + `(key, value) VALUES($1, 1)`
	_, err := db.Exec(sqlIns, patchKey)
	return err
}

func createPropertyTable(db *sql.DB, tableName string) error {
	sqlCreateTable := `
		CREATE TABLE ` + tableName + ` (
			key   TEXT PRIMARY KEY,
			value TEXT
		);
    `

	_, err := db.Exec(sqlCreateTable)
	if err != nil {
		return err
	}

	err = insertPatchLevel(db, tableName)
	return err
}

func getCurrentLevel(db *sql.DB, tableName string) (int, error) {
	sqlSel := "SELECT value FROM " + tableName + " WHERE key = $1"
	fmt.Printf("sql=%s\n", sqlSel)
	fmt.Printf("key=%s\n", patchKey)

	var level int
	row := db.QueryRow(sqlSel, patchKey)
	if err := row.Scan(&level); err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("*** NO ROWS\n")

			err := insertPatchLevel(db, tableName)
			return 0, err
		}
		return 0, err
	}

	return level, nil
}

func Patch(db *sql.DB, level int, opts *Options) (int, error) {
	if opts == nil {
		opts = &Options{}
	}

	// set some option defaults
	if opts.Dir == "" {
		opts.Dir = "."
	}
	if opts.PropertyTable == "" {
		opts.PropertyTable = "property"
	}
	fmt.Printf("Options=%#v\n", opts)

	// read the dir list
	files, err := ioutil.ReadDir(opts.Dir)
	if err != nil {
		return 0, err
	}

	// check the files are what we expect
	patchSet := make(map[int]*patch)
	for _, file := range files {
		fmt.Printf("file=%s\n", file.Name())

		patchDirection := ""
		if strings.HasSuffix(file.Name(), "-forward.sql") {
			fmt.Printf("- forward\n")
			patchDirection = "forward"
		} else if strings.HasSuffix(file.Name(), "-reverse.sql") {
			fmt.Printf("- reverse\n")
			patchDirection = "reverse"
		} else {
			fmt.Printf("- unknown\n")
			continue
		}

		// what is the patch level of this filename
		patchLevel := file.Name()
		patchLevel = strings.TrimSuffix(patchLevel, "-"+patchDirection+".sql")

		// check this is a number
		n, err := strconv.Atoi(patchLevel)
		if err != nil {
			return 0, err
		}

		// read this file in
		if _, ok := patchSet[n]; ok {
			fmt.Printf("patch already exists, so adding the other direction (hopefully)\n")
		} else {
			fmt.Printf("no such patch in patchset yet - adding\n")
			patchSet[n] = &patch{}
		}

		filename := path.Join(opts.Dir, file.Name())
		sql, err := ioutil.ReadFile(filename)
		if err != nil {
			return 0, err
		}

		p := patchSet[n]
		if patchDirection == "forward" {
			p.Forward = string(sql)
		}
		if patchDirection == "reverse" {
			p.Reverse = string(sql)
		}
		fmt.Printf("p=%#v\n", p)
	}

	fmt.Printf("PatchSet=%v\n", patchSet)

	// ToDo: check that we have all patches (both Forward and Reverse) up to the level required

	// firstly, figure out if the property table exists
	sqlPropertyTableExists := `
        SELECT EXISTS (
            SELECT
                *
            FROM
                information_schema.tables
            WHERE
                table_schema = 'public'
            AND
                table_name = $1
        );
    `

	// check to see if the property table exists
	var propertyTableExists bool
	row := db.QueryRow(sqlPropertyTableExists, opts.PropertyTable)
	if err := row.Scan(&propertyTableExists); err != nil {
		return 0, err
	}
	if propertyTableExists == false {
		fmt.Printf("Creating property table\n")
		err := createPropertyTable(db, opts.PropertyTable)
		if err != nil {
			return 0, err
		}
	}

	// current patch level
	currentLevel, err := getCurrentLevel(db, opts.PropertyTable)
	if err != nil {
		return 0, err
	}
	fmt.Printf("-> Current Level = %d\n", currentLevel)

	// figure out which direction we're actually going to go in
	direction := ""
	step := 0
	if currentLevel < level {
		direction = "forward"
		step = 1
	}
	if currentLevel > level {
		direction = "reverse"
		step = -1
	}
	if currentLevel == level {
		fmt.Printf("Nothing to do, currently at the same level %d\n", level)
		return level, nil
	}
	fmt.Printf("-> Direction = %s\n", direction)

	for num := currentLevel; num != level; num += step {
		fmt.Printf("- doing from %d to %d\n", num, num+step)

		// loop through all of the patches we know about
		fmt.Printf("-> BEGIN ...\n")
		tx, err := db.Begin()
		if err != nil {
			return 0, err
		}
		fmt.Printf("-> BEGIN done\n")

		// update with this patch
		sql := ""
		if direction == "forward" {
			sql = patchSet[num+step].Forward
		}
		if direction == "reverse" {
			sql = patchSet[num].Reverse
		}
		fmt.Printf("-----> SQL = %s\n", sql)
		_, err = db.Exec(sql)
		if err != nil {
			fmt.Printf("-> ROLLBACK ...\n")
			err2 := tx.Rollback()
			if err2 != nil {
				log.Fatal(err2)
			}
			fmt.Printf("-> ROLLBACK done\n")
			return num, err
		}

		// update the property table
		_, err = db.Exec(`UPDATE property SET value = $1 WHERE key = $2`, num+step, patchKey)
		if err != nil {
			fmt.Printf("-> ROLLBACK ...\n")
			err2 := tx.Rollback()
			if err2 != nil {
				log.Fatal(err2)
			}
			fmt.Printf("-> ROLLBACK done\n")
			return num, err
		}

		// commit
		fmt.Printf("-> COMMIT ...\n")
		err = tx.Commit()
		if err != nil {
			err2 := tx.Rollback()
			if err2 != nil {
				log.Fatal(err2)
			}
			return 0, err
		}
		fmt.Printf("-> COMMIT done\n")
	}

	return level, nil
}
