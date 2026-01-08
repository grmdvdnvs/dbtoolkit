# dbtoolkit

CLI para migraciones y administración de sesiones.

Comandos:
- `dbtoolkit migrate --source <dsn> --target <dsn> --tables a,b --strategy incremental --dry-run`
- `dbtoolkit sessions list --db oracle --older-than 2h`
- `dbtoolkit sessions kill --db oracle --user app_user --preview=false`

Notas:
- Añadir drivers (`godror`, `lib/pq`) según las bases usadas.
- Completar adaptadores DB y mapeos de esquemas para producción.