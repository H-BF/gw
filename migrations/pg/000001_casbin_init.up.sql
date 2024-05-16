create extension if not exists "uuid-ossp";;

create table if not exists casbin_rule
(
    id uuid default gen_random_uuid() primary key,
    ptype text not null,
    v0 text not null,
    v1 text not null,
    v2 text,
    v3 text,
    v4 text,
    v5 text
);

create index if not exists prtype_idx on casbin_rule(ptype);
create index if not exists v0_idx on casbin_rule(v0);
create index if not exists v1_idx on casbin_rule(v1);
create index if not exists v2_idx on casbin_rule(v2);

insert into casbin_rule (ptype, v0, v1, v2) values
    ('p', 'alice', '*', 'admin'),
    ('p', 'bob', 'bob-res', 'owner'),
    ('p', 'bob', '*', 'client'),
    ('p', 'clare', 'clare-res', 'owner'),
    ('p', 'clare', '*', 'client'),
    ('g', 'admin', 'read', ''),
    ('g', 'admin', 'edit', ''),
    ('g', 'admin', 'reference', ''),
    ('g', 'owner', 'read', ''),
    ('g', 'owner', 'edit', ''),
    ('g', 'owner', 'reference', ''),
    ('g', 'client', 'read', '');