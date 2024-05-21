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
    startTime int not null,
    endTime int not null,
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

create table PSUserAssignment (
    userEmail text not null,
    assignmentName text not null,
    courseName text not null,
    missing int not null,
    collected int not null,
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

create table PSAssignment (
    name text not null,
    courseName text not null,
    assignmentTypeName text not null,
    description text,
    duedate int not null,
    category text not null,
    primary key (name, courseName),
    foreign key (courseName, assignmentTypeName) references PSAssignmentType(courseName, name)
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
    duedate int not null,
    category text,
    primary key (name, courseId),
    foreign key (courseId) references MoodleCourse(id)
        on update cascade
        on delete cascade
);

create table WeightCourse (
    name text not null primary key
);

create table WeightCourseAssignmentType (
    courseName text not null,
    name text not null,
    weight real not null,
    primary key (courseName, name),
    foreign key (courseName) references WeightCourse(name)
        on update cascade
        on delete cascade
);

