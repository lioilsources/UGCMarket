# Ol1n Market (market.ol1n.com). Docker doesn't run on the Mac: the DB is on JODA, see .env.example.
#   make api           # market-api on :8196 against market_dev on JODA, ingest from UGC_URL
#   make web           # SvelteKit dev on :5173 → API :8196
#   make test          # go vet + Go tests (e2e wipes market_test) + svelte-check
#   make sqlc          # regenerate api/internal/db after a SQL change
#   make deploy-db     # only Postgres on JODA (dev/test DB for the Mac), nothing public
#   make deploy-nas    # whole stack on JODA = publishes market.ol1n.com

-include .env
export

NAS     ?= joda
NAS_DIR ?= ~/deploy/RobloxMarket
STACK   := cd $(NAS_DIR)/deploy && docker compose -f market-stack.yaml
RSYNC   := rsync -az --delete --exclude node_modules --exclude .svelte-kit --exclude build --exclude data
EVENT    = ssh $(NAS) "echo '{\"ts\":\"$$(date -u +%FT%TZ)\",\"who\":\"$$(whoami)@$$(hostname -s)\",\"repo\":\"RobloxMarket\",\"action\":\"$(1)\",\"commit\":\"$$(git rev-parse --short HEAD 2>/dev/null || echo none)\"}' >> ~/deploy/events.jsonl"

.PHONY: api web test sqlc sync deploy-db deploy-nas tunnel-config status

api:
	cd api && MARKET_ADDR=:8196 DATABASE_URL='$(DEV_DATABASE_URL)' MARKET_DATA=../data UGC_URL=$(UGC_URL) \
		UGC_POLL_INTERVAL=0 ADMIN_EMAILS=$${ADMIN_EMAILS:-admin@example.com} go run ./cmd/market-api

web:
	cd web && API_URL=http://localhost:8196 npm run dev

test:
	cd api && go vet ./... && MARKET_TEST_DATABASE_URL='$(TEST_DATABASE_URL)' go test -count=1 ./...
	cd web && npm run check

sqlc:
	cd api && sqlc generate

# .env and cloudflared/{config.yml,credentials.json} live only on JODA, rsync won't overwrite them.
sync:
	ssh $(NAS) 'mkdir -p $(NAS_DIR)/deploy/cloudflared'
	$(RSYNC) api web $(NAS):$(NAS_DIR)/
	$(RSYNC) --exclude .env --exclude cloudflared/credentials.json --exclude cloudflared/config.yml \
		deploy/ $(NAS):$(NAS_DIR)/deploy/
	@ssh $(NAS) 'test -s $(NAS_DIR)/deploy/.env' || (echo "chybí $(NAS_DIR)/deploy/.env (viz deploy/.env.example)" && false)
	@ssh $(NAS) 'test -d "$$(sed -n "s/^DATA_DIR=//p" $(NAS_DIR)/deploy/.env)"' \
		|| (echo "DATA_DIR z deploy/.env na JODĚ neexistuje (snap docker by ho tiše založil jinde)" && false)

deploy-db: sync
	ssh $(NAS) '$(STACK) up -d db'
	$(call EVENT,deploy db)

# Finds the UUID of the existing market-nas tunnel (read-only) and prepares config.yml + credentials.json.
tunnel-config:
	ssh $(NAS) 'set -e; cd $(NAS_DIR)/deploy/cloudflared; \
		id=$$(cloudflared tunnel --config /dev/null list -n market-nas -o json | python3 -c "import sys,json;print(json.load(sys.stdin)[0][\"id\"])"); \
		install -m 0444 ~/.cloudflared/$$id.json credentials.json; \
		sed "s/<TUNNEL-ID>/$$id/" config.yml.example > config.yml; echo "tunnel market-nas = $$id"'

deploy-nas: sync tunnel-config
	ssh $(NAS) '$(STACK) up -d --build && $(STACK) ps --format "{{.Name}} {{.Status}}"'
	$(call EVENT,deploy)

status:
	@curl -s -m 8 -o /dev/null -w "market.ol1n.com: %{http_code}\n" https://market.ol1n.com/api/catalog/facets
	@ssh $(NAS) '$(STACK) ps --format "{{.Name}} {{.Status}}"'
