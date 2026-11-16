-- +goose Up

-- 1. Normalise: a request's pet is its prescription's pet. Storing pet_id here
--    as well meant the two could disagree, so derive it with a JOIN instead.
--    (Dropping the column drops its index too.)
ALTER TABLE prescription_requests DROP COLUMN pet_id;

CREATE INDEX prescription_requests_prescription_id_idx
    ON prescription_requests (prescription_id);

-- 2. At most one *open* request per prescription. A partial index only
--    contains rows matching its WHERE clause, so finished requests (Ready,
--    Rejected) don't count and are kept as history.
CREATE UNIQUE INDEX prescription_requests_one_open_idx
    ON prescription_requests (prescription_id)
    WHERE status IN ('Pending', 'Approved');

-- 3. Like pet_details: the joins every request query needs, written once.
CREATE VIEW prescription_request_details AS
SELECT r.id,
       r.status,
       r.comments,
       r.rejection_reason,
       r.created_at,
       r.prescription_id,
       pr.dosage,
       m.name     AS medicine_name,
       p.id       AS pet_id,
       p.owner_id AS owner_id
FROM prescription_requests r
JOIN prescriptions pr ON pr.id = r.prescription_id
JOIN medicines m      ON m.id = pr.medicine_id
JOIN pets p           ON p.id = pr.pet_id;

-- +goose Down
DROP VIEW prescription_request_details;
DROP INDEX prescription_requests_one_open_idx;
DROP INDEX prescription_requests_prescription_id_idx;

ALTER TABLE prescription_requests
    ADD COLUMN pet_id uuid REFERENCES pets (id) ON DELETE CASCADE;
UPDATE prescription_requests r
SET pet_id = pr.pet_id
FROM prescriptions pr
WHERE pr.id = r.prescription_id;
ALTER TABLE prescription_requests ALTER COLUMN pet_id SET NOT NULL;
CREATE INDEX prescription_requests_pet_id_idx ON prescription_requests (pet_id);
