INSERT INTO apps (id, name, secret)
VALUES (2, 'testtest', 'testtest-secret')
ON CONFLICT DO NOTHING;