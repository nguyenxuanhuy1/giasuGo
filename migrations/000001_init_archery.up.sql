create table users (
   id         serial primary key,
   username   text not null unique,
   password   text not null,
   role       text not null default 'user',
   avatar     text,
   locked     boolean default false,
   coin       integer default 0,
   created_at timestamptz default now()
);
create table images (
   id            bigserial primary key,
   image_url     text not null,
   blur_url      text,
   tiny_blur_url text,
   public_id     text not null,
   image_type    text not null,
   owner_id      integer,
   created_at    timestamptz default now(),
   constraint fk_images_user foreign key ( owner_id )
      references users ( id )
         on delete cascade
);
create table posts (
   id          bigserial primary key,
   image_id    bigint
      references images ( id )
         on delete set null,
   name        varchar(255) not null,
   description varchar(255),
   topic       varchar(100) not null,
   prompt      text,
   hot_level   smallint default 0 check ( hot_level between 0 and 9 ),
   hot_at      timestamptz,
   created_at  timestamptz default now(),
   updated_at  timestamptz default now()
);

create index idx_posts_topic_created on
   posts (
      topic,
      created_at
   desc );

create index idx_posts_hot on
   posts (
      hot_level
   desc,
      hot_at
   desc );

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