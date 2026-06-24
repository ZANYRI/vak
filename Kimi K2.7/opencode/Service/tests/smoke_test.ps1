$base = 'http://localhost:8080/api/v1'
$headers = @{ 'Content-Type' = 'application/json' }

function Invoke-Api($method, $path, $body=$null, $extraHeaders=@{}) {
    $url = "$base$path"
    $h = $headers.Clone()
    foreach ($k in $extraHeaders.Keys) { $h[$k] = $extraHeaders[$k] }
    if ($body) {
        return Invoke-RestMethod -Method $method -Uri $url -Headers $h -Body ($body | ConvertTo-Json -Depth 5)
    }
    return Invoke-RestMethod -Method $method -Uri $url -Headers $h
}

Write-Host '1. Register admin'
$user = Invoke-Api POST '/auth/register' @{
    email = 'admin@example.com'
    password = 'password123'
    name = 'Admin User'
    role = 'admin'
}
Write-Host "  user id: $($user.user.id), role: $($user.user.role)"

Write-Host '2. Login'
$login = Invoke-Api POST '/auth/login' @{ email = 'admin@example.com'; password = 'password123' }
$token = $login.access_token
Write-Host "  token received: $($token.Substring(0,20))..."

Write-Host '3. Current user'
$me = Invoke-Api GET '/auth/me' -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  me: $($me.user.email), role: $($me.user.role)"

Write-Host '4. Create per-seat plan'
$plan = Invoke-Api POST '/plans' @{
    name = 'Team Plan'
    currency = 'USD'
    billing_interval = 'monthly'
    base_price_cents = 1200
    trial_days = 0
    is_active = $true
    prices = @(@{
        model = 'per_seat'
        seat_price_cents = 700
        included_seats = 3
    })
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  plan id: $($plan.id), base: $($plan.base_price_cents)"

Write-Host '5. Create customer'
$customer = Invoke-Api POST '/customers' @{
    email = 'customer@example.com'
    name = 'ACME Inc'
    country = 'US'
    region = 'CA'
    currency = 'USD'
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  customer id: $($customer.id)"

Write-Host '6. Create subscription (5 seats -> 2 extra)'
$subscription = Invoke-Api POST '/subscriptions' @{
    customer_id = $customer.id
    plan_id = $plan.id
    quantity = 1
    seats = 5
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  subscription id: $($subscription.id), status: $($subscription.status)"

Write-Host '7. Report usage'
$usage = Invoke-Api POST '/usage' @{
    customer_id = $customer.id
    subscription_id = $subscription.id
    metric = 'api_requests'
    quantity = 1500
    idempotency_key = 'usage-demo-1'
    recorded_at = (Get-Date -Format o)
} -extraHeaders @{ 'Authorization' = "Bearer $token"; 'Idempotency-Key' = 'usage-demo-1' }
Write-Host "  usage event id: $($usage.id)"

Write-Host '8. Create coupon (30% off)'
$coupon = Invoke-Api POST '/coupons' @{
    code = 'SUMMER30'
    type = 'percentage'
    percent_off = 30
    is_active = $true
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  coupon id: $($coupon.id)"

Write-Host '9. Apply coupon to subscription'
$applied = Invoke-Api POST "/subscriptions/$($subscription.id)/apply-coupon" @{
    coupon_id = $coupon.id
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  result: $($applied.message)"

Write-Host '10. Create tax rule (10%)'
$tax = Invoke-Api POST '/tax-rules' @{
    country = 'US'
    region = 'CA'
    tax_name = 'CA Sales Tax'
    tax_rate_basis_points = 1000
    is_active = $true
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  tax rule id: $($tax.id)"

Write-Host '11. Generate invoice'
$invoice = Invoke-Api POST '/invoices/generate' @{
    subscription_id = $subscription.id
    customer_id = $customer.id
} -extraHeaders @{ 'Authorization' = "Bearer $token"; 'Idempotency-Key' = 'inv-demo-1' }
Write-Host "  invoice id: $($invoice.id), subtotal: $($invoice.subtotal_cents), discount: $($invoice.discount_cents), tax: $($invoice.tax_cents), total: $($invoice.total_cents), lines: $($invoice.lines.Count)"
foreach ($ln in $invoice.lines) {
    Write-Host "    - $($ln.type): $($ln.amount_cents) ($($ln.description))"
}

Write-Host '12. Finalize invoice'
$finalized = Invoke-Api POST "/invoices/$($invoice.id)/finalize" $null -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  status: $($finalized.status)"

Write-Host '13. Simulate payment (card ending 0000 -> success)'
$payment = Invoke-Api POST '/payments/simulate' @{
    invoice_id = $invoice.id
    card_number = '4242424242420000'
} -extraHeaders @{ 'Authorization' = "Bearer $token"; 'Idempotency-Key' = 'pay-demo-1' }
Write-Host "  payment id: $($payment.id), status: $($payment.status), amount: $($payment.amount_cents)"

Write-Host '14. Get invoice after payment'
$paidInvoice = Invoke-Api GET "/invoices/$($invoice.id)" -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  invoice status: $($paidInvoice.status), paid_at: $($paidInvoice.paid_at)"

Write-Host '15. Create cheaper plan and change subscription (proration)'
$cheapPlan = Invoke-Api POST '/plans' @{
    name = 'Basic Plan'
    currency = 'USD'
    billing_interval = 'monthly'
    base_price_cents = 500
    trial_days = 0
    is_active = $true
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
$change = Invoke-Api POST "/subscriptions/$($subscription.id)/change-plan" @{
    plan_id = $cheapPlan.id
} -extraHeaders @{ 'Authorization' = "Bearer $token" }
Write-Host "  proration amount cents: $($change.proration.amount_cents)"

Write-Host '16. Generate second invoice after plan change'
$invoice2 = Invoke-Api POST '/invoices/generate' @{
    subscription_id = $subscription.id
    customer_id = $customer.id
} -extraHeaders @{ 'Authorization' = "Bearer $token"; 'Idempotency-Key' = 'inv-demo-2' }
Write-Host "  invoice2 total: $($invoice2.total_cents)"

Write-Host '17. List jobs in queue DB'
$jobsJson = docker exec billing-service-postgres-1 psql -U billing -d billing -t -A -c "SELECT queue, status, COUNT(*) FROM jobs GROUP BY queue, status;"
Write-Host "  jobs: $jobsJson"

Write-Host '18. Idempotency check (same usage key returns existing event id)'
$usage2 = Invoke-Api POST '/usage' @{
    customer_id = $customer.id
    subscription_id = $subscription.id
    metric = 'api_requests'
    quantity = 9999
    idempotency_key = 'usage-demo-1'
} -extraHeaders @{ 'Authorization' = "Bearer $token"; 'Idempotency-Key' = 'usage-demo-1' }
if ($usage2.id -eq $usage.id) {
    Write-Host "  OK: idempotency preserved, same id $($usage2.id)"
} else {
    Write-Host "  FAIL: idempotency broken, $($usage.id) vs $($usage2.id)"
}

Write-Host "`nAll checks completed."
