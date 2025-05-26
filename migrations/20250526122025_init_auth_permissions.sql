-- +goose Up
-- 5. Seed permissions
INSERT INTO auth.permissions (code, description) VALUES
  ('leads.create', 'Create leads'),
  ('leads.edit', 'Edit leads'),
  ('leads.view', 'View leads'),
  ('leads.delete', 'Delete leads'),
  ('projects.create', 'Create projects'),
  ('projects.edit', 'Edit projects'),
  ('projects.view', 'View projects'),
  ('invoices.create', 'Create invoices'),
  ('invoices.send', 'Send invoices'),
  ('invoices.view', 'View invoices'),
  ('users.manage', 'Manage users'),
  ('users.invite', 'Invite users');

-- 6. Seed permission groups
INSERT INTO auth.permission_groups (name, description) VALUES
  ('Admin', 'Full access to everything'),
  ('User', 'Standard user access'),
  ('Viewer', 'Read-only access');

-- 7. Assign all permissions to Admin
INSERT INTO auth.group_permissions (group_id, permission_id)
SELECT pg.id, p.id
FROM auth.permission_groups pg, auth.permissions p
WHERE pg.name = 'Admin';

-- 8. Assign selected permissions to User
INSERT INTO auth.group_permissions (group_id, permission_id)
SELECT pg.id, p.id
FROM auth.permission_groups pg
JOIN auth.permissions p ON p.code IN (
  'leads.create', 'leads.edit', 'leads.view',
  'projects.create', 'projects.edit', 'projects.view',
  'invoices.create', 'invoices.view'
)
WHERE pg.name = 'User';

-- 9. Assign limited view-only permissions to Viewer
INSERT INTO auth.group_permissions (group_id, permission_id)
SELECT pg.id, p.id
FROM auth.permission_groups pg
JOIN auth.permissions p ON p.code IN (
  'leads.view', 'projects.view', 'invoices.view'
)
WHERE pg.name = 'Viewer';

-- +goose Down
ALTER TABLE auth.users DROP COLUMN group_id;
DROP TABLE IF EXISTS auth.group_permissions;
DROP TABLE IF EXISTS auth.permission_groups;
DROP TABLE IF EXISTS auth.permissions;

