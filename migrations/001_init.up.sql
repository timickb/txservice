create table accounts (
  id varchar not null
      constraint accounts_pk primary key,
  balance bigint not null default 0
);

create table transactions (
  id varchar not null
      constraint transactions_pk
          primary key,
  date date not null default now(),
  sender_id varchar not null
      constraint tr_sender_id_fk
          references accounts,
  receiver_id varchar not null
      constraint tr_receiver_id_fk
          references accounts,
  amount bigint not null
);