create table Author (
    id int not null primary key,
    name text not null
);

create table Book (
    id int not null primary key,
    authorId int not null,
    name text not null,
    foreign key (authorId) references Author(id)
);

create table AuthorCard (
    authorId int not null primary key,
    age int,
    sign text,
    foreign key (authorId) references Author(id)
);

create table AuthorExtraCard (
    authorId int primary key,
    age int,
    sign text,
    foreign key (authorId) references Author(id)
);

-- authorId is not unique therefore
-- one author may have many books
select * from Book
inner join Author on Book.authorId = Author.id;

-- it is the same the other way around
select * from Author
inner join Book on Author.id = Book.authorId;

-- one Author may only have one AuthorCard because
-- the authorId field is unique
select * from Author
inner join AuthorCard on AuthorCard.authorId = Author.id;

-- if you're using inner join, there will exist no
-- situation in which you will have a nullable table 
-- so `type Author = { Book: { ... } | undefined }`
-- is impossible (in inner join)

-- the algorithm for determining current-to-n relationship is as follows
-- 1. if there is a non-unique or composite foreign key on the target table referencing the current it is current-to-many
-- 2. if there is a unique foreign key on the target table referencing the current it is current-to-one
-- 3. if there is no foreign key on the target table referencing the current it is current-to-one
