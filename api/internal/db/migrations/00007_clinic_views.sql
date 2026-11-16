-- +goose Up

-- 1. Who wrote each prescription. Nullable: older rows (and the dev seed
--    script) have no author. ON DELETE SET NULL keeps the prescription if the
--    staff account is ever deleted.
ALTER TABLE prescriptions
    ADD COLUMN prescribed_by uuid REFERENCES users (id) ON DELETE SET NULL;

-- 2. The clinic's view of a prescription request: everything the request
--    queue shows (medicine, pet, owner) in one row, so the queries that list
--    requests and the one that returns a request after updating it share the
--    same joins. Like prescription_request_details, but for staff, so it
--    includes the owner's contact details.
CREATE VIEW clinic_request_details AS
SELECT r.id,
       r.status,
       r.comments,
       r.rejection_reason,
       r.created_at,
       pr.id       AS prescription_id,
       pr.dosage,
       pr.quantity,
       m.name      AS medicine_name,
       p.id        AS pet_id,
       p.name      AS pet_name,
       s.name      AS species_name,
       u.id        AS owner_id,
       u.name      AS owner_name,
       u.surname   AS owner_surname,
       u.email     AS owner_email
FROM prescription_requests r
JOIN prescriptions pr ON pr.id = r.prescription_id
JOIN medicines m      ON m.id = pr.medicine_id
JOIN pets p           ON p.id = pr.pet_id
JOIN species s        ON s.id = p.species_id
JOIN users u          ON u.id = p.owner_id;

-- +goose Down
DROP VIEW clinic_request_details;
ALTER TABLE prescriptions DROP COLUMN prescribed_by;
