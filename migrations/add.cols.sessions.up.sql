ALTER TABLE sessions 
ADD COLUMN is_recurring BOOLEAN NOT NULL;

ALTER TABLE	sessions 
ADD COLUMN plan_id varchar(255);
