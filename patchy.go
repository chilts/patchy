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
	"path"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func Hi() {
	fmt.Printf("Hi\n")
}

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

func Patch(db *sql.DB, level int, opts *Options) (int, error) {
	if opts == nil {
		opts = &Options{}
	}

	// read the dir list
	if opts.Dir == "" {
		opts.Dir = "."
	}
	files, err := ioutil.ReadDir(opts.Dir)
	if err != nil {
		return 0, err
	}

	// check the files are what we expect
	patchSet := make(map[int]*patch)
	for _, file := range files {
		fmt.Printf("file=%s\n", file.Name())

		direction := ""
		if strings.HasPrefix(file.Name(), "forward-") {
			fmt.Printf("- forward\n")
			direction = "forward"
		} else if strings.HasPrefix(file.Name(), "reverse-") {
			fmt.Printf("- reverse\n")
			direction = "reverse"
		} else {
			fmt.Printf("- unknown\n")
			continue
		}

		// what is the patch level of this filename
		patchLevel := file.Name()
		patchLevel = strings.TrimPrefix(patchLevel, direction+"-")
		patchLevel = strings.TrimSuffix(patchLevel, ".sql")

		// check this is a number
		n, err := strconv.Atoi(patchLevel)
		if err != nil {
			return 0, err
		}

		// read this file in
		if _, ok := patchSet[n]; ok {
			fmt.Printf("patch exists\n")
		} else {
			fmt.Printf("no such patch in patchset yet\n")
			patchSet[n] = &patch{}
		}

		filename := path.Join(opts.Dir, file.Name())
		sql, err := ioutil.ReadFile(filename)
		if err != nil {
			return 0, err
		}

		p := patchSet[n]
		if direction == "forward" {
			p.Forward = string(sql)
		}
		if direction == "reverse" {
			p.Reverse = string(sql)
		}
		fmt.Printf("p=%#v\n", p)
	}

	fmt.Printf("PatchSet=%v\n", patchSet)

	// remember some things when we check the property table
	// var patchRowExists bool
	var currentLevel int

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
                table_name = 'property'
        );
    `

	queryExists, err := db.Query(sqlPropertyTableExists)
	if err != nil {
		return 0, err
	}

	propertyTableExists, err := oneRowOneColBool(queryExists)
	if err != nil {
		if err == sql.ErrNoRows {

		} else {
			// no table exists probably
		}
	}

	// figure out the current patch level
	key := 21
	rows, err := db.Query("SELECT value FROM property WHERE key = $1", key)
	defer rows.Close()

	// loop through all rows (though we expect one at the most)
	for rows.Next() {
		propertyTableExists = true
		err = rows.Scan(&currentLevel)
		if err != nil {
			return 0, rows.Err()
		}
	}

	// get any error encountered during iteration
	if rows.Err() != nil {
		return 0, rows.Err()
	}

	if propertyTableExists {

	} else {

	}

	// readPatchDir,
	// getCurrentPatch,
	// checkAllPatchFilesExist,
	// begin,

	// nextPatch,
	// writeCurrentLevel,
	// commit,

	return 1, nil
}
