-- Reads go through the pet_details view (migration 00003), which joins in
-- species and breed. That replaces PocketBase's `expand: 'species,breed'`.
--
-- Every query has `owner_id = ...` in its WHERE clause. That is the access
-- rule: you can only see or change your own pets. Another user's pet id
-- simply matches no rows, and the API answers 404.
--
-- @name is sqlc's named-parameter syntax. It becomes $1, $2 ... in the real
-- SQL and gives the generated Go struct readable field names.

-- name: ListPetsByOwner :many
SELECT *
FROM pet_details
WHERE owner_id = @owner_id
ORDER BY name;

-- name: GetPetForOwner :one
SELECT *
FROM pet_details
WHERE id = @id
  AND owner_id = @owner_id;

-- name: CreatePet :one
INSERT INTO pets (owner_id, name, species_id, breed_id, date_of_birth, gender,
                  neutered, colour, imported, import_country_code, microchip_number)
VALUES (@owner_id, @name, @species_id, @breed_id, @date_of_birth, @gender,
        @neutered, @colour, @imported, @import_country_code, @microchip_number)
RETURNING id;

-- :execrows returns how many rows changed, so 0 means "not found / not yours".
-- name: UpdatePet :execrows
UPDATE pets
SET name                = @name,
    species_id          = @species_id,
    breed_id            = @breed_id,
    date_of_birth       = @date_of_birth,
    gender              = @gender,
    neutered            = @neutered,
    colour              = @colour,
    imported            = @imported,
    import_country_code = @import_country_code,
    microchip_number    = @microchip_number,
    updated_at          = now()
WHERE id = @id
  AND owner_id = @owner_id;

-- name: DeletePet :execrows
DELETE FROM pets
WHERE id = @id
  AND owner_id = @owner_id;
