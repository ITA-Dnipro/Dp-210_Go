INSERT INTO users (id, name, email, ROLE, password_hash)
   VALUES ('00000000-0000-0000-0000-000000000001', 'admin', 'admin@admin.com', 'admin', '$2a$04$bKg2MF8nknGSB8BeSLHd2OEEJiwlIBf.34w02sHiGybRbbWgJ24vC')
ON CONFLICT
   DO NOTHING;

INSERT INTO users (id, name, email, ROLE, password_hash)
   VALUES ('00000000-0000-0000-0000-000000000002', 'operator', 'operator@us.com', 'operator', '$2a$04$t1YjOvExDuBejdc6bDojh./Pbv8dhsfVsz2Nagimzf4ur1Tw5vxl.')
ON CONFLICT
   DO NOTHING;

INSERT INTO users (id, name, email, ROLE, password_hash)
   VALUES ('10000000-0000-0000-0000-000000000001', 'doctor1', 'doctor1@us.com', 'viewer', '$2a$04$4TqnK9D7naJXOzT1ImQhW.0vtohYWkdCJx/xq5rkMLBLmIaWTxQTu')
ON CONFLICT
   DO NOTHING;

INSERT INTO users (id, name, email, ROLE, password_hash)
   VALUES ('10000000-0000-0000-0000-000000000002', 'doctor2', 'doctor2@us.com', 'viewer', '$2a$04$4TqnK9D7naJXOzT1ImQhW.0vtohYWkdCJx/xq5rkMLBLmIaWTxQTu')
ON CONFLICT
   DO NOTHING;

INSERT INTO users (id, name, email, ROLE, password_hash)
   VALUES ('20000000-0000-0000-0000-000000000001', 'user1', 'user1@us.com', 'viewer', '$2a$04$4TqnK9D7naJXOzT1ImQhW.0vtohYWkdCJx/xq5rkMLBLmIaWTxQTu')
ON CONFLICT
   DO NOTHING;

INSERT INTO users (id, name, email, ROLE, password_hash)
   VALUES ('20000000-0000-0000-0000-000000000002', 'user2', 'user2@us.com', 'viewer', '$2a$04$4TqnK9D7naJXOzT1ImQhW.0vtohYWkdCJx/xq5rkMLBLmIaWTxQTu')
ON CONFLICT
   DO NOTHING;
