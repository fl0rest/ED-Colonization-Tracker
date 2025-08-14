-- name: ListEvents :many
select *
from events
order by time desc
;

-- name: AddDepotEvent :exec
INSERT INTO depotEvents (completion, time, marketId, raw_text)
VALUES (:completion, :time, :marketId, :raw_text);

-- name: AddDockEvent :exec
INSERT INTO dockEvents (time, marketId, systemName, stationName)
VALUES (:time, :marketId, :systemName, :stationName);

-- name: AddEvent :exec
INSERT INTO events (time, completion, marketId, systemName, stationName, raw_resources)
VALUES (:time, :completion, :marketId, :systemName, :stationName, :raw_resources);

-- name: GetLatestEvent :one
select *
from events
order by time desc
limit 1
;

-- name: UpsertResource :exec
insert into resources (id, eventId, name, required, provided, diff, payment, time)
values (?, ?, ?, ?, ?, ?, ?, ?) on conflict(id) do update set
  eventId = excluded.eventId,
  required = excluded.required,
  provided = excluded.provided,
  diff = excluded.diff,
  payment = excluded.payment,
  time = excluded.time
  where excluded.provided != provided;

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

-- name: GetInaraId :one
select id
from resourceids
where name like '%' ||:query || '%'
;
