CREATE TYPE root.user_role AS ENUM ('admin');

ALTER TYPE root.user_role ADD VALUE 'coach';
ALTER TYPE root.user_role ADD VALUE 'athlete';
