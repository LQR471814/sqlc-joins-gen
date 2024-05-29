-- name: CreateUser :exec
insert into User(email, gpa) values (?, ?);
