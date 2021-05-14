-- name: GetBalance :many
SELECT 
    symbol,
    CAST(SUM(CASE 
            WHEN t.side = 'buy' THEN t.amount
            ELSE t.amount * -1
        END
    ) AS DOUBLE PRECISION) AS total
FROM "transaction" as t
WHERE t.user_id = $1
GROUP BY symbol;

-- name: GetCostBasis :one
WITH A AS (
    SELECT 
        row_number() OVER (ORDER BY date) n,
        side,
        amount,
        price,
        SUM(
            CASE 
                WHEN side = 'buy' THEN amount
                ELSE amount * -1
            END
        ) OVER (ORDER BY date, id) as current_amount
    FROM "transaction"
    WHERE symbol = $1 AND user_id = $2
),
R AS (
    SELECT 
        n,
        current_amount,
        price as running_total
    FROM A 
    WHERE n = 1
    UNION ALL 
    SELECT 
        A.n, 
        A.current_amount,
        CASE 
            WHEN A.side = 'buy' THEN (R.current_amount * R.running_total + A.amount*A.price)/(R.current_amount+A.amount)
            ELSE running_total
        END as running_total
    FROM R
        JOIN A 
            ON A.n = R.n+1
)

SELECT cast(R.running_total as BIGINT) as cost_basis FROM R
ORDER BY n DESC
LIMIT 1;