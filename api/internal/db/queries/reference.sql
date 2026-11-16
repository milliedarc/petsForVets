-- Reference data for the pet form's dropdowns.

-- name: ListSpecies :many
SELECT id, name, pet_type_id
FROM species
ORDER BY name;

-- Generic breeds ("Mixed breed") first, then alphabetical.
-- pet_type_id is joined in because the form filters breeds by pet type.
-- name: ListBreeds :many
SELECT b.id, b.name, b.species_id, s.pet_type_id, b.generic
FROM breeds b
JOIN species s ON s.id = b.species_id
ORDER BY b.generic DESC, b.name;
