# Postgres 18 migration (production bind mount)

Use this when upgrading a **deployed** instance that uses the generated [`docker-compose.yml`](../../templates/project/docker-compose.yml) bind mount at `postgres/volume` (PG 18+ expects data under `/var/lib/postgresql`, not the old `/var/lib/postgresql/data`-only layout). See [refactor backlog — Database / Postgres](refactor-backlog.md) for template context.

**Assumptions:** you can tolerate downtime; you replace the Postgres image to 18 and reset the data directory; placeholders below are filled for your host and container names.

1. **Dump everything** from the **old** Postgres container (still on the pre-18 image while it runs):

   ```bash
   docker exec -t <pg_container> pg_dumpall -U postgres > dump_pg.sql
   ```

2. **Edit the dump:** remove the line containing `CREATE ROLE postgres` (the new cluster already defines the `postgres` superuser; replaying that statement fails).

3. **Archive the old data directory** (path matches generated server layout: `$serverPath/postgres/volume`):

   ```bash
   mv /home/deploy/<project_name>/postgres/volume /home/deploy/<project_name>/postgres/volume.pg-before-18
   mkdir -p /home/deploy/<project_name>/postgres/volume
   ```

4. **Start Postgres 18** with the updated compose image (empty `volume` dir initializes a fresh PG18 data area):

   ```bash
   docker compose up -d postgres
   ```

5. **Restore** into the **new** container (`<pg18_container>` is typically the new `postgres` service container name):

   ```bash
   cat dump_pg.sql | docker exec -i <pg18_container> psql -U postgres -v ON_ERROR_STOP=1
   ```

For a logical dump/restore workflow without replacing the bind mount (e.g. local dev patterns), see [`templates/project/restoreDump.sh`](../../templates/project/restoreDump.sh).
