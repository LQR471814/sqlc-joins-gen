package gen

import (
	"slices"
	"sqlc-joins-gen/lib/querycfg"
	"sqlc-joins-gen/lib/schema"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetSelectFields(t *testing.T) {
	schemas := schema.TESTING_SCHEMAS
	testCases := []struct {
		schema      schema.Schema
		parentTable string
		qclause     querycfg.Clause
		expected    []SqlSelectField
	}{
		{
			schema:      schemas[0],
			parentTable: "Author",
			qclause: querycfg.Clause{
				Columns: map[string]bool{
					"name": true,
				},
				With: map[string]querycfg.Clause{
					"Book": {
						Columns: map[string]bool{
							"id":   true,
							"name": true,
						},
					},
					"BookAuthorRelevanceRating": {},
				},
			},
			expected: []SqlSelectField{
				SqlSelectField{
					Table: "Author",
					Attr:  "name",
					As:    "Author_name",
				},
				SqlSelectField{
					Table: "Book",
					Attr:  "id",
					As:    "Book_id",
				},
				SqlSelectField{
					Table: "Book",
					Attr:  "name",
					As:    "Book_name",
				},
				SqlSelectField{
					Table: "BookAuthorRelevanceRating",
					Attr:  "authorId",
					As:    "BookAuthorRelevanceRating_authorId",
				},
				SqlSelectField{
					Table: "BookAuthorRelevanceRating",
					Attr:  "bookId",
					As:    "BookAuthorRelevanceRating_bookId",
				},
				SqlSelectField{
					Table: "BookAuthorRelevanceRating",
					Attr:  "ratedBy",
					As:    "BookAuthorRelevanceRating_ratedBy",
				},
				SqlSelectField{
					Table: "BookAuthorRelevanceRating",
					Attr:  "rating",
					As:    "BookAuthorRelevanceRating_rating",
				},
			},
		},
	}

	for _, test := range testCases {
		fromSchema := FromSchema{Schema: test.schema}
		fields := []SqlSelectField{}
		fromSchema.getSelectFields(
			test.qclause,
			test.schema.MustFindTable(test.parentTable),
			&fields,
		)

		result := make([]string, len(fields))
		for i, f := range fields {
			result[i] = f.As
		}
		slices.Sort(result)

		diff := cmp.Diff(processed, test.expected)
		if diff != "" {
			t.Fatalf("different test result than expected:\n%s", diff)
		}
	}
}

func TestGenerateJoinClause(t *testing.T) {
	schemas := TESTING_SCHEMAS
	testCases := []struct {
		schema   Schema
		source   int
		target   int
		expected string
	}{
		{
			schema:   schemas[0].Schema,
			source:   0,
			target:   1,
			expected: "inner join Book on Book.authorId = Author.id",
		},
		{
			schema:   schemas[0].Schema,
			source:   1,
			target:   2,
			expected: "inner join BookAuthorRelevance on BookAuthorRelevance.bookId = Book.id",
		},
		{
			schema:   schemas[0].Schema,
			source:   2,
			target:   3,
			expected: "inner join BookAuthorRelevanceRating on BookAuthorRelevanceRating.authorId = BookAuthorRelevance.authorId and BookAuthorRelevanceRating.bookId = BookAuthorRelevance.bookId",
		},
		{
			schema:   schemas[0].Schema,
			source:   0,
			target:   3,
			expected: "inner join BookAuthorRelevanceRating on BookAuthorRelevanceRating.ratedBy = Author.id",
		},
	}
	for _, test := range testCases {
		generator := compositeQueryGenerator{schema: test.schema}
		result := generator.generateJoinClause(test.source, test.target)
		if result != test.expected {
			t.Fatalf(
				"unexpected differences in generated join clause:\n- expect: %s\n- got:    %s\n",
				test.expected,
				result,
			)
		}
	}
}
