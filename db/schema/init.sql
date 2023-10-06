-- Create tables for the new models if they don't exist
CREATE TABLE IF NOT EXISTS id_uuid_entity (
  ID UUID PRIMARY KEY,
  UUID UUID DEFAULT UUID_generate_v4() NOT NULL
);

CREATE TABLE IF NOT EXISTS task_metadata_entity (
  UUID UUID DEFAULT UUID_generate_v4() NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  task_type VARCHAR(255) NOT NULL,
  assignee VARCHAR(255),
  status VARCHAR(255) NOT NULL,
  priority VARCHAR(255) NOT NULL,
  risk VARCHAR(255) NOT NULL,
  tags TEXT[],
  dates_uuid UUID,
  story_point_uuid UUID,
  iterations_uuid UUID,
  relations_uuid UUID,
  task_entity_uuid UUID
);

CREATE TABLE IF NOT EXISTS task_content_entity (
  UUID UUID DEFAULT UUID_generate_v4() NOT NULL PRIMARY KEY,
  context TEXT,
  description TEXT,
  links_uuid UUID,
  comments_uuid UUID,
  task_entity_uuid UUID
);

CREATE TABLE IF NOT EXISTS task_story_points_entity (
  UUID UUID DEFAULT UUID_generate_v4() NOT NULL PRIMARY KEY,
  complexity VARCHAR(255) NOT NULL,
  uncertainty VARCHAR(255) NOT NULL,
  dependency VARCHAR(255) NOT NULL,
  effort VARCHAR(255) NOT NULL,
  total INT,
  tme_story_points_uuid UUID
);

-- Insert sample data into the new tables

-- Insert data into id_uuid_entity table
INSERT INTO id_uuid_entity (ID, UUID)
VALUES
  ('1', '00000000-0000-0000-0000-000000000001'),
  ('2', '00000000-0000-0000-0000-000000000002');

-- Insert data into task_metadata_entity table
INSERT INTO task_metadata_entity (
  UUID,
  name,
  task_type,
  assignee,
  status,
  priority,
  risk,
  tags,
  dates_uuid,
  story_point_uuid,
  iterations_uuid,
  relations_uuid,
  task_entity_uuid
) VALUES (
  '3e506f19-e4a6-4f6a-bb10-5a5f3bf787e2',
  'Task 1',
  'TASK',
  'Assignee 1',
  'TODO',
  'P1',
  'LOW',
  '{"Tag 1", "Tag 2"}',
  '3e506f19-e4a6-4f6a-bb10-5a5f3bf787e3', -- Replace with an actual CommonDateEntity UUID
  '3e506f19-e4a6-4f6a-bb10-5a5f3bf787e4', -- Replace with an actual TaskStoryPointsEntity UUID
  '3e506f19-e4a6-4f6a-bb10-5a5f3bf787e5', -- Replace with an actual IdUuidStatusEntity UUID
  '3e506f19-e4a6-4f6a-bb10-5a5f3bf787e6', -- Replace with an actual TaskRelationsEntity UUID
  '3e506f19-e4a6-4f6a-bb10-5a5f3bf787e7' -- Replace with an actual TaskEntity UUID
), (
  '4e506f19-e4a6-4f6a-bb10-5a5f3bf787e8',
  'Task 2',
  'TASK',
  'Assignee 2',
  'IN_PROGRESS',
  'P2',
  'MEDIUM',
  '{"Tag 3", "Tag 4"}',
  '4e506f19-e4a6-4f6a-bb10-5a5f3bf787e9', -- Replace with an actual CommonDateEntity UUID
  '4e506f19-e4a6-4f6a-bb10-5a5f3bf787ea', -- Replace with an actual TaskStoryPointsEntity UUID
  '4e506f19-e4a6-4f6a-bb10-5a5f3bf787eb', -- Replace with an actual IdUuidStatusEntity UUID
  '4e506f19-e4a6-4f6a-bb10-5a5f3bf787ec', -- Replace with an actual TaskRelationsEntity UUID
  '4e506f19-e4a6-4f6a-bb10-5a5f3bf787ed' -- Replace with an actual TaskEntity UUID
);

-- Insert data into task_content_entity table
INSERT INTO task_content_entity (
  UUID,
  context,
  description,
  links_uuid,
  comments_uuid,
  task_entity_uuid
) VALUES (
  '5e506f19-e4a6-4f6a-bb10-5a5f3bf787ee',
  'Context for Task 1',
  'Description for Task 1',
  '5e506f19-e4a6-4f6a-bb10-5a5f3bf787ef', -- Replace with an actual NameUrlEntity UUID
  '5e506f19-e4a6-4f6a-bb10-5a5f3bf787f0', -- Replace with an actual IdUuidEntity UUID
  '5e506f19-e4a6-4f6a-bb10-5a5f3bf787f1' -- Replace with an actual TaskEntity UUID
), (
  '6e506f19-e4a6-4f6a-bb10-5a5f3bf787f2',
  'Context for Task 2',
  'Description for Task 2',
  '6e506f19-e4a6-4f6a-bb10-5a5f3bf787f3', -- Replace with an actual NameUrlEntity UUID
  '6e506f19-e4a6-4f6a-bb10-5a5f3bf787f4', -- Replace with an actual IdUuidEntity UUID
  '6e506f19-e4a6-4f6a-bb10-5a5f3bf787f5' -- Replace with an actual TaskEntity UUID
);

-- Insert data into task_story_points_entity table
INSERT INTO task_story_points_entity (
  UUID,
  complexity,
  uncertainty,
  dependency,
  effort,
  total,
  tme_story_points_uuid
) VALUES (
  '7e506f19-e4a6-4f6a-bb10-5a5f3bf787f6',
  'M',
  'SOME',
  'FEW',
  'M',
  10,
  '7e506f19-e4a6-4f6a-bb10-5a5f3bf787f7' -- Replace with an actual TaskMetadataEntity UUID
), (
  '8e506f19-e4a6-4f6a-bb10-5a5f3bf787f8',
  'L',
  'LOW',
  'NONE',
  'L',
  15,
  '8e506f19-e4a6-4f6a-bb10-5a5f3bf787f9' -- Replace with an actual TaskMetadataEntity UUID
);

-- Update the users table to include UUIDs for the new entities
UPDATE users
SET task_metadata_uuid = '3e506f19-e4a6-4f6a-bb10-5a5f3bf787e2' -- Replace with an actual TaskMetadataEntity UUID
WHERE UUID = '00000000-0000-0000-0000-000000000001';

UPDATE users
SET task_metadata_uuid = '4e506f19-e4a6-4f6a-bb10-5a5f3bf787e8' -- Replace with an actual TaskMetadataEntity UUID
WHERE UUID = '00000000-0000-0000-0000-000000000002';
