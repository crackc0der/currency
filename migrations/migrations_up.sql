create table if not exists currency (
    id serial primary key,
    currency_name varchar(255) not null,
    price float not null,
    price_min float not null,
    price_max float not null,
    changes_per_hour float not null
);

create unique index currency_name_index on currency(currency_name);