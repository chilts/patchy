test: simple

simple:
	echo "-------------------------------------------------------------------------------"
	PGUSER=frisk PGDATABASE=frisk go run ../bin/patchy.go --level 1 --patch-dir simple
	echo "-------------------------------------------------------------------------------"
	PGUSER=frisk PGDATABASE=frisk go run ../bin/patchy.go --level 2 --patch-dir simple
	echo "-------------------------------------------------------------------------------"
	PGUSER=frisk PGDATABASE=frisk go run ../bin/patchy.go --level 0 --patch-dir simple
	echo "-------------------------------------------------------------------------------"
	PGUSER=frisk PGDATABASE=frisk go run ../bin/patchy.go --level 2 --patch-dir simple
	echo "-------------------------------------------------------------------------------"
	PGUSER=frisk PGDATABASE=frisk go run ../bin/patchy.go --level 0 --patch-dir simple
	echo "-------------------------------------------------------------------------------"

.PHONY: simple
