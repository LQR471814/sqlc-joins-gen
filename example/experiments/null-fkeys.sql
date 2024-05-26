create table Author (
    id int not null primary key,
    name text not null
);

create table Book (
    id int not null primary key,
    authorId int,
    name text not null,
    foreign key (authorId) references Author(id)
);

-- this is valid, nullable fields can reference and join on
-- not-null fields as foreign keys.
