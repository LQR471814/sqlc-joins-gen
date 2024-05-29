-- name: CreateAuthor :exec
insert into Author(id, name) values (?, ?);
