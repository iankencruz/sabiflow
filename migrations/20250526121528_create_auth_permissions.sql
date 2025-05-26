
-- +goose Up
-- Create permissions table
CREATE TABLE auth.permissions (
    id SERIAL PRIMARY KEY,
    code TEXT UNIQUE NOT NULL,
    description TEXT
);

-- Create permission groups (e.g. Admin, User, Viewer)
CREATE TABLE auth.permission_groups (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT
);

-- Map group to permission
CREATE TABLE auth.group_permissions (
    group_id INTEGER NOT NULL REFERENCES auth.permission_groups(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES auth.permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (group_id, permission_id)
);

-- Add group_id to users table
ALTER TABLE auth.users
ADD COLUMN group_id INTEGER REFERENCES auth.permission_groups(id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE auth.users DROP COLUMN group_id;
DROP TABLE IF EXISTS auth.group_permissions;
DROP TABLE IF EXISTS auth.permission_groups;
DROP TABLE IF EXISTS auth.permissions;

