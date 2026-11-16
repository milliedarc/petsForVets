-- Development-only sample data: gives every pet of the dev account two
-- repeat prescriptions, so the prescriptions page has something to show
-- until the clinic side (where staff write prescriptions) exists.
-- Safe to run more than once: pets that already have prescriptions are skipped.
--
--   docker-compose exec -T db psql -U petsforvets < api/dev/seed_dev_data.sql

INSERT INTO prescriptions (pet_id, medicine_id, dosage, quantity)
SELECT p.id, m.id, d.dosage, d.quantity
FROM pets p
JOIN users u ON u.id = p.owner_id
CROSS JOIN (VALUES
    ('Bravecto', 'One chewable tablet every 12 weeks', 1),
    ('Drontal',  'One tablet every 3 months',          2)
) AS d (medicine_name, dosage, quantity)
JOIN medicines m ON m.name = d.medicine_name
WHERE u.email = 'dev@petsforvets.test'
  AND NOT EXISTS (SELECT 1 FROM prescriptions x WHERE x.pet_id = p.id);
