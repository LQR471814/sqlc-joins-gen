package types

import "sqlc-joins-gen/lib/utils"

func init() {
	// unordered sorting support for schema.go
	utils.AddCustomSort(func(val *Table) string {
		return val.Name
	})
	utils.AddCustomSort(func(val *Column) string {
		return val.Name
	})
	utils.AddCustomSort(func(val ForeignKey) string {
		return val.TargetTable.Name
	})
	utils.AddCustomSort(func(val ForeignColumn) string {
		return val.SourceColumn.Name
	})

	// unordered sorting support for query.go
	utils.AddCustomSort(func(val Method) string {
		return val.Name
	})
	utils.AddCustomSort(func(val QueryColumn) string {
		return val.Column.Name
	})
	utils.AddCustomSort(func(val QueryOrderBy) string {
		return val.Column.Name
	})
	utils.AddCustomSort(func(val QueryWith) string {
		return val.Table.Name
	})
}

func authorBookSchema() Schema {
	authorId := &Column{
		Name:     "id",
		Type:     INT,
		Nullable: false,
	}
	author := &Table{
		Name: "Author",
		Columns: []*Column{
			authorId,
			{
				Name:     "name",
				Type:     TEXT,
				Nullable: false,
			},
		},
		PrimaryKey:  []*Column{authorId},
		ForeignKeys: nil,
	}

	bookId := &Column{
		Name:     "id",
		Type:     INT,
		Nullable: false,
	}
	bookAuthorId := &Column{
		Name:     "authorId",
		Type:     INT,
		Nullable: false,
	}
	book := &Table{
		Name: "Book",
		Columns: []*Column{
			bookId,
			bookAuthorId,
			{
				Name:     "name",
				Type:     TEXT,
				Nullable: false,
			},
		},
		PrimaryKey: []*Column{bookId},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: author,
				On: []ForeignColumn{
					{
						SourceColumn: bookAuthorId,
						TargetColumn: authorId,
					},
				},
			},
		},
	}

	bookAuthorRelevanceAuthorId := &Column{
		Name:     "authorId",
		Type:     INT,
		Nullable: false,
	}
	bookAuthorRelevanceBookId := &Column{
		Name:     "bookId",
		Type:     INT,
		Nullable: false,
	}
	bookAuthorRelevance := &Table{
		Name: "BookAuthorRelevance",
		Columns: []*Column{
			bookAuthorRelevanceAuthorId,
			bookAuthorRelevanceBookId,
			{
				Name:     "relevance",
				Type:     REAL,
				Nullable: true,
			},
		},
		PrimaryKey: []*Column{bookAuthorRelevanceAuthorId, bookAuthorRelevanceBookId},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: author,
				On: []ForeignColumn{
					{
						SourceColumn: bookAuthorRelevanceAuthorId,
						TargetColumn: authorId,
					},
				},
			},
			{
				TargetTable: book,
				On: []ForeignColumn{
					{
						SourceColumn: bookAuthorRelevanceBookId,
						TargetColumn: bookId,
					},
				},
			},
		},
	}

	bookAuthorRelevanceRatingAuthorId := &Column{
		Name:     "authorId",
		Type:     INT,
		Nullable: false,
	}
	bookAuthorRelevanceRatingBookId := &Column{
		Name:     "bookId",
		Type:     INT,
		Nullable: false,
	}
	bookAuthorRelevanceRatingRatedBy := &Column{
		Name:     "ratedBy",
		Type:     INT,
		Nullable: true,
	}
	bookAuthorRelevanceRating := &Table{
		Name: "BookAuthorRelevanceRating",
		Columns: []*Column{
			bookAuthorRelevanceRatingAuthorId,
			bookAuthorRelevanceRatingBookId,
			bookAuthorRelevanceRatingRatedBy,
			{
				Name:     "rating",
				Type:     REAL,
				Nullable: true,
			},
		},
		PrimaryKey: []*Column{
			bookAuthorRelevanceRatingAuthorId,
			bookAuthorRelevanceRatingBookId,
			bookAuthorRelevanceRatingRatedBy,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: bookAuthorRelevance,
				On: []ForeignColumn{
					{
						SourceColumn: bookAuthorRelevanceRatingAuthorId,
						TargetColumn: bookAuthorRelevanceAuthorId,
					},
					{
						SourceColumn: bookAuthorRelevanceRatingBookId,
						TargetColumn: bookAuthorRelevanceBookId,
					},
				},
			},
			{
				TargetTable: author,
				On: []ForeignColumn{
					{
						SourceColumn: bookAuthorRelevanceRatingRatedBy,
						TargetColumn: authorId,
					},
				},
			},
		},
	}

	bookMetadataBookId := &Column{
		Name:     "bookId",
		Type:     INT,
		Nullable: false,
	}
	bookMetadata := &Table{
		Name: "BookMetadata",
		Columns: []*Column{
			bookMetadataBookId,
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
		PrimaryKey: []*Column{bookMetadataBookId},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: book,
				On: []ForeignColumn{
					{
						SourceColumn: bookMetadataBookId,
						TargetColumn: bookId,
					},
				},
			},
		},
	}

	return Schema{
		Tables: []*Table{
			author,
			book,
			bookAuthorRelevance,
			bookAuthorRelevanceRating,
			bookMetadata,
		},
	}
}

func userDataSchema() Schema {
	userEmail := &Column{
		Name:     "email",
		Type:     TEXT,
		Nullable: false,
	}
	user := &Table{
		Name: "User",
		Columns: []*Column{
			userEmail,
			{
				Name:     "gpa",
				Type:     REAL,
				Nullable: false,
			},
		},
		PrimaryKey: []*Column{userEmail},
	}

	psCourseName := &Column{
		Name:     "name",
		Type:     TEXT,
		Nullable: false,
	}
	psCourse := &Table{
		Name: "PSCourse",
		Columns: []*Column{
			psCourseName,
		},
		PrimaryKey: []*Column{psCourseName},
	}

	psUserCourseUserEmail := &Column{
		Name:     "userEmail",
		Type:     TEXT,
		Nullable: false,
	}
	psUserCourseCourseName := &Column{
		Name:     "courseName",
		Type:     TEXT,
		Nullable: false,
	}
	psUserCourse := &Table{
		Name: "PSUserCourse",
		Columns: []*Column{
			psUserCourseUserEmail,
			psUserCourseCourseName,
		},
		PrimaryKey: []*Column{psUserCourseUserEmail, psUserCourseCourseName},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: user,
				On: []ForeignColumn{
					{
						SourceColumn: psUserCourseUserEmail,
						TargetColumn: userEmail,
					},
				},
			},
			{
				TargetTable: psCourse,
				On: []ForeignColumn{
					{
						SourceColumn: psUserCourseCourseName,
						TargetColumn: psCourseName,
					},
				},
			},
		},
	}
	psUserMeetingUserEmail := &Column{
		Name:     "userEmail",
		Type:     TEXT,
		Nullable: false,
	}
	psUserMeetingCourseName := &Column{
		Name:     "courseName",
		Type:     TEXT,
		Nullable: false,
	}
	psUserMeetingStartTime := &Column{
		Name:     "startTime",
		Type:     INT,
		Nullable: false,
	}
	psUserMeetingEndTime := &Column{
		Name:     "endTime",
		Type:     INT,
		Nullable: false,
	}
	psUserMeeting := &Table{
		Name: "PSUserMeeting",
		Columns: []*Column{
			psUserMeetingUserEmail,
			psUserMeetingCourseName,
			psUserMeetingStartTime,
			psUserMeetingEndTime,
		},
		PrimaryKey: []*Column{
			psUserMeetingUserEmail,
			psUserMeetingCourseName,
			psUserMeetingStartTime,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: psUserCourse,
				On: []ForeignColumn{
					{
						SourceColumn: psUserMeetingUserEmail,
						TargetColumn: psUserCourseUserEmail,
					},
					{
						SourceColumn: psUserMeetingCourseName,
						TargetColumn: psUserCourseCourseName,
					},
				},
			},
		},
	}

	psAssignmentTypeCourseName := &Column{
		Name:     "courseName",
		Type:     TEXT,
		Nullable: false,
	}
	psAssignmentTypeName := &Column{
		Name:     "name",
		Type:     TEXT,
		Nullable: false,
	}
	psAssignmentType := &Table{
		Name: "PSAssignmentType",
		Columns: []*Column{
			psAssignmentTypeCourseName,
			psAssignmentTypeName,
		},
		PrimaryKey: []*Column{
			psAssignmentTypeCourseName,
			psAssignmentTypeName,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: psCourse,
				On: []ForeignColumn{
					{
						SourceColumn: psAssignmentTypeCourseName,
						TargetColumn: psCourseName,
					},
				},
			},
		},
	}

	psAssignmentName := &Column{
		Name:     "name",
		Type:     TEXT,
		Nullable: false,
	}
	psAssignmentCourseName := &Column{
		Name:     "courseName",
		Type:     TEXT,
		Nullable: false,
	}
	psAssignmentAssignmentTypeName := &Column{
		Name:     "assignmentTypeName",
		Type:     TEXT,
		Nullable: false,
	}
	psAssignment := &Table{
		Name: "PSAssignment",
		Columns: []*Column{
			psAssignmentName,
			psAssignmentCourseName,
			psAssignmentAssignmentTypeName,
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
		PrimaryKey: []*Column{
			psAssignmentName,
			psAssignmentCourseName,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: psAssignmentType,
				On: []ForeignColumn{
					{
						SourceColumn: psAssignmentCourseName,
						TargetColumn: psAssignmentTypeCourseName,
					},
					{
						SourceColumn: psAssignmentAssignmentTypeName,
						TargetColumn: psAssignmentTypeName,
					},
				},
			},
		},
	}

	psUserAssignmentUserEmail := &Column{
		Name:     "userEmail",
		Type:     TEXT,
		Nullable: false,
	}
	psUserAssignmentAssignmentName := &Column{
		Name:     "assignmentName",
		Type:     TEXT,
		Nullable: false,
	}
	psUserAssignmentCourseName := &Column{
		Name:     "courseName",
		Type:     TEXT,
		Nullable: false,
	}
	psUserAssignment := &Table{
		Name: "PSUserAssignment",
		Columns: []*Column{
			psUserAssignmentUserEmail,
			psUserAssignmentAssignmentName,
			psUserAssignmentCourseName,
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
		PrimaryKey: []*Column{
			psUserAssignmentUserEmail,
			psUserAssignmentAssignmentName,
			psUserAssignmentCourseName,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: psAssignment,
				On: []ForeignColumn{
					{
						SourceColumn: psUserAssignmentCourseName,
						TargetColumn: psAssignmentCourseName,
					},
					{
						SourceColumn: psUserAssignmentAssignmentName,
						TargetColumn: psAssignmentName,
					},
				},
			},
			{
				TargetTable: psUserCourse,
				On: []ForeignColumn{
					{
						SourceColumn: psUserAssignmentCourseName,
						TargetColumn: psUserCourseCourseName,
					},
					{
						SourceColumn: psUserAssignmentUserEmail,
						TargetColumn: psUserCourseUserEmail,
					},
				},
			},
		},
	}

	moodleCourseId := &Column{
		Name:     "id",
		Type:     TEXT,
		Nullable: false,
	}
	moodleCourse := &Table{
		Name: "MoodleCourse",
		Columns: []*Column{
			moodleCourseId,
			{
				Name:     "courseName",
				Type:     TEXT,
				Nullable: false,
			},
			{
				Name:     "teacher",
				Type:     TEXT,
				Nullable: true,
			},
			{
				Name:     "zoom",
				Type:     TEXT,
				Nullable: true,
			},
		},
		PrimaryKey: []*Column{moodleCourseId},
	}

	moodleUserCourseCourseId := &Column{
		Name:     "courseId",
		Type:     TEXT,
		Nullable: false,
	}
	moodleUserCourseUserEmail := &Column{
		Name:     "userEmail",
		Type:     TEXT,
		Nullable: false,
	}
	moodleUserCourse := &Table{
		Name: "MoodleUserCourse",
		Columns: []*Column{
			moodleUserCourseCourseId,
			moodleUserCourseUserEmail,
		},
		PrimaryKey: []*Column{
			moodleUserCourseCourseId,
			moodleUserCourseUserEmail,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: moodleCourse,
				On: []ForeignColumn{
					{
						SourceColumn: moodleUserCourseCourseId,
						TargetColumn: moodleCourseId,
					},
				},
			},
			{
				TargetTable: user,
				On: []ForeignColumn{
					{
						SourceColumn: moodleUserCourseUserEmail,
						TargetColumn: userEmail,
					},
				},
			},
		},
	}

	moodlePageCourseId := &Column{
		Name:     "courseId",
		Type:     TEXT,
		Nullable: false,
	}
	moodlePageUrl := &Column{
		Name:     "url",
		Type:     TEXT,
		Nullable: false,
	}
	moodlePage := &Table{
		Name: "MoodlePage",
		Columns: []*Column{
			moodlePageCourseId,
			moodlePageUrl,
			{
				Name:     "content",
				Type:     TEXT,
				Nullable: false,
			},
		},
		PrimaryKey: []*Column{
			moodlePageCourseId,
			moodlePageUrl,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: moodleCourse,
				On: []ForeignColumn{
					{
						SourceColumn: moodlePageCourseId,
						TargetColumn: moodleCourseId,
					},
				},
			},
		},
	}

	moodleAssignmentName := &Column{
		Name:     "name",
		Type:     TEXT,
		Nullable: false,
	}
	moodleAssignmentCourseId := &Column{
		Name:     "courseId",
		Type:     TEXT,
		Nullable: false,
	}
	moodleAssignment := &Table{
		Name: "MoodleAssignment",
		Columns: []*Column{
			moodleAssignmentName,
			moodleAssignmentCourseId,
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
				Nullable: true,
			},
		},
		PrimaryKey: []*Column{
			moodleAssignmentName,
			moodleAssignmentCourseId,
		},
		ForeignKeys: []ForeignKey{
			{
				TargetTable: moodleCourse,
				On: []ForeignColumn{
					{
						SourceColumn: moodleAssignmentCourseId,
						TargetColumn: moodleCourseId,
					},
				},
			},
		},
	}

	return Schema{
		Tables: []*Table{
			user,
			psCourse,
			psUserCourse,
			psUserMeeting,
			psAssignmentType,
			psAssignment,
			psUserAssignment,
			moodleCourse,
			moodleUserCourse,
			moodlePage,
			moodleAssignment,
		},
	}
}

var TESTING_SCHEMAS = []Schema{
	authorBookSchema(),
	userDataSchema(),
}

func getUserDataMethod() Method {
	schema := userDataSchema()

	user := schema.MustFindTable("User")
	psUserCourse := schema.MustFindTable("PSUserCourse")
	psUserAssignment := schema.MustFindTable("PSUserAssignment")
	moodleUserCourse := schema.MustFindTable("MoodleUserCourse")
	moodleCourse := schema.MustFindTable("MoodleCourse")

	return Method{
		Name:   "getUserData",
		Table:  user,
		Return: MANY,
		Query: Query{
			Columns: []QueryColumn{
				{
					Column:  user.FindColumn("gpa"),
					Enabled: true,
				},
			},
			With: []QueryWith{
				{
					Table: psUserCourse,
					Query: Query{
						Columns: []QueryColumn{
							{
								Column:  psUserCourse.MustFindColumn("courseName"),
								Enabled: true,
							},
						},
						With: []QueryWith{
							{
								Table: psUserAssignment,
								Query: Query{
									With: []QueryWith{
										{Table: schema.MustFindTable("PSAssignment")},
									},
								},
							},
							{Table: schema.MustFindTable("PSUserMeeting")},
						},
					},
				},
				{
					Table: moodleUserCourse,
					Query: Query{
						Columns: []QueryColumn{
							{
								Column:  moodleUserCourse.FindColumn("courseId"),
								Enabled: true,
							},
							{
								Column:  moodleUserCourse.FindColumn("userEmail"),
								Enabled: false,
							},
						},
						With: []QueryWith{
							{
								Table: moodleCourse,
								Query: Query{
									With: []QueryWith{
										{Table: schema.MustFindTable("MoodlePage")},
										{Table: schema.MustFindTable("MoodleAssignment")},
									},
								},
							},
						},
					},
				},
			},
			Where: "User.email = ?",
			OrderBy: []QueryOrderBy{
				{
					Column:  user.FindColumn("gpa"),
					OrderBy: ASC,
				},
			},
		},
	}
}

var TESTING_METHODS = []Method{
	getUserDataMethod(),
}
