-- name: ListEvents :many
select *
from events
order by time desc
;

-- name: AddEvent :one
INSERT INTO events (raw_text, time)
VALUES (:raw_text, :time)
RETURNING id, raw_text, time;

-- name: ListResources :many
select *
from resources
order by diff desc
;

-- name: FindResourceId :one
select id
from resources
where name =:name
;

-- name: FindResourceName :one
select name
from resources
where id =:id
;

-- name: ListResource :one
select *
from resources
where name like '%' ||:query || '%' or id like '%' ||:query || '%'
;
