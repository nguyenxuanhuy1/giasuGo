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
   
-- 1. subjects
create table subjects (
   id   serial primary key,
   code varchar(20) unique not null,
   name varchar(100) not null
);

-- 2. exams
create table exams (
   id          bigserial primary key,
   user_id     integer not null
      references users ( id ),
   subject_id  integer not null
      references subjects ( id ),
   title       varchar(255),
   ai_raw_json jsonb not null,
   created_at  timestamptz default now()
);

create index idx_exams_user on
   exams (
      user_id
   );

create index idx_exams_subject on
   exams (
      subject_id
   );

-- 3. questions
create table questions (
   id            bigserial primary key,
   exam_id       bigint not null
      references exams ( id )
         on delete cascade,
   content       text not null,
   question_type text not null check ( question_type in ( 'single',
                                                          'multiple',
                                                          'essay' ) ),
   explanation   text,
   order_no      integer
);

create index idx_questions_exam on
   questions (
      exam_id
   );

-- 4. answers
create table answers (
   id          bigserial primary key,
   question_id bigint not null
      references questions ( id )
         on delete cascade,
   content     text not null,
   is_correct  boolean default false
);

create index idx_answers_question_correct on
   answers (
      question_id,
      is_correct
   );

-- 5. exam_attempts
create table exam_attempts (
   id          bigserial primary key,
   exam_id     bigint not null
      references exams ( id ),
   user_id     integer not null
      references users ( id ),
   started_at  timestamptz default now(),
   finished_at timestamptz,
   score       numeric(5,2)
);

create index idx_attempts_user on
   exam_attempts (
      user_id
   );

create index idx_attempts_user_exam on
   exam_attempts (
      user_id,
      exam_id
   );