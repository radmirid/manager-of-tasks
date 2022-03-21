set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASSWORD psql -h "$host" -U "postgres" -c '\q'; do
  >&2 echo "PostgreSQL is sleeping"
  sleep 1
done

>&2 echo "PostgreSQL is running - command execution"
exec $cmd