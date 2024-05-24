package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestSchema struct {
	Source []byte
	Schema Schema
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
		Schema: Schema{
			Tables: []Table{
				{
					Name: "Author",
					Columns: []Column{
						{
							Name:     "id",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "name",
							Type:     TEXT,
							Nullable: false,
						},
					},
					PrimaryKey:  []int{0},
					ForeignKeys: []ForeignKey{},
				},
				{
					Name: "Book",
					Columns: []Column{
						{
							Name:     "id",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "authorId",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "name",
							Type:     TEXT,
							Nullable: false,
						},
					},
					PrimaryKey: []int{0},
					ForeignKeys: []ForeignKey{
						{
							TargetTable: 0,
							On: []ForeignColumn{
								{
									SourceColumn: 1,
									TargetColumn: 0,
								},
							},
						},
					},
				},
				{
					Name: "BookAuthorRelevance",
					Columns: []Column{
						{
							Name:     "authorId",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "bookId",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "relevance",
							Type:     REAL,
							Nullable: true,
						},
					},
					PrimaryKey: []int{0, 1},
					ForeignKeys: []ForeignKey{
						{
							TargetTable: 0,
							On: []ForeignColumn{
								{
									SourceColumn: 0,
									TargetColumn: 0,
								},
							},
						},
						{
							TargetTable: 1,
							On: []ForeignColumn{
								{
									SourceColumn: 1,
									TargetColumn: 0,
								},
							},
						},
					},
				},
				{
					Name: "BookAuthorRelevanceRating",
					Columns: []Column{
						{
							Name:     "authorId",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "bookId",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "ratedBy",
							Type:     INT,
							Nullable: false,
						},
						{
							Name:     "rating",
							Type:     REAL,
							Nullable: true,
						},
					},
					PrimaryKey: []int{0, 1, 2},
					ForeignKeys: []ForeignKey{
						{
							TargetTable: 2,
							On: []ForeignColumn{
								{
									SourceColumn: 0,
									TargetColumn: 0,
								},
								{
									SourceColumn: 1,
									TargetColumn: 1,
								},
							},
						},
						{
							TargetTable: 0,
							On: []ForeignColumn{
								{
									SourceColumn: 2,
									TargetColumn: 0,
								},
							},
						},
					},
				},
			},
		},
	},
}

func TestNewSchema(t *testing.T) {
	for _, test := range TESTING_SCHEMAS {
		result, err := NewSchema([]byte(test.Source))
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
