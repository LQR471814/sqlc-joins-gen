package sqlite

import (
	"sqlc-joins-gen/lib/schema"
	"sqlc-joins-gen/lib/utils"
	"testing"
)

type TestSchema struct {
	Source []byte
	Schema schema.Schema
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
	ratedBy integer,
	rating real,
	primary key (authorId, bookId, ratedBy),
	foreign key (authorId, bookId) references BookAuthorRelevance(authorId, bookId),
	foreign key (ratedBy) references Author(id)
);

create table BookMetadata (
	bookId integer not null primary key,
	pages integer not null,
	rating real,
	foreign key (bookId) references Book(id)
);
`),
		Schema: schema.TESTING_SCHEMAS[0],
	},
	{
		Source: []byte(`
create table User (
    email text not null primary key,
    gpa real not null
);

create table PSCourse (
    name text not null primary key
);

create table PSUserCourse (
    userEmail text not null,
    courseName text not null,
    primary key (userEmail, courseName),
    foreign key (userEmail) references User(email)
        on update cascade
        on delete cascade,
    foreign key (courseName) references PSCourse(name)
        on update cascade
        on delete cascade
);

create table PSUserMeeting (
    userEmail text not null,
    courseName text not null,
    startTime integer not null,
    endTime integer not null,
    primary key (userEmail, courseName, startTime),
    foreign key (userEmail, courseName) references PSUserCourse(userEmail, courseName)
        on update cascade
        on delete cascade
);

create table PSAssignmentType (
    courseName text not null,
    name text not null,
    primary key (name, courseName),
    foreign key (courseName) references PSCourse(name)
        on update cascade
        on delete cascade
);

create table PSAssignment (
    name text not null,
    courseName text not null,
    assignmentTypeName text not null,
    description text,
    duedate integer not null,
    category text not null,
    primary key (name, courseName),
    foreign key (courseName, assignmentTypeName) references PSAssignmentType(courseName, name)
        on update cascade
        on delete cascade
);

create table PSUserAssignment (
    userEmail text not null,
    assignmentName text not null,
    courseName text not null,
    missing integer not null,
    collected integer not null,
    scored real,
    total real,
    primary key (userEmail, assignmentName, courseName),
    foreign key (assignmentName, courseName) references PSAssignment(name, courseName)
        on update cascade
        on delete cascade,
    foreign key (courseName, userEmail) references PSUserCourse(courseName, userEmail)
        on update cascade
        on delete cascade
);

create table MoodleCourse (
    id text not null primary key,
    courseName text not null,
    teacher text,
    zoom text
);

create table MoodleUserCourse (
    courseId text not null,
    userEmail text not null,
    primary key (courseId, userEmail),
    foreign key (courseId) references MoodleCourse(id)
        on update cascade
        on delete cascade,
    foreign key (userEmail) references User(email)
        on update cascade
        on delete cascade
);

create table MoodlePage (
    courseId text not null,
    url text not null,
    content text not null,
    primary key (url, courseId),
    foreign key (courseId) references MoodleCourse(id)
        on update cascade
        on delete cascade
);

create table MoodleAssignment (
    name text not null,
    courseId text not null,
    description text,
    duedate integer not null,
    category text,
    primary key (name, courseId),
    foreign key (courseId) references MoodleCourse(id)
        on update cascade
        on delete cascade
);`,
		),
		Schema: schema.TESTING_SCHEMAS[1],
	},
}

func TestParseSchema(t *testing.T) {
	for _, test := range TESTING_SCHEMAS {
		result, err := ParseSchema([]byte(test.Source))
		if err != nil {
			t.Fatal("failed to construct schema:", err)
		}

		diff := utils.DiffUnordered(test.Schema, result)
		if diff != "" {
			t.Fatal(
				"unexpected differences in expected schema parse:",
				diff,
			)
		}
	}
}
