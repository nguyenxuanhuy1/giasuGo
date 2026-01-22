create table users (
   id         serial primary key,
   username   text,
   role       text not null default 'user',
   avatar     text,
   locked     boolean default false,
   coin       integer default 0,
   google_id  text unique,  -- ID từ Google (sub claim)
   email      text unique,  -- Email từ Google
   created_at timestamptz default now()
);

-- Index để tìm kiếm
create index idx_users_google_id on
   users (
      google_id
   );
create index idx_users_email on
   users (
      email
   );

create table exam_sets (
   id         bigserial primary key,
   name       varchar(255) not null,        -- tên đề
   created_by integer not null
      references users ( id )
         on delete cascade,
   is_public  boolean default false,         -- PUBLIC / PRIVATE

   created_at timestamptz default now()
);

create table exam_questions (
   id             bigserial primary key,
   exam_set_id    bigint not null
      references exam_sets ( id )
         on delete cascade,
   content        text not null,
   type           text not null,              -- single | multiple | essay
   level          text default 'question',
   options        jsonb,                       -- FE gửi
   correct_answer jsonb,                       -- FE gửi

   order_no       integer not null,
   parent_order   integer,
   created_at     timestamptz default now()
);

create table exam_attempts (
   id           bigserial primary key,
   exam_set_id  bigint not null
      references exam_sets ( id )
         on delete cascade,
   user_id      integer not null
      references users ( id )
         on delete cascade,
   started_at   timestamptz default now(),
   submitted_at timestamptz
);

create table exam_answers (
   id          bigserial primary key,
   attempt_id  bigint not null
      references exam_attempts ( id )
         on delete cascade,
   question_id bigint not null
      references exam_questions ( id )
         on delete cascade,
   answer      jsonb not null                -- đáp án user chọn
);