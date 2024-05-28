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
		{
			schema:   schema.TESTING_SCHEMAS[1],
			table:    schema.TESTING_SCHEMAS[1].Tables[8],
			fkey:     schema.TESTING_SCHEMAS[1].Tables[8].ForeignKeys[0],
			expected: false,
		},
		{
			schema:   schema.TESTING_SCHEMAS[1],
			table:    schema.TESTING_SCHEMAS[1].Tables[5],
			fkey:     schema.TESTING_SCHEMAS[1].Tables[5].ForeignKeys[0],
			expected: false,
		},
	}

	for _, test := range testCases {
		g := GenManager{Schema: test.schema}
		result := g.isUniqueFkey(test.table, test.fkey)
		diff := utils.DiffUnordered(test.expected, result)
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
		{
			schema:      schemas[1],
			parentTable: "User",
			query:       querycfg.TESTING_METHODS[0].Query,
			expected: []string{
				"User_email",
				"User_gpa",
				"PSUserCourse_courseName",
				"PSUserCourse_userEmail",
				"PSUserAssignment_userEmail",
				"PSUserAssignment_assignmentName",
				"PSUserAssignment_courseName",
				"PSUserAssignment_missing",
				"PSUserAssignment_collected",
				"PSUserAssignment_scored",
				"PSUserAssignment_total",
				"PSAssignment_name",
				"PSAssignment_courseName",
				"PSAssignment_assignmentTypeName",
				"PSAssignment_description",
				"PSAssignment_duedate",
				"PSAssignment_category",
				"PSUserMeeting_userEmail",
				"PSUserMeeting_courseName",
				"PSUserMeeting_startTime",
				"PSUserMeeting_endTime",
				"MoodleUserCourse_courseId",
				"MoodleUserCourse_userEmail",
				"MoodleCourse_id",
				"MoodleCourse_courseName",
				"MoodleCourse_teacher",
				"MoodleCourse_zoom",
				"MoodlePage_courseId",
				"MoodlePage_url",
				"MoodlePage_content",
				"MoodleAssignment_name",
				"MoodleAssignment_courseId",
				"MoodleAssignment_description",
				"MoodleAssignment_duedate",
				"MoodleAssignment_category",
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
				Name:  "getBooksAndAuthor",
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
				// index: 0
				{
					DefName:   "getBooksAndAuthor",
					TableName: "Book",
					Fields: []PlFieldDef{
						{
							Name:           "authorId",
							TableFieldName: "authorId",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name:           "id",
							TableFieldName: "id",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name: "Author",
							Type: PlType{
								IsRowDef: true,
								Array:    false,
								RowDef:   1,
							},
						},
					},
				},
				// index: 1
				{
					DefName:   "getBooksAndAuthor0",
					TableName: "Author",
					Fields: []PlFieldDef{
						{
							Name:           "id",
							TableFieldName: "id",
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
				Name:  "getAuthorAndBooks",
				Query: querycfg.Query{
					With: map[string]querycfg.Query{
						"Book": {},
					},
				},
			},
			expected: []PlRowDef{
				{
					DefName:   "getAuthorAndBooks",
					TableName: "Author",
					Fields: []PlFieldDef{
						{
							Name:           "id",
							TableFieldName: "id",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name:           "name",
							TableFieldName: "name",
							Type: PlType{
								Primitive: STRING,
							},
						},
						{
							Name: "Book",
							Type: PlType{
								Array:    true,
								IsRowDef: true,
								RowDef:   1,
							},
						},
					},
				},
				{
					DefName:   "getAuthorAndBooks0",
					TableName: "Book",
					Fields: []PlFieldDef{
						{
							Name:           "id",
							TableFieldName: "id",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name:           "authorId",
							TableFieldName: "authorId",
							Type: PlType{
								Primitive: INT,
							},
						},
						{
							Name:           "name",
							TableFieldName: "name",
							Type: PlType{
								Primitive: STRING,
							},
						},
					},
				},
			},
		},
		{
			schema: schema.TESTING_SCHEMAS[1],
			method: querycfg.TESTING_METHODS[0],
			expected: []PlRowDef{
				// index: 0
				{
					TableName: "User",
					DefName:   "getUserData",
					Fields: []PlFieldDef{
						{
							TableFieldName: "gpa",
							Name:           "gpa",
							Type:           PlType{Primitive: FLOAT},
						},
						{
							TableFieldName: "email",
							Name:           "email",
							Type:           PlType{Primitive: STRING},
						},
						{
							Name: "PSUserCourse",
							Type: PlType{
								IsRowDef: true,
								Array:    true,
								RowDef:   1,
							},
						},
						{
							Name: "MoodleUserCourse",
							Type: PlType{
								IsRowDef: true,
								Array:    true,
								RowDef:   2,
							},
						},
					},
				},
				// index: 1
				{
					TableName: "PSUserCourse",
					DefName:   "getUserData0",
					Fields: []PlFieldDef{
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           PlType{Primitive: STRING},
						},
						{
							Name: "PSUserAssignment",
							Type: PlType{
								IsRowDef: true,
								Array:    true,
								RowDef:   3,
							},
						},
						{
							Name: "PSUserMeeting",
							Type: PlType{
								IsRowDef: true,
								Array:    true,
								RowDef:   4,
							},
						},
					},
				},
				// index: 2
				{
					TableName: "MoodleUserCourse",
					DefName:   "getUserData1",
					Fields: []PlFieldDef{
						{
							TableFieldName: "courseId",
							Name:           "courseId",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           PlType{Primitive: STRING},
						},
						{
							Name: "MoodleCourse",
							Type: PlType{
								IsRowDef: true,
								RowDef:   5,
							},
						},
					},
				},
				// index: 3
				{
					TableName: "PSUserAssignment",
					DefName:   "getUserData00",
					Fields: []PlFieldDef{
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "assignmentName",
							Name:           "assignmentName",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "missing",
							Name:           "missing",
							Type:           PlType{Primitive: INT},
						},
						{
							TableFieldName: "collected",
							Name:           "collected",
							Type:           PlType{Primitive: INT},
						},
						{
							TableFieldName: "scored",
							Name:           "scored",
							Type:           PlType{Primitive: FLOAT, Nullable: true},
						},
						{
							TableFieldName: "total",
							Name:           "total",
							Type:           PlType{Primitive: FLOAT, Nullable: true},
						},
						{
							Name: "PSAssignment",
							Type: PlType{
								IsRowDef: true,
								RowDef:   6,
							},
						},
					},
				},
				// index: 4
				{
					TableName: "PSUserMeeting",
					DefName:   "getUserData01",
					Fields: []PlFieldDef{
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "startTime",
							Name:           "startTime",
							Type:           PlType{Primitive: INT},
						},
						{
							TableFieldName: "endTime",
							Name:           "endTime",
							Type:           PlType{Primitive: INT},
						},
					},
				},
				// index: 5
				{
					TableName: "MoodleCourse",
					DefName:   "getUserData10",
					Fields: []PlFieldDef{
						{
							TableFieldName: "id",
							Name:           "id",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "teacher",
							Name:           "teacher",
							Type:           PlType{Primitive: STRING, Nullable: true},
						},
						{
							TableFieldName: "zoom",
							Name:           "zoom",
							Type:           PlType{Primitive: INT, Nullable: true},
						},
						{
							Name: "MoodlePage",
							Type: PlType{
								IsRowDef: true,
								Array:    true,
								RowDef:   7,
							},
						},
						{
							Name: "MoodleAssignment",
							Type: PlType{
								IsRowDef: true,
								Array:    true,
								RowDef:   8,
							},
						},
					},
				},
				// index: 6
				{
					TableName: "PSAssignment",
					DefName:   "getUserData000",
					Fields: []PlFieldDef{
						{
							TableFieldName: "name",
							Name:           "name",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "assignmentTypeName",
							Name:           "assignmentTypeName",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "description",
							Name:           "description",
							Type:           PlType{Primitive: STRING, Nullable: true},
						},
						{
							TableFieldName: "duedate",
							Name:           "duedate",
							Type:           PlType{Primitive: INT},
						},
						{
							TableFieldName: "category",
							Name:           "category",
							Type:           PlType{Primitive: STRING},
						},
					},
				},
				// index: 7
				{
					TableName: "MoodlePage",
					DefName:   "getUserData100",
					Fields: []PlFieldDef{
						{
							TableFieldName: "courseId",
							Name:           "courseId",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "url",
							Name:           "url",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "content",
							Name:           "content",
							Type:           PlType{Primitive: STRING},
						},
					},
				},
				// index: 8
				{
					TableName: "MoodleAssignment",
					DefName:   "getUserData101",
					Fields: []PlFieldDef{
						{
							TableFieldName: "name",
							Name:           "name",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "courseId",
							Name:           "courseId",
							Type:           PlType{Primitive: STRING},
						},
						{
							TableFieldName: "description",
							Name:           "description",
							Type:           PlType{Primitive: STRING, Nullable: true},
						},
						{
							TableFieldName: "duedate",
							Name:           "duedate",
							Type:           PlType{Primitive: INT},
						},
						{
							TableFieldName: "category",
							Name:           "category",
							Type:           PlType{Primitive: STRING, Nullable: true},
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
