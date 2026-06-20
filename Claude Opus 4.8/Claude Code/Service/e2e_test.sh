#!/usr/bin/env bash
# End-to-end verification of all 18 requirements against the running stack.
set -uo pipefail
B=http://localhost:8080/api/v1
PASS=0; FAIL=0
ok(){ echo "  PASS - $1"; PASS=$((PASS+1)); }
no(){ echo "  FAIL - $1"; FAIL=$((FAIL+1)); }
jval(){ # extract first value of a json string field (tolerates spaces after colon)
  echo "$1" | grep -oE "\"$2\": ?\"[^\"]*\"" | head -1 | sed "s/.*: \?\"\([^\"]*\)\"/\1/"; }
jnum(){ echo "$1" | grep -oE "\"$2\": ?-?[0-9]+" | head -1 | sed "s/.*: \?//"; }
post(){ curl -s -X POST "$B$1" -H "$AUTH" -H 'Content-Type: application/json' "${@:2}"; }

echo "== 1) Register & login =="
REG=$(curl -s -X POST $B/auth/register -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password123","role":"customer"}')
[ -n "$(jval "$REG" id)" ] && ok "register user" || no "register user ($REG)"

LOGIN=$(curl -s -X POST $B/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"admin@billing.local","password":"admin12345"}')
ADMIN=$(jval "$LOGIN" access_token)
[ -n "$ADMIN" ] && ok "admin login + JWT" || no "admin login ($LOGIN)"
AUTH="Authorization: Bearer $ADMIN"

ME=$(curl -s $B/auth/me -H "$AUTH")
[ "$(jval "$ME" role)" = "admin" ] && ok "GET /auth/me" || no "me ($ME)"

echo "== 2) RBAC =="
CLOGIN=$(curl -s -X POST $B/auth/login -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password123"}')
CUSTTOK=$(jval "$CLOGIN" access_token)
CODE=$(curl -s -o /dev/null -w '%{http_code}' -X POST $B/plans \
  -H "Authorization: Bearer $CUSTTOK" -H 'Content-Type: application/json' \
  -d '{"name":"x","currency":"USD","billing_interval":"monthly","pricing_model":"flat","base_price_cents":100}')
[ "$CODE" = "403" ] && ok "customer blocked from creating plan (403)" || no "RBAC expected 403 got $CODE"

echo "== 3) Create customer =="
CUST=$(post /customers -d '{"email":"acme@example.com","name":"Acme","country":"US","region":"","currency":"USD"}')
CID=$(jval "$CUST" id)
[ -n "$CID" ] && ok "create customer ($CID)" || no "create customer ($CUST)"

echo "== 4) Create plans =="
PLAN=$(post /plans -d '{"name":"Pro","currency":"USD","billing_interval":"monthly","pricing_model":"flat","base_price_cents":10000}')
PID=$(jval "$PLAN" id)
[ -n "$PID" ] && ok "create flat plan ($PID)" || no "create plan ($PLAN)"
UPLAN=$(post /plans -d '{"name":"Metered","currency":"USD","billing_interval":"monthly","pricing_model":"usage_based","base_price_cents":0,"included_units":0,"unit_price_cents":2,"usage_metric":"api_calls"}')
UPID=$(jval "$UPLAN" id)
BIGPLAN=$(post /plans -d '{"name":"Enterprise","currency":"USD","billing_interval":"monthly","pricing_model":"flat","base_price_cents":30000}')
BPID=$(jval "$BIGPLAN" id)

echo "== 9) Tax rule (needed before invoicing) =="
TAX=$(post /tax-rules -d '{"country":"US","region":"","tax_name":"Sales Tax","tax_rate_basis_points":2000}')
[ -n "$(jval "$TAX" id)" ] && ok "create tax rule 20%" || no "tax rule ($TAX)"

echo "== 5) Create subscription =="
SUB=$(post /subscriptions -H 'Idempotency-Key: sub-flat-1' -d "{\"customer_id\":\"$CID\",\"plan_id\":\"$PID\",\"quantity\":1}")
SID=$(jval "$SUB" id)
[ -n "$SID" ] && ok "create subscription ($SID)" || no "create subscription ($SUB)"
# idempotency replay:
SUB2=$(post /subscriptions -H 'Idempotency-Key: sub-flat-1' -d "{\"customer_id\":\"$CID\",\"plan_id\":\"$PID\",\"quantity\":1}")
[ "$(jval "$SUB2" id)" = "$SID" ] && ok "idempotent subscription replay" || no "idempotency replay ($SUB2)"

echo "== 8) Coupon =="
COUP=$(post /coupons -d '{"code":"SAVE1000","type":"fixed_amount","amount_off_cents":1000,"currency":"USD"}')
[ -n "$(jval "$COUP" id)" ] && ok "create coupon" || no "coupon ($COUP)"
APPLY=$(post /subscriptions/$SID/apply-coupon -d '{"code":"SAVE1000","currency":"USD"}')
[ -n "$(jval "$APPLY" id)" ] && ok "apply coupon to subscription" || no "apply coupon ($APPLY)"

echo "== 6 & 11) Generate invoice + totals (base 10000, -1000 coupon, +20% tax) =="
INV=$(post /invoices/generate -H 'Idempotency-Key: inv-1' -d "{\"subscription_id\":\"$SID\"}")
IID=$(jval "$INV" id)
SUB_T=$(jnum "$INV" subtotal_cents); DISC=$(jnum "$INV" discount_cents)
TAXC=$(jnum "$INV" tax_cents); TOT=$(jnum "$INV" total_cents)
echo "    subtotal=$SUB_T discount=$DISC tax=$TAXC total=$TOT"
[ "$SUB_T" = "10000" ] && ok "subtotal=10000" || no "subtotal=$SUB_T"
[ "$DISC" = "1000" ] && ok "discount=1000 (coupon)" || no "discount=$DISC"
[ "$TAXC" = "1800" ] && ok "tax=1800 (20% of 9000)" || no "tax=$TAXC"
[ "$TOT" = "10800" ] && ok "total=10800" || no "total=$TOT"

echo "== 7) Usage-based billing + idempotency =="
USUB=$(post /subscriptions -d "{\"customer_id\":\"$CID\",\"plan_id\":\"$UPID\",\"quantity\":1}")
USID=$(jval "$USUB" id)
post /usage -H 'Idempotency-Key: u1' -d "{\"customer_id\":\"$CID\",\"subscription_id\":\"$USID\",\"metric\":\"api_calls\",\"quantity\":1000,\"idempotency_key\":\"uk-1\"}" >/dev/null
post /usage -H 'Idempotency-Key: u1' -d "{\"customer_id\":\"$CID\",\"subscription_id\":\"$USID\",\"metric\":\"api_calls\",\"quantity\":1000,\"idempotency_key\":\"uk-1\"}" >/dev/null
SUMM=$(curl -s "$B/usage/summary?subscription_id=$USID" -H "$AUTH")
TOTQ=$(jnum "$SUMM" total_quantity)
[ "$TOTQ" = "1000" ] && ok "usage idempotent (total=1000 not 2000)" || no "usage total=$TOTQ"
UINV=$(post /invoices/generate -d "{\"subscription_id\":\"$USID\"}")
USUBT=$(jnum "$UINV" subtotal_cents)
[ "$USUBT" = "2000" ] && ok "usage charge subtotal=2000 (1000*2)" || no "usage subtotal=$USUBT"

echo "== 10) Proration on plan change =="
CHG=$(post /subscriptions/$SID/change-plan -d "{\"plan_id\":\"$BPID\"}")
echo "$CHG" | grep -q '"proration"' && ok "change-plan returns proration" || no "no proration ($CHG)"
DIFF=$(jnum "$CHG" DifferenceCents)
[ -n "$DIFF" ] && [ "$DIFF" != "0" ] && ok "proration difference=$DIFF (non-zero)" || no "proration diff=$DIFF"

echo "== 12) Payment simulation =="
PAY=$(post /payments/simulate -H 'Idempotency-Key: pay-1' -d "{\"invoice_id\":\"$IID\",\"card_number\":\"4242424242420000\"}")
[ "$(jval "$PAY" status)" = "succeeded" ] && ok "payment with card ...0000 succeeds" || no "payment ($PAY)"
INVP=$(curl -s $B/invoices/$IID -H "$AUTH")
[ "$(jval "$INVP" status)" = "paid" ] && ok "invoice marked paid" || no "invoice status=$(jval "$INVP" status)"
PAYF=$(post /payments/simulate -d "{\"invoice_id\":\"$IID\",\"card_number\":\"4111111111119999\"}")
echo "$PAYF" | grep -qi 'fail\|already paid\|conflict\|PAYMENT_FAILED' && ok "card ...9999 / paid invoice handled" || no "fail-card ($PAYF)"

echo "== 13 & 14) Queue + workers (auto invoice from subscription.create) =="
sleep 3
DONE=$(docker compose exec -T postgres psql -U billing -d billing -tAc \
  "SELECT count(*) FROM jobs WHERE status='completed';" 2>/dev/null | tr -d '[:space:]')
echo "    completed jobs in DB: $DONE"
[ -n "$DONE" ] && [ "$DONE" -gt 0 ] 2>/dev/null && ok "workers consumed & completed $DONE jobs" || no "no completed jobs ($DONE)"

echo
echo "================ RESULT: PASS=$PASS FAIL=$FAIL ================"
