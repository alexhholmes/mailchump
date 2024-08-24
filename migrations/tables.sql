create table newsletters
(
    id              uuid default gen_random_uuid()         not null,
    owner_id        uuid                                   not null,
    title           varchar(300)                           not null,
    slug            varchar(300)                           not null,
    description     varchar(900)                           not null,
    created         timestamp with time zone default now() not null,
    updated         timestamp with time zone default now() not null,
    post_count      bigint                   default 0     not null,
    hidden          boolean                  default false not null,
    deleted         boolean                  default false not null,
    recovery_window timestamp with time zone,
    constraint newsletters_pk
        primary key (id),
    constraint newsletters_pk_2
        unique (title)
);

alter table newsletters
    owner to username;

create index newsletters_owner_id_index
    on newsletters (owner_id);

create table users
(
    id         uuid default gen_random_uuid() not null,
    user_name  varchar(30)                    not null,
    password   varchar(128)                   not null,
    first_name varchar(30)                    not null,
    last_name  varchar(30)                    not null,
    email      varchar(255)                   not null,
    constraint users_pk
        primary key (id),
    constraint users_pk_2
        unique (user_name)
);

alter table users
    owner to username;

create table posts
(
    id              uuid default gen_random_uuid()         not null,
    owner_id        uuid                                   not null,
    title           varchar(300)                           not null,
    slug            varchar(300)                           not null,
    description     varchar(900)                           not null,
    content         text                                   not null,
    created         timestamp with time zone default now() not null,
    updated         timestamp with time zone default now() not null,
    hidden          boolean                  default false not null,
    deleted         boolean                  default false not null,
    recovery_window timestamp,
    constraint posts_pk
        primary key (id)
);

alter table posts
    owner to username;

create table authors
(
    id      bigserial,
    user_id uuid not null,
    post_id uuid not null,
    constraint authors_pk
        primary key (id),
    constraint authors_users_id_fk
        foreign key (user_id) references users,
    constraint authors_posts_id_fk
        foreign key (post_id) references posts
);

alter table authors
    owner to username;

create index authors_user_id_post_id_index
    on authors (user_id, post_id);

create table subscriptions
(
    id            bigserial,
    user_id       uuid not null,
    newsletter_id uuid not null,
    constraint subscriptions_pk
        primary key (id),
    constraint subscriptions_users_id_fk
        foreign key (user_id) references users,
    constraint subscriptions_newsletters_id_fk
        foreign key (newsletter_id) references newsletters
);

alter table subscriptions
    owner to username;

create index subscriptions_user_id_newsletter_id_index
    on subscriptions (user_id, newsletter_id);


