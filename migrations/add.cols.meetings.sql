ALTER TABLE meetings
ADD COLUMN is_recurring BOOLEAN NOT NULL;

ALTER TABLE meetings	
ADD COLUMN plan_id varchar(255);

UPDATE TABLE meetings
SET plan_id = "" 
WHERE id = 1;
