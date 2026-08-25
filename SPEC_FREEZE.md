# Foundation specification freeze

- Profile: `compact_10`; target capacity 10 independent runtime boundaries; backend only.
- Domain boundary: multi-province vocational undergraduate admission from annual plan locking through professional-group allocation, application, review, parallel ranking, admission, withdrawal, supplementary admission and auditable notifications. Excluded-topic review passed; this is not an RBAC product, inventory system, order system, dashboard, or appointment system.
- Persistence: SQLite with `modernc.org/sqlite`, five idempotent migrations, nine related tables, foreign keys, unique keys and indexes. Submission transaction reserves both group and plan capacity and writes the application, audit event and notification job atomically.
- State and concurrency: plan draft→locked→published→closed; application submitted→reviewing→admitted/rejected/withdrawn; conditional capacity updates, version checks, deterministic rank ordering, idempotency key and conflict errors.
- Context/errors: every HTTP/service/repository/worker operation accepts context; cancellation and wrapped sentinel errors are preserved. HTTP maps errors to stable status responses with request IDs and panic recovery.
- Identity: password hashing, login, expiring server sessions, logout revocation, inactive/expired rejection, and admin/officer/reviewer/viewer role differences.
- Worker/audit/recovery: cancellable retry worker, backoff, permanent failure record, decision/audit tables, idempotent restart migration and file-database recovery test.
- HTTP/ops: `/healthz`, `/readyz`, auth endpoints, protected plan/application endpoints, structured logs, graceful shutdown, environment configuration, Docker default entrypoint `/app/admission`.
- Testing: domain, service, real SQLite migration/rollback/restart, HTTP contract, idempotency, capacity, state transitions, cancellation, worker retry and pagination tests; race, vet and build gates executed.
- Scale: 43 production Go files, 12 packages, 2000 physical production lines, 1516 test lines; no generated/vendor padding or frontend.
- Future boundaries: plan locking, rule validation, capacity reservation, idempotent submit, review transition, ranking allocation, audit persistence, session lifecycle, worker retry and restart recovery can each support one independent candidate without pre-seeding a defect.
