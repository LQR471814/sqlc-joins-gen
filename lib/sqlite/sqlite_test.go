package sqlite

import (
	"sqlc-joins-gen/lib/schema"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestSchema struct {
	Source []byte
	Schema schema.Schema
}

var TESTING_SCHEMAS = []TestSchema{
	{
		Source: []byte(`
create table Author (
	id integer not null primary key autoincrement,
	name text not null
);

create table Book (
	id integer not null primary key autoincrement,
	authorId integer not null,
	name text not null,
	foreign key (authorId) references Author(id)
);

-- authorId is not necessarily the same author as the author of the book
create table BookAuthorRelevance (
	authorId integer not null,
	bookId integer not null,
	relevance real,
	primary key (authorId, bookId),
	foreign key (authorId) references Author(id),
	foreign key (bookId) references Book(id)
);

create table BookAuthorRelevanceRating (
	authorId integer not null,
	bookId integer not null,
	ratedBy integer not null,
	rating real,
	primary key (authorId, bookId, ratedBy),
	foreign key (authorId, bookId) references BookAuthorRelevance(authorId, bookId),
	foreign key (ratedBy) references Author(id)
);
`),
		Schema: schema.TESTING_SCHEMAS[0],
	},
}

func TestParseSchema(t *testing.T) {
	for _, test := range TESTING_SCHEMAS {
		result, err := ParseSchema([]byte(test.Source))
		if err != nil {
			t.Fatal("failed to construct schema:", err)
		}

		diff := cmp.Diff(test.Schema, result)
		if diff != "" {
			t.Fatal(
				"unexpected differences in expected schema parse:",
				diff,
			)
		}
	}
}
