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

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	"github.com/chilts/patchy"
)

var level = flag.Int("level", 0, "patch to this level")
var patchDir = flag.String("patch-dir", ".", "directory containing the patch files")
var propertyTable = flag.String("property-table", "property", "key/val table to hold the current patch level")

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	if *level < 0 {
		fmt.Printf("Provide a patch level (--level)")
		return
	}

	if *patchDir == "" {
		fmt.Printf("Provide a patch level (--patch-dir)")
		return
	}

	fmt.Printf("Patching to Level = %d\n", *level)
	fmt.Printf("Patch Directory   = %s\n", *patchDir)
	fmt.Printf("Property Table    = %s\n", *propertyTable)

	// open the connection to the database
	// ToDo: allow all of the same Postgres options from the command line here
	// db, err := sql.Open("postgres", "user=zentype dbname=zentype")
	db, err := sql.Open("postgres", "")
	check(err)

	// patch this database
	opts := patchy.Options{
		Dir:           *patchDir,
		PropertyTable: *propertyTable,
	}
	newLevel, err := patchy.Patch(db, *level, &opts)
	check(err)

	fmt.Printf("bin/patchy: Database patched to level %d\n", newLevel)
}
