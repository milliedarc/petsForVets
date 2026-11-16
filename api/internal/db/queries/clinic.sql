-- Queries for clinic staff. Unlike the pet owner queries, these have no
-- owner_id filter: there is one clinic, and its staff look after every pet.
-- The API only lets staff (users.veterinary = true) reach them; see
-- requireStaff in internal/httpapi/clinic.go.
--
-- sqlc.narg('x') is an optional ("nullable") parameter: the Go function takes
-- a pointer, and nil means "don't filter on this". `@x::text IS NULL OR ...`
-- is the usual SQL trick for an optional WHERE condition.

-- name: GetClinicSummary :one
SELECT (SELECT count(*) FROM prescription_requests WHERE status = 'Pending')  AS pending_requests,
       (SELECT count(*) FROM prescription_requests WHERE status = 'Approved') AS approved_requests,
       (SELECT count(*) FROM prescription_requests WHERE status = 'Ready')    AS ready_requests,
       (SELECT count(*) FROM pets)                                            AS patients,
       (SELECT count(DISTINCT owner_id) FROM pets)                            AS clients;

-- Oldest first, so the queue is worked through in the order requests came in.
-- name: ListClinicRequests :many
SELECT *
FROM clinic_request_details
WHERE (sqlc.narg('status')::text IS NULL OR status = sqlc.narg('status'))
  AND (sqlc.narg('pet_id')::uuid IS NULL OR pet_id = sqlc.narg('pet_id'))
ORDER BY created_at
LIMIT 500;

-- name: GetClinicRequest :one
SELECT *
FROM clinic_request_details
WHERE id = @id;

-- Moves a request to a new status, but only from one of the statuses allowed
-- before it (e.g. only a Pending request can be Approved). If the request has
-- moved on in the meantime, nothing matches and no row comes back, so two
-- staff members can't both act on the same request.
-- name: UpdateClinicRequestStatus :one
UPDATE prescription_requests
SET status           = @new_status::text,
    rejection_reason = CASE WHEN @new_status::text = 'Rejected' THEN @rejection_reason::text ELSE rejection_reason END,
    comments         = COALESCE(sqlc.narg('comments')::text, comments),
    approved_by      = CASE
                           WHEN @new_status::text IN ('Approved', 'Ready') AND approved_by IS NULL THEN @staff_id::uuid
                           ELSE approved_by
                       END,
    rejected_by      = CASE WHEN @new_status::text = 'Rejected' THEN @staff_id::uuid ELSE rejected_by END
WHERE id = @id
  AND status = ANY (@from_statuses::text[])
RETURNING id;

-- Patients, with their owner and how many requests are waiting on the clinic.
-- sqlc.embed(pd) puts the pet_details columns into a nested PetDetail struct,
-- so the existing newPetResponse() can turn it into JSON.
-- name: ListClinicPets :many
SELECT sqlc.embed(pd),
       u.name    AS owner_name,
       u.surname AS owner_surname,
       u.email   AS owner_email,
       (SELECT count(*)
        FROM prescription_requests r
        JOIN prescriptions pr ON pr.id = r.prescription_id
        WHERE pr.pet_id = pd.id
          AND r.status IN ('Pending', 'Approved')) AS open_requests
FROM pet_details pd
JOIN users u ON u.id = pd.owner_id
WHERE sqlc.narg('search')::text IS NULL
   OR pd.name ILIKE sqlc.narg('search')
   OR u.name ILIKE sqlc.narg('search')
   OR u.surname ILIKE sqlc.narg('search')
   OR (u.name || ' ' || u.surname) ILIKE sqlc.narg('search')
   OR u.email ILIKE sqlc.narg('search')
   OR pd.microchip_number ILIKE sqlc.narg('search')
ORDER BY pd.name, u.surname
LIMIT 100;

-- name: GetClinicPet :one
SELECT sqlc.embed(pd),
       u.name    AS owner_name,
       u.surname AS owner_surname,
       u.email   AS owner_email,
       (SELECT count(*)
        FROM prescription_requests r
        JOIN prescriptions pr ON pr.id = r.prescription_id
        WHERE pr.pet_id = pd.id
          AND r.status IN ('Pending', 'Approved')) AS open_requests
FROM pet_details pd
JOIN users u ON u.id = pd.owner_id
WHERE pd.id = @id;

-- name: ListClinicPetPrescriptions :many
SELECT pr.id,
       pr.dosage,
       pr.quantity,
       pr.created_at,
       m.id          AS medicine_id,
       m.name        AS medicine_name,
       m.description AS medicine_description,
       staff.name    AS prescribed_by_name,
       staff.surname AS prescribed_by_surname
FROM prescriptions pr
JOIN medicines m      ON m.id = pr.medicine_id
LEFT JOIN users staff ON staff.id = pr.prescribed_by
WHERE pr.pet_id = @pet_id
ORDER BY pr.created_at DESC;

-- name: CreatePrescription :one
INSERT INTO prescriptions (pet_id, medicine_id, dosage, quantity, prescribed_by)
VALUES (@pet_id, @medicine_id, @dosage, @quantity, @prescribed_by)
RETURNING id;

-- name: ListMedicines :many
SELECT id, name, description
FROM medicines
ORDER BY name;
