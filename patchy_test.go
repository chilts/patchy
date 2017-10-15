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
	"testing"

	_ "github.com/lib/pq"
)

func TestSimple(t *testing.T) {
	db, err := sql.Open("postgres", "user=patchy dbname=patchy")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// opts for this test
	opts := Options{
		Dir: "example/simple",
	}

	// test a simple forward patch
	level, err := Patch(db, 2, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if level != 2 {
		t.Fatalf("Level expected %q, but got %q", 2, level)
	}

	// test a simple patch
	level, err = Patch(db, 0, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if level != 0 {
		t.Fatalf("Level expected %v, but got %v", 0, level)
	}

	// test a simple patch
	level, err = Patch(db, 1, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if level != 1 {
		t.Fatalf("Level expected %v, but got %v", 1, level)
	}

	// test a simple patch
	level, err = Patch(db, 2, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if level != 2 {
		t.Fatalf("Level expected %v, but got %v", 2, level)
	}

	// test a simple patch
	level, err = Patch(db, 1, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if level != 1 {
		t.Fatalf("Level expected %v, but got %v", 1, level)
	}

	// test a simple patch
	level, err = Patch(db, 0, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if level != 0 {
		t.Fatalf("Level expected %v, but got %v", 0, level)
	}

}
