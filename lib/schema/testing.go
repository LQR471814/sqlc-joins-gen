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
}
