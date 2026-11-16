-- Ownership is checked by joining to pets and filtering on owner_id, the
-- same rule as for pets themselves.

-- name: ListPrescriptionsForPet :many
SELECT pr.id,
       pr.dosage,
       pr.quantity,
       m.id          AS medicine_id,
       m.name        AS medicine_name,
       m.description AS medicine_description
FROM prescriptions pr
JOIN medicines m ON m.id = pr.medicine_id
JOIN pets p      ON p.id = pr.pet_id
WHERE pr.pet_id = @pet_id
  AND p.owner_id = @owner_id
ORDER BY m.name;

-- name: ListPrescriptionRequestsForPet :many
SELECT *
FROM prescription_request_details
WHERE pet_id = @pet_id
  AND owner_id = @owner_id
ORDER BY created_at DESC;

-- name: GetPrescriptionRequestForOwner :one
SELECT *
FROM prescription_request_details
WHERE id = @id
  AND owner_id = @owner_id;

-- INSERT ... SELECT: the row is only inserted if the SELECT finds the
-- prescription *and* it belongs to one of the user's pets. No match means
-- nothing is inserted and nothing is returned (the API answers 404).
-- name: CreatePrescriptionRequest :one
INSERT INTO prescription_requests (prescription_id)
SELECT pr.id
FROM prescriptions pr
JOIN pets p ON p.id = pr.pet_id
WHERE pr.id = @prescription_id
  AND p.owner_id = @owner_id
RETURNING id;

-- Only pending requests can be cancelled. DELETE ... USING is Postgres's way
-- of joining other tables into a DELETE.
-- name: CancelPrescriptionRequest :execrows
DELETE FROM prescription_requests r
USING prescriptions pr, pets p
WHERE r.id = @id
  AND pr.id = r.prescription_id
  AND p.id = pr.pet_id
  AND p.owner_id = @owner_id
  AND r.status = 'Pending';
