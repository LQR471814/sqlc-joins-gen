create table Author (
    id integer not null primary key,
    name text not null,
    age integer,
    alive integer not null
);

create table Book (
    id integer not null primary key,
    authorId integer not null,
    name text not null,
    foreign key (authorId) references Author(id)
);
