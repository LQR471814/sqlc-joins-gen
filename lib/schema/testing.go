package schema

var TESTING_SCHEMAS = []Schema{
	{
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
				ForeignKeys: nil,
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
						Nullable: true,
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
			{
				Name: "BookMetadata",
				Columns: []Column{
					{
						Name:     "bookId",
						Type:     INT,
						Nullable: false,
					},
					{
						Name:     "pages",
						Type:     INT,
						Nullable: false,
					},
					{
						Name:     "rating",
						Type:     REAL,
						Nullable: true,
					},
				},
				PrimaryKey: []int{0},
				ForeignKeys: []ForeignKey{
					{
						TargetTable: 1,
						On: []ForeignColumn{
							{
								SourceColumn: 0,
								TargetColumn: 0,
							},
						},
					},
				},
			},
		},
	},
	{
		Tables: []Table{
			{
				Name: "User",
				Columns: []Column{
					{
						Name:     "email",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "gpa",
						Type:     REAL,
						Nullable: false,
					},
				},
				PrimaryKey: []int{0},
			},
			{
				Name: "PSCourse",
				Columns: []Column{
					{
						Name:     "name",
						Type:     TEXT,
						Nullable: false,
					},
				},
				PrimaryKey: []int{0},
			},
			{
				Name: "PSUserCourse",
				Columns: []Column{
					{
						Name:     "userEmail",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "courseName",
						Type:     TEXT,
						Nullable: false,
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
				Name: "PSUserMeeting",
				Columns: []Column{
					{
						Name:     "userEmail",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "courseName",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "startTime",
						Type:     INT,
						Nullable: false,
					},
					{
						Name:     "endTIme",
						Type:     INT,
						Nullable: false,
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
				},
			},
			{
				Name: "PSAssignmentType",
				Columns: []Column{
					{
						Name:     "courseName",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "name",
						Type:     TEXT,
						Nullable: false,
					},
				},
				PrimaryKey: []int{0, 1},
				ForeignKeys: []ForeignKey{
					{
						TargetTable: 1,
						On: []ForeignColumn{
							{
								SourceColumn: 0,
								TargetColumn: 0,
							},
						},
					},
				},
			},
			{
				Name: "PSAssignment",
				Columns: []Column{
					{
						Name:     "name",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "courseName",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "assignmentTypeName",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "description",
						Type:     TEXT,
						Nullable: true,
					},
					{
						Name:     "duedate",
						Type:     INT,
						Nullable: false,
					},
					{
						Name:     "category",
						Type:     TEXT,
						Nullable: false,
					},
				},
				PrimaryKey: []int{0, 1},
				ForeignKeys: []ForeignKey{
					{
						TargetTable: 4,
						On: []ForeignColumn{
							{
								SourceColumn: 1,
								TargetColumn: 0,
							},
							{
								SourceColumn: 0,
								TargetColumn: 1,
							},
						},
					},
				},
			},
			{
				Name: "PSUserAssignment",
				Columns: []Column{
					{
						Name:     "userEmail",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "assignmentName",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "courseName",
						Type:     TEXT,
						Nullable: false,
					},
					{
						Name:     "missing",
						Type:     INT,
						Nullable: false,
					},
					{
						Name:     "collected",
						Type:     INT,
						Nullable: false,
					},
					{
						Name:     "scored",
						Type:     REAL,
						Nullable: true,
					},
					{
						Name:     "total",
						Type:     REAL,
						Nullable: true,
					},
				},
				PrimaryKey: []int{0, 1, 2},
				ForeignKeys: []ForeignKey{
					{
						TargetTable: 5,
						On: []ForeignColumn{
							{
								SourceColumn: 1,
								TargetColumn: 0,
							},
							{
								SourceColumn: 2,
								TargetColumn: 1,
							},
						},
					},
				},
			},
		},
	},
}
