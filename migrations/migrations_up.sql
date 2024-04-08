create table if not exists currency (
    id serial primary key,
    currency_name varchar(255) not null,
    price float not null,
    price_min float not null,
    price_max float not null,
    changes_per_hour float not null,
    last_update time(0) default now()
);

create unique index currency_name_index on currency(currency_name);