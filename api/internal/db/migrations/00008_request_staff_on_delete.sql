-- +goose Up

-- approved_by / rejected_by (migration 00001) had no ON DELETE rule, so a
-- staff account that had ever approved or rejected a request could not be
-- deleted. Keep the request and forget who handled it instead, the same as
-- prescriptions.prescribed_by (migration 00007).
ALTER TABLE prescription_requests
    DROP CONSTRAINT prescription_requests_approved_by_fkey,
    ADD CONSTRAINT prescription_requests_approved_by_fkey
        FOREIGN KEY (approved_by) REFERENCES users (id) ON DELETE SET NULL,
    DROP CONSTRAINT prescription_requests_rejected_by_fkey,
    ADD CONSTRAINT prescription_requests_rejected_by_fkey
        FOREIGN KEY (rejected_by) REFERENCES users (id) ON DELETE SET NULL;

-- +goose Down
ALTER TABLE prescription_requests
    DROP CONSTRAINT prescription_requests_approved_by_fkey,
    ADD CONSTRAINT prescription_requests_approved_by_fkey
        FOREIGN KEY (approved_by) REFERENCES users (id),
    DROP CONSTRAINT prescription_requests_rejected_by_fkey,
    ADD CONSTRAINT prescription_requests_rejected_by_fkey
        FOREIGN KEY (rejected_by) REFERENCES users (id);
