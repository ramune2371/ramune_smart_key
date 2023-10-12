create table if not exists smart_key.user_info(
  user_uuid varchar(36) not null primary key,
  user_name varchar(100) not null,
  created_at timestamp not null
);

create table if not exists smart_key.client_certificate (
 certificate_uuid varchar(36) not null primary key,
 certificate_fingerprint varchar(40) not null,
 user_uuid varchar(36) not null,
 last_accessed timestamp not null,
 verify boolean not null,
 foreign key (user_uuid) references user_info(user_uuid) on delete cascade
);

create table if not exists smart_key.operation_history(
  operation_id int not null primary key auto_increment,
  certificate_uuid varchar(36) not null,
  operation_type int not null,
  operation_result boolean not null,
  error_code varchar(3),
  operation_time timestamp not null,
  foreign key (certificate_uuid) references client_certificate(certificate_uuid) on delete cascade
);


