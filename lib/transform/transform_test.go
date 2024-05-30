package transform

import (
	"sqlc-joins-gen/lib/outputs"
	"sqlc-joins-gen/lib/types"
	"sqlc-joins-gen/lib/utils"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestIsUniqueFkey(t *testing.T) {
	testCases := []struct {
		schema   types.Schema
		table    *types.Table
		fkey     types.ForeignKey
		expected bool
	}{
		{
			schema:   types.TESTING_SCHEMAS[0],
			table:    types.TESTING_SCHEMAS[0].Tables[1],
			fkey:     types.TESTING_SCHEMAS[0].Tables[1].ForeignKeys[0],
			expected: false,
		},
		{
			schema:   types.TESTING_SCHEMAS[0],
			table:    types.TESTING_SCHEMAS[0].Tables[3],
			fkey:     types.TESTING_SCHEMAS[0].Tables[3].ForeignKeys[0],
			expected: false,
		},
		{
			schema:   types.TESTING_SCHEMAS[0],
			table:    types.TESTING_SCHEMAS[0].Tables[3],
			fkey:     types.TESTING_SCHEMAS[0].Tables[3].ForeignKeys[1],
			expected: false,
		},
		{
			schema:   types.TESTING_SCHEMAS[0],
			table:    types.TESTING_SCHEMAS[0].Tables[4],
			fkey:     types.TESTING_SCHEMAS[0].Tables[4].ForeignKeys[0],
			expected: true,
		},
		{
			schema:   types.TESTING_SCHEMAS[1],
			table:    types.TESTING_SCHEMAS[1].Tables[8],
			fkey:     types.TESTING_SCHEMAS[1].Tables[8].ForeignKeys[0],
			expected: false,
		},
		{
			schema:   types.TESTING_SCHEMAS[1],
			table:    types.TESTING_SCHEMAS[1].Tables[5],
			fkey:     types.TESTING_SCHEMAS[1].Tables[5].ForeignKeys[0],
			expected: false,
		},
	}

	for _, test := range testCases {
		g := FromSchema{Schema: test.schema}
		result := g.isUniqueFkey(test.table, test.fkey)
		diff := utils.DiffUnordered(test.expected, result)
		if diff != "" {
			t.Fatalf(
				"got different result than expected:\n%s\ntable: %s\nfkey target: %p",
				diff, test.table.Name, test.fkey.TargetTable,
			)
		}
	}
}

func TestGetSelectFields(t *testing.T) {
	schemas := types.TESTING_SCHEMAS

	authorTable := schemas[0].MustFindTable("Author")
	bookTable := schemas[0].MustFindTable("Book")
	bookAuthorRelevanceRatingTable := schemas[0].MustFindTable("BookAuthorRelevanceRating")

	testCases := []struct {
		schema   types.Schema
		table    *types.Table
		query    types.Query
		expected []string
	}{
		{
			schema: schemas[0],
			table:  authorTable,
			query: types.Query{
				Columns: []types.QueryColumn{
					{
						Column:  authorTable.FindColumn("name"),
						Enabled: true,
					},
				},
				With: []types.QueryWith{
					{
						Table: bookTable,
						Query: types.Query{
							Columns: []types.QueryColumn{
								{
									Column:  bookTable.FindColumn("id"),
									Enabled: true,
								},
								{
									Column:  bookTable.FindColumn("name"),
									Enabled: true,
								},
							},
						},
					},
					{
						Table: bookAuthorRelevanceRatingTable,
						Query: types.Query{
							Columns: []types.QueryColumn{
								{
									Column:  bookAuthorRelevanceRatingTable.FindColumn("rating"),
									Enabled: false,
								},
							},
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
			schema: schemas[1],
			table:  schemas[1].MustFindTable("User"),
			query:  types.TESTING_METHODS[0].Query,
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
		fromSchema := FromSchema{Schema: test.schema}
		var fields []outputs.SqlSelectField
		fromSchema.getSelectFields(test.query, test.table, &fields)

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
	schemas := types.TESTING_SCHEMAS
	testCases := []struct {
		schema   types.Schema
		source   *types.Table
		target   *types.Table
		expected outputs.SqlJoinLine
	}{
		{
			schema: schemas[0],
			source: schemas[0].MustFindTable("Author"),
			target: schemas[0].MustFindTable("Book"),
			expected: outputs.SqlJoinLine{
				Table: "Book",
				On: []outputs.SqlJoinOn{
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
			source: schemas[0].MustFindTable("Book"),
			target: schemas[0].MustFindTable("BookAuthorRelevance"),
			expected: outputs.SqlJoinLine{
				Table: "BookAuthorRelevance",
				On: []outputs.SqlJoinOn{
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
			source: schemas[0].MustFindTable("BookAuthorRelevance"),
			target: schemas[0].MustFindTable("BookAuthorRelevanceRating"),
			expected: outputs.SqlJoinLine{
				Table: "BookAuthorRelevanceRating",
				On: []outputs.SqlJoinOn{
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
			source: schemas[0].MustFindTable("Author"),
			target: schemas[0].MustFindTable("BookAuthorRelevanceRating"),
			expected: outputs.SqlJoinLine{
				Table: "BookAuthorRelevanceRating",
				On: []outputs.SqlJoinOn{
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
		generator := FromSchema{Schema: test.schema}
		result, err := generator.getJoinLine(test.source, test.target)
		if err != nil {
			t.Fatal(err)
		}
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
	schemas := types.TESTING_SCHEMAS
	testCases := []struct {
		schema   types.Schema
		method   types.Method
		expected []*outputs.PlRowDef
	}{
		{
			schema: schemas[0],
			method: types.Method{
				Table: schemas[0].MustFindTable("Book"),
				Name:  "getBooksAndAuthor",
				Query: types.Query{
					Columns: []types.QueryColumn{
						{
							Column:  schemas[0].MustFindTable("Book").FindColumn("authorId"),
							Enabled: true,
						},
					},
					With: []types.QueryWith{
						{
							Table: schemas[0].MustFindTable("Author"),
							Query: types.Query{
								Columns: []types.QueryColumn{
									{
										Column:  schemas[0].MustFindTable("Author").FindColumn("name"),
										Enabled: false,
									},
								},
							},
						},
					},
				},
			},
			expected: []*outputs.PlRowDef{
				// index: 0
				{
					DefName:   "getBooksAndAuthor",
					TableName: "Book",
					Fields: []*outputs.PlFieldDef{
						{
							Name:           "authorId",
							TableFieldName: "authorId",
							Type: outputs.PlType{
								Primitive: outputs.INT,
							},
						},
						{
							Name:           "id",
							TableFieldName: "id",
							Type: outputs.PlType{
								Primitive: outputs.INT,
							},
						},
						{
							Name: "Author",
							Type: outputs.PlType{
								Array: false,
							},
							IsRowDef: true,
							RowDef:   1,
						},
					},
				},
				// index: 1
				{
					DefName:   "getBooksAndAuthor0",
					TableName: "Author",
					Fields: []*outputs.PlFieldDef{
						{
							Name:           "id",
							TableFieldName: "id",
							Type: outputs.PlType{
								Primitive: outputs.INT,
							},
						},
					},
				},
			},
		},
		{
			schema: schemas[0],
			method: types.Method{
				Table: schemas[0].MustFindTable("Author"),
				Name:  "getAuthorAndBooks",
				Query: types.Query{
					With: []types.QueryWith{
						{
							Table: schemas[0].MustFindTable("Book"),
						},
					},
				},
			},
			expected: []*outputs.PlRowDef{
				{
					DefName:    "getAuthorAndBooks",
					TableName:  "Author",
					PrimaryKey: []*outputs.PlFieldDef{},
					Fields: []*outputs.PlFieldDef{
						{
							Name:           "id",
							TableFieldName: "id",
							Type: outputs.PlType{
								Primitive: outputs.INT,
							},
						},
						{
							Name:           "name",
							TableFieldName: "name",
							Type: outputs.PlType{
								Primitive: outputs.STRING,
							},
						},
						{
							Name: "Book",
							Type: outputs.PlType{
								Array: true,
							},
							IsRowDef: true,
							RowDef:   1,
						},
					},
				},
				{
					DefName:    "getAuthorAndBooks0",
					TableName:  "Book",
					PrimaryKey: []*outputs.PlFieldDef{},
					Fields: []*outputs.PlFieldDef{
						{
							Name:           "id",
							TableFieldName: "id",
							Type: outputs.PlType{
								Primitive: outputs.INT,
							},
						},
						{
							Name:           "authorId",
							TableFieldName: "authorId",
							Type: outputs.PlType{
								Primitive: outputs.INT,
							},
						},
						{
							Name:           "name",
							TableFieldName: "name",
							Type: outputs.PlType{
								Primitive: outputs.STRING,
							},
						},
					},
				},
			},
		},
		{
			schema: types.TESTING_SCHEMAS[1],
			method: types.TESTING_METHODS[0],
			expected: []*outputs.PlRowDef{
				// index: 0
				{
					TableName: "User",
					DefName:   "getUserData",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "gpa",
							Name:           "gpa",
							Type:           outputs.PlType{Primitive: outputs.FLOAT},
						},
						{
							TableFieldName: "email",
							Name:           "email",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							Name: "PSUserCourse",
							Type: outputs.PlType{
								Array: true,
							},
							IsRowDef: true,
							RowDef:   1,
						},
						{
							Name: "MoodleUserCourse",
							Type: outputs.PlType{
								Array: true,
							},
							IsRowDef: true,
							RowDef:   2,
						},
					},
				},
				// index: 1
				{
					TableName: "PSUserCourse",
					DefName:   "getUserData0",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							Name: "PSUserAssignment",
							Type: outputs.PlType{
								Array: true,
							},
							IsRowDef: true,
							RowDef:   3,
						},
						{
							Name:     "PSUserMeeting",
							Type:     outputs.PlType{Array: true},
							IsRowDef: true,
							RowDef:   4,
						},
					},
				},
				// index: 2
				{
					TableName: "MoodleUserCourse",
					DefName:   "getUserData1",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "courseId",
							Name:           "courseId",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							Name:     "MoodleCourse",
							IsRowDef: true,
							RowDef:   5,
						},
					},
				},
				// index: 3
				{
					TableName: "PSUserAssignment",
					DefName:   "getUserData00",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "assignmentName",
							Name:           "assignmentName",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "missing",
							Name:           "missing",
							Type:           outputs.PlType{Primitive: outputs.INT},
						},
						{
							TableFieldName: "collected",
							Name:           "collected",
							Type:           outputs.PlType{Primitive: outputs.INT},
						},
						{
							TableFieldName: "scored",
							Name:           "scored",
							Type:           outputs.PlType{Primitive: outputs.FLOAT, Nullable: true},
						},
						{
							TableFieldName: "total",
							Name:           "total",
							Type:           outputs.PlType{Primitive: outputs.FLOAT, Nullable: true},
						},
						{
							Name:     "PSAssignment",
							IsRowDef: true,
							RowDef:   6,
						},
					},
				},
				// index: 4
				{
					TableName: "PSUserMeeting",
					DefName:   "getUserData01",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "userEmail",
							Name:           "userEmail",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "startTime",
							Name:           "startTime",
							Type:           outputs.PlType{Primitive: outputs.INT},
						},
						{
							TableFieldName: "endTime",
							Name:           "endTime",
							Type:           outputs.PlType{Primitive: outputs.INT},
						},
					},
				},
				// index: 5
				{
					TableName: "MoodleCourse",
					DefName:   "getUserData10",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "id",
							Name:           "id",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "teacher",
							Name:           "teacher",
							Type:           outputs.PlType{Primitive: outputs.STRING, Nullable: true},
						},
						{
							TableFieldName: "zoom",
							Name:           "zoom",
							Type:           outputs.PlType{Primitive: outputs.STRING, Nullable: true},
						},
						{
							Name: "MoodlePage",
							Type: outputs.PlType{
								Array: true,
							},
							IsRowDef: true,
							RowDef:   7,
						},
						{
							Name: "MoodleAssignment",
							Type: outputs.PlType{
								Array: true,
							},
							IsRowDef: true,
							RowDef:   8,
						},
					},
				},
				// index: 6
				{
					TableName: "PSAssignment",
					DefName:   "getUserData000",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "name",
							Name:           "name",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "courseName",
							Name:           "courseName",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "assignmentTypeName",
							Name:           "assignmentTypeName",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "description",
							Name:           "description",
							Type:           outputs.PlType{Primitive: outputs.STRING, Nullable: true},
						},
						{
							TableFieldName: "duedate",
							Name:           "duedate",
							Type:           outputs.PlType{Primitive: outputs.INT},
						},
						{
							TableFieldName: "category",
							Name:           "category",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
					},
				},
				// index: 7
				{
					TableName: "MoodlePage",
					DefName:   "getUserData100",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "courseId",
							Name:           "courseId",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "url",
							Name:           "url",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "content",
							Name:           "content",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
					},
				},
				// index: 8
				{
					TableName: "MoodleAssignment",
					DefName:   "getUserData101",
					Fields: []*outputs.PlFieldDef{
						{
							TableFieldName: "name",
							Name:           "name",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "courseId",
							Name:           "courseId",
							Type:           outputs.PlType{Primitive: outputs.STRING},
						},
						{
							TableFieldName: "description",
							Name:           "description",
							Type:           outputs.PlType{Primitive: outputs.STRING, Nullable: true},
						},
						{
							TableFieldName: "duedate",
							Name:           "duedate",
							Type:           outputs.PlType{Primitive: outputs.INT},
						},
						{
							TableFieldName: "category",
							Name:           "category",
							Type:           outputs.PlType{Primitive: outputs.STRING, Nullable: true},
						},
					},
				},
			},
		},
	}

	for _, test := range testCases {
		generator := FromSchema{Schema: test.schema}
		var result []*outputs.PlRowDef
		generator.GetRowDefs(test.method, &result)

		diff := utils.DiffUnordered(
			test.expected, result,
			cmpopts.IgnoreFields(outputs.PlRowDef{}, "PrimaryKey"),
			cmpopts.IgnoreFields(outputs.PlRowDef{}, "Parent"),
			cmpopts.IgnoreFields(outputs.PlRowDef{}, "ParentField"),
		)

		if diff != "" {
			t.Fatalf(
				"unexpected differences in generated struct def:\n%s",
				diff,
			)
		}
	}
}
