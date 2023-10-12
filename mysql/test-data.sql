insert into user_info (user_uuid,user_name,created_at) values
  ('b4d88c57-6514-4775-8f2c-2ee22d2e391e','narumi',current_timestamp),
  ('2559570b-83d0-4be2-aa73-9da770350809','yuka',current_timestamp),
  ('08a70fca-c9f2-43c8-b77f-1a3fe1dedd24','tsubasa',current_timestamp),
  ('71482dba-028b-459b-92e1-b7eeaf6f4247','ruka',current_timestamp);

insert into client_certificate (certificate_uuid,certificate_fingerprint,user_uuid,last_accessed,verify) values
  ('48ed0faf-7366-4f8a-8d7e-a4769f1889fa','+VvO9Cju1dEf3DxCD5ez7IK3GhzY2rB30pspfQ==','b4d88c57-6514-4775-8f2c-2ee22d2e391e',current_timestamp,true),
  ('ac6d7f1b-cdc6-47b8-9dd2-6550fc1936c8','k64o7diksk0zqoqtmg2u5ok49ag0548lghntyg==','2559570b-83d0-4be2-aa73-9da770350809',current_timestamp,true),
  ('a425ded8-84dc-4dc4-81ce-51e4918a8a44','8e6bnlju+i0ag9pjioezccrayoyhup8t9sijcw==','08a70fca-c9f2-43c8-b77f-1a3fe1dedd24',current_timestamp,true),
  ('485b78c1-2a3a-42b8-a111-4e5887712e20','xnmm2qpnqrrphndiznmxonutfhrdlswtws6e5w==','71482dba-028b-459b-92e1-b7eeaf6f4247',current_timestamp,true);

insert operation_history (certificate_uuid,operation_type,operation_result,error_code,operation_time) values
  ('48ed0faf-7366-4f8a-8d7e-a4769f1889fa',0,true,null,current_timestamp),
  ('ac6d7f1b-cdc6-47b8-9dd2-6550fc1936c8',1,true,null,current_timestamp),
  ('a425ded8-84dc-4dc4-81ce-51e4918a8a44',2,false,'999',current_timestamp),
  ('485b78c1-2a3a-42b8-a111-4e5887712e20',0,true,null,current_timestamp);

