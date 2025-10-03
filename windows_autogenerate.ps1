docker run --rm -v ${PWD}:/src -w /src sqlc/sqlc generate
Start-Sleep -Seconds 5
mockgen --package mockdatabase --destination autogenerate/database/mock/store.go github.com/WhiteMaks/go_academy_portal_service/autogenerate/database Store