package main

import (
	"slices"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateQuerySelect(t *testing.T) {
	schemas := TESTING_SCHEMAS
	testCases := []struct {
		schema      Schema
		parentTable string
		query       CompositeQueryClause
		expected    []string
	}{
		{
			schema:      schemas[0].Schema,
			parentTable: "Author",
			query: CompositeQueryClause{
				Columns: map[string]bool{
					"name": true,
				},
				With: map[string]CompositeQueryClause{
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
		generator := compositeQueryGenerator{schema: test.schema}
		result := generator.generateSelectFields(test.query, test.parentTable)

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
