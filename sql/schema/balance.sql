-- name: GetBalanceByUser :many
WITH RECURSIVE ordered_in AS (
    SELECT 
        t.*,
        ROW_NUMBER() OVER (PARTITION BY t.symbol ORDER BY t.date) AS rn
    FROM "transaction" t
    WHERE t.side = 'buy' AND t.user_id = $1
), running_totals as (
    SELECT symbol,amount,price,amount::numeric as total, 0::numeric as prev_total, rn 
    FROM ordered_in
    WHERE rn = 1
    UNION ALL
    SELECT rt.symbol, oi.amount, oi.price, rt.total + oi.amount, rt.total, oi.rn
    FROM
        running_totals rt
            INNER JOIN
        ordered_in oi
            ON
                rt.symbol = oi.symbol AND
                rt.rn = oi.rn - 1
), total_out AS (
    SELECT 
        symbol,
        SUM(amount) AS amount
    FROM "transaction"
    WHERE side='sell' AND user_id = $1
    GROUP BY symbol
)
SELECT
    rt.symbol,
    CAST(SUM(CASE WHEN prev_total > COALESCE(out.amount,0) THEN rt.amount ELSE rt.total - COALESCE(out.amount,0) END * price) AS DOUBLE PRECISION) AS cost_basis,
    CAST(SUM(CASE WHEN prev_total > COALESCE(out.amount,0) THEN rt.amount ELSE rt.total - COALESCE(out.amount,0) END) AS DOUBLE PRECISION) AS amount
FROM
    running_totals rt
        LEFT JOIN
    total_out out
        ON
            rt.symbol = out.symbol
WHERE
    rt.total > COALESCE(out.amount, 0) 
GROUP BY rt.symbol;
