-- Development-only demo data for the clinic side of the app: one staff
-- account and three clients with pets, prescriptions and repeat requests in
-- every status, so the clinic dashboard has something to show.
--
--   docker-compose exec -T db psql -U petsforvets < api/dev/seed_clinic_demo.sql
--
-- Demo logins (all accounts use the same password):
--   vet@example.test       clinic staff (veterinary = true)
--   alex@example.test      client
--   sam@example.test       client
--   jordan@example.test    client
--   password: demo-clinic-2026
--
-- Safe to run more than once: existing demo accounts are left as they are,
-- and pets are only added to demo clients that have none.
-- Remove everything again with:
--   DELETE FROM users WHERE email LIKE '%@example.test';

-- crypt(..., gen_salt('bf')) makes a bcrypt hash, the same format the API
-- creates with golang.org/x/crypto/bcrypt, so these accounts can log in.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO users (email, password_hash, name, surname, veterinary, app_mode)
VALUES
    ('vet@example.test',    crypt('demo-clinic-2026', gen_salt('bf', 10)), 'Priya',  'Shah',   true,  'clinicTeam'),
    ('alex@example.test',   crypt('demo-clinic-2026', gen_salt('bf', 10)), 'Alex',   'Morgan', false, 'petOwner'),
    ('sam@example.test',    crypt('demo-clinic-2026', gen_salt('bf', 10)), 'Sam',    'Okafor', false, 'petOwner'),
    ('jordan@example.test', crypt('demo-clinic-2026', gen_salt('bf', 10)), 'Jordan', 'Lee',    false, 'petOwner')
ON CONFLICT (email) DO NOTHING;

-- Pets. The breed is looked up by name within the species.
INSERT INTO pets (owner_id, name, species_id, breed_id, date_of_birth, gender, neutered, colour, microchip_number)
SELECT u.id, d.pet_name, s.id, b.id, d.dob::date, d.gender, d.neutered, d.colour, d.microchip
FROM (VALUES
    ('alex@example.test',   'Bella',   'Dog',    'Labrador Retriever', '2016-03-12', 'female', true,  'Yellow',          '826098100000001'),
    ('alex@example.test',   'Mochi',   'Cat',    'British Shorthair',  '2021-07-01', 'male',   true,  'Blue',            '826098100000002'),
    ('sam@example.test',    'Rocky',   'Dog',    'Boxer',              '2014-11-20', 'male',   true,  'Fawn',            '826098100000003'),
    ('sam@example.test',    'Clover',  'Rabbit', 'Holland Lop',        '2022-04-05', 'female', false, 'Brown and white', ''),
    ('jordan@example.test', 'Luna',    'Cat',    'Siamese',            '2012-09-15', 'female', true,  'Seal point',      '826098100000005')
) AS d (email, pet_name, species_name, breed_name, dob, gender, neutered, colour, microchip)
JOIN users u        ON u.email = d.email
JOIN species s      ON s.name = d.species_name
LEFT JOIN breeds b  ON b.species_id = s.id AND b.name = d.breed_name
WHERE NOT EXISTS (SELECT 1 FROM pets p WHERE p.owner_id = u.id);

-- Prescriptions written by the demo vet (only for demo pets without any).
INSERT INTO prescriptions (pet_id, medicine_id, dosage, quantity, prescribed_by)
SELECT p.id, m.id, d.dosage, d.quantity, vet.id
FROM (VALUES
    ('Bella',  'Metacam',  '1.5 ml with food once a day',            1),
    ('Bella',  'Bravecto', 'One chewable tablet every 12 weeks',     1),
    ('Mochi',  'Drontal',  'One tablet every 3 months',              2),
    ('Rocky',  'Vetmedin', '5 mg tablet twice a day, before food',   2),
    ('Rocky',  'Apoquel',  '16 mg tablet once a day',                1),
    ('Clover', 'Metacam',  '0.3 ml by mouth once a day for 5 days',  1),
    ('Luna',   'Metacam',  '0.5 ml with food once a day',            1)
) AS d (pet_name, medicine_name, dosage, quantity)
JOIN pets p      ON p.name = d.pet_name
JOIN users owner ON owner.id = p.owner_id AND owner.email LIKE '%@example.test'
JOIN medicines m ON m.name = d.medicine_name
CROSS JOIN (SELECT id FROM users WHERE email = 'vet@example.test') AS vet
WHERE NOT EXISTS (SELECT 1 FROM prescriptions x WHERE x.pet_id = p.id);

-- Repeat requests in every status, at different ages, so the queue has
-- an order to it. Only added for demo prescriptions that have no requests.
INSERT INTO prescription_requests (prescription_id, status, rejection_reason, created_at, approved_by, rejected_by)
SELECT pr.id,
       d.status,
       d.reason,
       now() - d.age::interval,
       CASE WHEN d.status IN ('Approved', 'Ready') THEN vet.id END,
       CASE WHEN d.status = 'Rejected' THEN vet.id END
FROM (VALUES
    ('Bella',  'Metacam',  'Pending',  '',                                                          '3 hours'),
    ('Rocky',  'Vetmedin', 'Pending',  '',                                                          '1 day'),
    ('Luna',   'Metacam',  'Pending',  '',                                                          '2 days'),
    ('Mochi',  'Drontal',  'Approved', '',                                                          '1 day'),
    ('Bella',  'Bravecto', 'Ready',    '',                                                          '4 days'),
    ('Rocky',  'Apoquel',  'Rejected', 'Rocky is due a check-up first. Please book an appointment.', '6 days')
) AS d (pet_name, medicine_name, status, reason, age)
JOIN pets p          ON p.name = d.pet_name
JOIN users owner     ON owner.id = p.owner_id AND owner.email LIKE '%@example.test'
JOIN prescriptions pr ON pr.pet_id = p.id
JOIN medicines m     ON m.id = pr.medicine_id AND m.name = d.medicine_name
CROSS JOIN (SELECT id FROM users WHERE email = 'vet@example.test') AS vet
WHERE NOT EXISTS (SELECT 1 FROM prescription_requests x WHERE x.prescription_id = pr.id);
