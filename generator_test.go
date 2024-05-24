package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateQuerySelect(t *testing.T) {
	schemas := TestingSchemas()
	testCases := []struct {
		schema      Schema
		parentTable string
		query       TableQuery
		expected    []string
	}{
		{
			schema:      schemas[0].Schema,
			parentTable: "Author",
			query: TableQuery{
				Columns: map[string]bool{
					"name": true,
				},
				With: map[string]TableQuery{
					"Book": {
						Columns: map[string]bool{
							"id":   true,
							"name": true,
						},
					},
					"BookAuthorRelevanceRating": {},
				},
			},
			expected: []string{
				"Author.name as Author_name",
				"Book.id as Book_id",
				"Book.name as Book_name",
				"BookAuthorRelevanceRating.authorId as BookAuthorRelevanceRating_authorId",
				"BookAuthorRelevanceRating.bookId as BookAuthorRelevanceRating_bookId",
				"BookAuthorRelevanceRating.ratedBy as BookAuthorRelevanceRating_ratedBy",
				"BookAuthorRelevanceRating.rating as BookAuthorRelevanceRating_rating",
			},
		},
	}

	for _, test := range testCases {
		generator := JoinsGenerator{
			Schema: test.schema,
			Joins: JoinsList{
				JoinQueryDef{
					Name:   "TestQuery",
					Return: "one",
					Table:  test.parentTable,
					Query:  test.query,
				},
			},
		}
		result := generator.generateQuerySelect(test.query, test.parentTable)

		processed := strings.Split(result, "\n")
		for i, e := range processed {
			processed[i] = strings.Trim(e, " \n\t,")
		}
		if processed[len(processed)-1] == "" {
			processed = processed[:len(processed)-1]
		}

		slices.Sort(processed)
		slices.Sort(test.expected)

		diff := cmp.Diff(processed, test.expected)
		if diff != "" {
			t.Fatalf("different test result than expected:\n%s", diff)
		}
	}
}
