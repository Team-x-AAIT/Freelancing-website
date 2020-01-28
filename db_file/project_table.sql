-- CREATE TABLE projects (
--   id varchar(12) primary key,
--   title varchar(200) not null,
--   description varchar(1000) not null,
--   details varchar(5000) not null,
--   category varchar(100) not null,
--   subcategory varchar(100) not null,
--   budget float,
--   worktype int,
--   closed boolean default false,
--   created_at timestamp
-- );

-- CREATE TABLE attached_files(
--     pid varchar(12),
--     name varchar(50) unique,
--     foreign key (pid) references projects(id)
-- );


-- CREATE TABLE categories (
--   id int primary key auto_increment,
--   name varchar(100) not null unique
-- );

-- CREATE TABLE subcategories (
--   cid int,
--   id int primary key auto_increment,
--   name varchar(100) not null unique,
--   FOREIGN KEY (cid) REFERENCES categories(id)
-- );


-- insert into categories values('1','mobile developer');
-- insert into subcategories values('1','1','front-end');
-- insert into subcategories values('1','2','back-end');


-- CREATE TABLE application_table(
--     pid varchar(12),
--     applicant_uid varchar(12),
--     proposal varchar(2000),
--     hired boolean default false,
--     FOREIGN KEY (pid) REFERENCES projects(id),
--     FOREIGN KEY (applicant_uid) REFERENCES users(uid)
-- );


-- CREATE TABLE application_history_table(
--     pid varchar(12),
--     applicant_uid varchar(12),
--     proposal varchar(2000),
--     hired boolean default false,
--     seen boolean default false,
--     applied_at timestamp
-- );


-- CREATE TABLE user_project_table (
--     uid varchar(12),
--     pid varchar(12),
--     FOREIGN KEY (uid) REFERENCES users(uid),
--     FOREIGN KEY (pid) REFERENCES projects(id)
-- );



-- create table match_tag_table (
--     uid varchar(12),
--     category varchar(100),
--     subcategory varchar(100),
--     worktype int,
--     foreign key (uid) references users(uid),
--     foreign key (category) references categories(name),
--     foreign key (subcategory) references subcategories(name)
-- )


-- link to the database ER diagram - https://dbdiagram.io/d/5e19c36394d9ab14375a1ddc