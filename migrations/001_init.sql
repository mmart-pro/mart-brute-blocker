-- +goose Up
create table if not exists white_list
(
    subnet inet not null,
    constraint ck_white_list_subnet exclude using gist (subnet inet_ops with &&),
    constraint pk_white_list_subnet primary key (subnet)
);

create index if not exists ix_white_list_subnet on white_list using gist (subnet inet_ops);

create table if not exists black_list (
	subnet inet not null,
    constraint ck_black_list_subnet exclude using gist (subnet inet_ops with &&),
    constraint pk_black_list_subnet primary key (subnet)
);

create index if not exists ix_black_list_subnet on black_list using gist (subnet inet_ops);
