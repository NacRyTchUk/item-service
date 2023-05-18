BEGIN;

create table if not exists campaigns
(
    id   serial,
    name varchar not null
);

create index campaigns_id_index
    on campaigns (id);

alter table campaigns
    add primary key (id);




create table if not exists items
(
    id          serial,
    campaign_id integer,
    name        varchar   not null,
    description varchar,
    priority    integer   not null,
    removed     bool      not null default false ,
    created_at  timestamp without time zone default (now() at time zone 'utc') not null,
    constraint items_pk
        primary key (id, campaign_id)
);

create function max_priority()
    returns integer as $$ select COALESCE((SELECT MAX(priority) FROM items), 0) + 1$$ language sql;

alter table items alter column priority set default max_priority();

insert into campaigns (id, name) values (default, 'Первая запись');

COMMIT;
