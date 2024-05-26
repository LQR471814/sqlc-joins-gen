package gen

import (
	"sqlc-joins-gen/lib/querycfg"
	"sqlc-joins-gen/lib/schema"
	"sqlc-joins-gen/lib/utils"
	"testing"
)

func TestIsUniqueFkey(t *testing.T) {
	testCases := []struct {
		schema   schema.Schema
		table    schema.Table
		fkey     schema.ForeignKey
		expected bool
	}{
		{
			schema:   schema.TESTING_SCHEMAS[0],
			table:    schema.TESTING_SCHEMAS[0].Tables[1],
			fkey:     schema.TESTING_SCHEMAS[0].Tables[1].ForeignKeys[0],
			expected: false,
		},
		{
			schema:   schema.TESTING_SCHEMAS[0],
			table:    schema.TESTING_SCHEMAS[0].Tables[3],
			fkey:     schema.TESTING_SCHEMAS[0].Tables[3].ForeignKeys[0],
			expected: false,
		},
		{
			schema:   schema.TESTING_SCHEMAS[0],
			table:    schema.TESTING_SCHEMAS[0].Tables[3],
			fkey:     schema.TESTING_SCHEMAS[0].Tables[3].ForeignKeys[1],
			expected: false,
		},
		{
			schema:   schema.TESTING_SCHEMAS[0],
			table:    schema.TESTING_SCHEMAS[0].Tables[4],
			fkey:     schema.TESTING_SCHEMAS[0].Tables[4].ForeignKeys[0],
			expected: true,
		},
	}

	for _, test := range testCases {
		g := GenManager{Schema: test.schema}
		result := g.isUniqueFkey(test.table, test.fkey)
		diff := utils.DiffUnordered(result, test.expected)
		if diff != "" {
			t.Fatalf(
				"got different result than expected:\n%s\ntable: %s\nfkey target: %d",
				diff, test.table.Name, test.fkey.TargetTable,
			)
		}
	}
}

func TestGetSelectFields(t *testing.T) {
	schemas := schema.TESTING_SCHEMAS
	testCases := []struct {
		schema      schema.Schema
		parentTable string
		query       querycfg.Query
		expected    []string
	}{
		{
			schema:      schemas[0],
			parentTable: "Author",
			query: querycfg.Query{
				Columns: map[string]bool{
					"name": true,
				},
				With: map[string]querycfg.Query{
					"Book": {
						Columns: map[string]bool{
							"id":   true,
							"name": true,
						},
					},
					"BookAuthorRelevanceRating": {
						Columns: map[string]bool{
							"rating": false,
						},
					},
				},
			},
			expected: []string{
				// this Author_id is included because it is a primary key
				"Author_id",
				"Author_name",
				"Book_id",
				"Book_name",
				"BookAuthorRelevanceRating_authorId",
				"BookAuthorRelevanceRating_bookId",
				"BookAuthorRelevanceRating_ratedBy",
			},
		},
	}

	for _, test := range testCases {
		fromSchema := GenManager{Schema: test.schema}
		var fields []SqlSelectField
		fromSchema.getSelectFields(
			test.query,
			test.schema.MustFindTable(test.parentTable),
			&fields,
		)

		result := make([]string, len(fields))
		for i, f := range fields {
			result[i] = f.As
		}

		diff := utils.DiffUnordered(test.expected, result)
		if diff != "" {
			t.Fatalf("different test result than expected:\n%s", diff)
		}
	}
}

func TestGetJoinLine(t *testing.T) {
	schemas := schema.TESTING_SCHEMAS
	testCases := []struct {
		schema   schema.Schema
		source   int
		target   int
		expected SqlJoinLine
	}{
		{
			schema: schemas[0],
			source: 0,
			target: 1,
			expected: SqlJoinLine{
				Table: "Book",
				On: []SqlJoinOn{
					{
						SourceTable: "Book",
						SourceAttr:  "authorId",
						TargetTable: "Author",
						TargetAttr:  "id",
					},
				},
			},
		},
		{
			schema: schemas[0],
			source: 1,
			target: 2,
			expected: SqlJoinLine{
				Table: "BookAuthorRelevance",
				On: []SqlJoinOn{
					{
						SourceTable: "BookAuthorRelevance",
						SourceAttr:  "bookId",
						TargetTable: "Book",
						TargetAttr:  "id",
					},
				},
			},
		},
		{
			schema: schemas[0],
			source: 2,
			target: 3,
			expected: SqlJoinLine{
				Table: "BookAuthorRelevanceRating",
				On: []SqlJoinOn{
					{
						SourceTable: "BookAuthorRelevanceRating",
						SourceAttr:  "authorId",
						TargetTable: "BookAuthorRelevance",
						TargetAttr:  "authorId",
					},
					{
						SourceTable: "BookAuthorRelevanceRating",
						SourceAttr:  "bookId",
						TargetTable: "BookAuthorRelevance",
						TargetAttr:  "bookId",
					},
				},
			},
		},
		{
			schema: schemas[0],
			source: 0,
			target: 3,
			expected: SqlJoinLine{
				Table: "BookAuthorRelevanceRating",
				On: []SqlJoinOn{
					{
						SourceTable: "BookAuthorRelevanceRating",
						SourceAttr:  "ratedBy",
						TargetTable: "Author",
						TargetAttr:  "id",
					},
				},
			},
		},
	}

	for _, test := range testCases {
		generator := GenManager{Schema: test.schema}
		result := generator.getJoinLine(test.source, test.target)
		diff := utils.DiffUnordered(test.expected, result)
		if diff != "" {
			t.Fatalf(
				"unexpected differences in generated join clause:\n%s",
				diff,
			)
		}
	}
}

func TestGetRowDef(t *testing.T) {
	schemas := schema.TESTING_SCHEMAS
	testCases := []struct {
		schema   schema.Schema
		method   querycfg.Method
		expected []PlRowDef
	}{
		{
			schema: schemas[0],
			method: querycfg.Method{
				Table: "Book",
				Query: querycfg.Query{
					Columns: map[string]bool{
						"authorId": true,
					},
					With: map[string]querycfg.Query{
						"Author": {
							Columns: map[string]bool{
								"name": false,
							},
						},
					},
				},
			},
			expected: []PlRowDef{
				{
					DefName:    "Book",
					TableName:  "Book",
					MethodRoot: true,
					Fields: []PlFieldDef{
						{
							Name: "authorId",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name: "id",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name: "Author",
							Type: PlType{
								IsStruct: true,
								Struct:   1,
							},
						},
					},
				},
				{
					DefName:    "Author",
					TableName:  "Author",
					MethodRoot: false,
					Fields: []PlFieldDef{
						{
							Name: "id",
							Type: PlType{
								Primitive: INT,
							},
						},
					},
				},
			},
		},
		{
			schema: schemas[0],
			method: querycfg.Method{
				Table: "Author",
				Name:  "some method",
				Query: querycfg.Query{
					With: map[string]querycfg.Query{
						"Book": {},
					},
				},
			},
			expected: []PlRowDef{
				{
					DefName:    "Author",
					TableName:  "Author",
					MethodRoot: true,
					MethodName: "some method",
					Fields: []PlFieldDef{
						{
							Name: "id",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name: "name",
							Type: PlType{
								Primitive: STRING,
							},
						},
						{
							Name: "Book",
							Type: PlType{
								Array:    true,
								IsStruct: true,
								Struct:   1,
							},
						},
					},
				},
				{
					DefName:    "Book",
					TableName:  "Book",
					MethodRoot: false,
					MethodName: "some method",
					Fields: []PlFieldDef{
						{
							Name: "id",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name: "authorId",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name: "name",
							Type: PlType{
								Primitive: STRING,
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		generator := GenManager{Schema: test.schema}
		var result []PlRowDef
		generator.getRowDefs(test.method, &result)

		diff := utils.DiffUnordered(test.expected, result)

		if diff != "" {
			t.Fatalf(
				"unexpected differences in generated struct def:\n%s",
				diff,
			)
		}
	}
}
