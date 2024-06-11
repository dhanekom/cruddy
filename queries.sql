-- name: GetPerson :one
select 
    email
  , first_name
  , id
  , last_name
from person
where
  id = ?
limit 1;

-- name: ListPerson :many
select *
from person;

-- name: CreatePerson :execresult
insert into person (
    email
  , first_name
  , id
  , last_name
) values (
    ?
  , ?
  , ?
  , ?
);

-- name: UpdatePerson :exec
update person
set
    email = ?
  , first_name = ?
  , id = ?
  , last_name = ?
where
  id = ?;

-- name: DeletePerson :exec
delete from person
where
  id = ?;

