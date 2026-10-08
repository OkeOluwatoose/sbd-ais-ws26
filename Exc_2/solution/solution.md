# Part 2 — Activities (graded)

## Activity 2.1 — PostgreSQL Analytical Queries (E-commerce)

>Generate the dataset (in another terminal, from the `ecommerce` folder):

>Load the generated data into PostgreSQL in a **new table** called `orders`.
Write the `CREATE TABLE` yourself (choose sensible types for each column) and
use the same `\COPY` pattern as in section 3.3.

### Loading code Solution
```sql
DROP TABLE IF EXISTS orders;

CREATE TABLE orders (
  id SERIAL PRIMARY KEY,
  customer_name TEXT,
  product_category TEXT,
  quantity INTEGER,
  price_per_unit NUMERIC(10,2),
  order_date DATE,
  country TEXT
);

\COPY orders(customer_name,product_category,quantity,price_per_unit,order_date,country) FROM '/data/orders_1M.csv' DELIMITER ',' CSV HEADER;
```

Using SQL ([list of supported SQL commands](https://www.postgresql.org/docs/current/sql-commands.html)),
answer the following questions:

> **A.** Which order has the highest `price_per_unit`?
#### query 
```sql
SELECT * FROM orders
ORDER BY price_per_unit DESC
LIMIT 1;
```
Result:
```bash
   id   | customer_name | product_category | quantity | price_per_unit | order_date | country 
--------+---------------+------------------+----------+----------------+------------+---------
 841292 | Emma Brown    | Automotive       |        3 |        2000.00 | 2024-10-11 | Italy
```

> **B.** What are the top 3 product categories with the highest total quantity
sold across all orders?
#### query
```sql
SELECT product_category, SUM(quantity) AS total_quantity
FROM orders
GROUP BY product_category
ORDER BY total_quantity DESC
    LIMIT 3;
```
Result:
```bash
 product_category | total_quantity 
------------------+----------------
 Health & Beauty  |         300842
 Electronics      |         300804
 Toys             |         300598
```

> **C.** What is the total revenue per product category?
(Revenue = `price_per_unit × quantity`)
#### query
```sql
SELECT
    product_category,
    SUM(price_per_unit * quantity) AS total_revenue
FROM orders
GROUP BY product_category
ORDER BY total_revenue DESC;
```
Result:
```bash
 product_category | total_revenue 
------------------+---------------
 Automotive       |  306589798.86
 Electronics      |  241525009.45
 Home & Garden    |   78023780.09
 Sports           |   61848990.83
 Health & Beauty  |   46599817.89
 Office Supplies  |   38276061.64
 Fashion          |   31566368.22
 Toys             |   23271039.02
 Grocery          |   15268355.66
 Books            |   12731976.04
```

> **D.** Who are the top 5 customers by total spending?
#### query
```sql
SELECT
    customer_name,
    SUM(price_per_unit * quantity) AS total_spending
FROM orders
GROUP BY customer_name
ORDER BY total_spending DESC
    LIMIT 5;
```
Result:
```bash
 customer_name  | total_spending 
----------------+----------------
 Carol Taylor   |      991179.18
 Nina Lopez     |      975444.95
 Daniel Jackson |      959344.48
 Carol Lewis    |      947708.57
 Daniel Young   |      946030.14
```


> **E.** Look at the spending totals in D — and at how many orders each of those
customers has. What do you notice? Open `ecommerce/dataset_generator.py` and
explain *why* the data looks like this.
#### query
```sql
SELECT
    customer_name,
    COUNT(*) AS number_of_orders,
    SUM(price_per_unit * quantity) AS total_spending
FROM orders
GROUP BY customer_name
ORDER BY total_spending DESC
    LIMIT 20;
```
Result:
```bash
 customer_name  | number_of_orders | total_spending 
----------------+------------------+----------------
 Carol Taylor   |             1028 |      991179.18
 Nina Lopez     |              980 |      975444.95
 Daniel Jackson |             1033 |      959344.48
 Carol Lewis    |              943 |      947708.57
 Daniel Young   |              973 |      946030.14
...
```
#### findings
We notice that the number of orders are pretty evenly distributed (each name has a number of orders between 900 and 1000.)

This even distribution of values is typical of synthetic datasets generated using random assignments, 
and when we look at the `dataset_generator.py` file, we can see that 
it does indeed use a random choice to assign customer names and orders.

## Activity 2.2 — Why Is This Self-Join So Slow?

Users of the system sometimes run naive queries such as:

```sql
SELECT COUNT(*)
FROM people_big p1
JOIN people_big p2
  ON p1.country = p2.country;
```

> **Do not run it on the full table in class** — use the smaller tables below.
> (If you are curious, run it at home and let it finish.)

**Step 1 — Measure how it grows.** Create three smaller copies of the table:

```sql
CREATE TABLE people_50k  AS SELECT * FROM people_big WHERE id <= 50000;
CREATE TABLE people_100k AS SELECT * FROM people_big WHERE id <= 100000;
CREATE TABLE people_200k AS SELECT * FROM people_big WHERE id <= 200000;
```

> Run the self-join (with `\timing on`) on each of the three tables and fill in:

| rows in table | join result (`COUNT(*)`) | time |
|---|---|---|
| 50 000 | 27501822 | 1771.973 ms (00:01.772) |
| 100 000 | 109946508 | 8867.400 ms (00:08.867) |
| 200 000 | 439395606 |41880.844 ms (00:41.881) |

When the input **doubles**, by what factor do the result and the time grow?
Use this to **predict** the result size and the runtime on `people_big` (1M rows).

#### Result size
From 50,000 → 100,000 rows:
$$ \frac{109,946,508}{27,501,822} \approx 4.00 $$
From 100,000 → 200,000:
$$ \frac{439,395,606}{109,946,508} \approx 4.00 $$

The join result grows by approximately a factor of 4 when the input doubles.
This makes sense for a self-join because the number of possible pairs is roughly proportional to $n^2$

#### Runtime size
From 50,000 → 100,000:
$$ \frac{8.867}{1.772} \approx 5.00 $$
From 100,000 → 200,000:
$$ \frac{41.881}{8.867} \approx 4.72 $$
So the runtime grows by approximately 4.7–5x when the input doubles.

#### Predicting `people_big`
From 200,000 → 1,000,000 rows is a factor of 5.

##### Result size
Since the result grows approximately quadratically:
$$ 439,395,606 \times 5^2 $$ $$ = 10,984,890,150 $$

Thus, we can predict approximately 10.98 billion rows in the join result.

##### Runtime size
The runtime grows with a factor of 4.72 to 5 for each doubling. 
That is not easy to approximate, so we can estimate the runtime quadratically using the 
200k runtime:

$$ 41.881 \times \left(\frac{1,000,000}{200,000}\right)^2 $$ $$ =41.881 \times 25 $$ $$ =1,047.025\text{ seconds} $$

That's approximately 17.45 minutes

**Step 2 — Does an index help?** Create an index on `country` of `people_100k`,
run `ANALYZE people_100k;`, and repeat the query. Did the time change? Use
`EXPLAIN ANALYZE` to support your answer. Why does (or doesn't) the index help?

Regular run
```sql
EXPLAIN ANALYZE
SELECT COUNT(*)
FROM people_100k p1
JOIN people_100k p2
  ON p1.country = p2.country;
```

Result
```bash
Finalize Aggregate  (cost=637894.46..637894.47 rows=1 width=8) (actual time=29473.795..29518.971 rows=1.00 loops=1)
   Buffers: shared hit=2048
   ->  Gather  (cost=637894.24..637894.45 rows=2 width=8) (actual time=29473.656..29518.838 rows=2.00 loops=1)
         Workers Planned: 1
         Workers Launched: 1
         Buffers: shared hit=2048
         ->  Partial Aggregate  (cost=636894.24..636894.25 rows=1 width=8) (actual time=29315.107..29315.114 rows=1.00 loops=2)
               Buffers: shared hit=2048
               ->  Parallel Hash Join  (cost=2347.54..475056.17 rows=64735229 width=0) (actual time=925.799..22121.372 rows=54973254.00 loops=2)
                     Hash Cond: (p1.country = p2.country)
                     Buffers: shared hit=2048
                     ->  Parallel Seq Scan on people_100k p1  (cost=0.00..1612.24 rows=58824 width=8) (actual time=0.018..29.185 rows=50000.00 loops=2)
                           Buffers: shared hit=1024
                     ->  Parallel Hash  (cost=1612.24..1612.24 rows=58824 width=8) (actual time=923.811..923.812 rows=50000.00 loops=2)
                           Buckets: 131072  Batches: 1  Memory Usage: 5248kB
                           Buffers: shared hit=1024
                           ->  Parallel Seq Scan on people_100k p2  (cost=0.00..1612.24 rows=58824 width=8) (actual time=763.333..885.029 rows=50000.00 loops=2)
                                 Buffers: shared hit=1024
 Planning Time: 5.413 ms
 JIT:
   Functions: 24
   Options: Inlining true, Optimization true, Expressions true, Deforming true
   Timing: Generation 14.241 ms (Deform 6.010 ms), Inlining 411.462 ms, Optimization 552.964 ms, Emission 561.901 ms, Total 1540.568 ms
 Execution Time: 29533.778 ms
```
Creating Index
```sql
CREATE INDEX idx_people_100k_country
ON people_100k(country);
```
Running analyze on `people_100k`
Result
```bash
Time: 159.364 ms
```
Running query again with index.
Results
```bash
Finalize Aggregate  (cost=640248.27..640248.28 rows=1 width=8) (actual time=16903.042..16905.872 rows=1.00 loops=1)
   Buffers: shared hit=95 read=91
   ->  Gather  (cost=640248.05..640248.26 rows=2 width=8) (actual time=16903.018..16905.853 rows=2.00 loops=1)
         Workers Planned: 1
         Workers Launched: 1
         Buffers: shared hit=95 read=91
         ->  Partial Aggregate  (cost=639248.05..639248.06 rows=1 width=8) (actual time=16813.394..16813.397 rows=1.00 loops=2)
               Buffers: shared hit=95 read=91
               ->  Parallel Hash Join  (cost=2332.10..477378.41 rows=64747855 width=0) (actual time=261.946..11094.502 rows=54973254.00 loops=2)
                     Hash Cond: (p1.country = p2.country)
                     Buffers: shared hit=95 read=91
                     ->  Parallel Index Only Scan using idx_people_100k_country on people_100k p1  (cost=0.29..1596.51 rows=58824 width=8) (actual time=0.017..26.553 rows=50000.00 loops=2)
                           Heap Fetches: 0
                           Index Searches: 1
                           Buffers: shared hit=93
                     ->  Parallel Hash  (cost=1596.51..1596.51 rows=58824 width=8) (actual time=259.316..259.317 rows=50000.00 loops=2)
                           Buckets: 131072  Batches: 1  Memory Usage: 5248kB
                           Buffers: shared hit=2 read=91
                           ->  Parallel Index Only Scan using idx_people_100k_country on people_100k p2  (cost=0.29..1596.51 rows=58824 width=8) (actual time=1.266..19.460 rows=50000.00 loops=2)
                                 Heap Fetches: 0
                                 Index Searches: 1
                                 Buffers: shared hit=2 read=91
 Planning:
   Buffers: shared hit=9 read=1
 Planning Time: 1.648 ms
 JIT:
   Functions: 16
   Options: Inlining true, Optimization true, Expressions true, Deforming true
   Timing: Generation 3.974 ms (Deform 0.323 ms), Inlining 220.804 ms, Optimization 65.578 ms, Emission 96.025 ms, Total 386.380 ms
 Execution Time: 16909.709 ms
```
Analysis:

|  | Before index | After index |
|---|---|---|
| Planning time | 5.413ms | 1.648 ms |
| Execution time | 29533.778 ms | 16909.709 ms |
| Scan type | Seq Sacn |Index Scan |

Conclusion:
Yes. An index helps in this case. 
The execution time decreased after creating the index. `EXPLAIN ANALYZE` shows that PostgreSQL used 
the index to locate rows matching the `country` condition instead of scanning the entire table. 

**Step 3 — Rewrite it.** The query only wants the *number* of matching pairs,
not the pairs themselves. If a country has *k* people, how many pairs does it
contribute to the join? Write a query that computes the **same number
without a join**. Check that it returns exactly the same result as the join on
`people_100k`, then run it on `people_big` and compare its runtime with your
prediction from Step 1.

Answer:
If a country has k people, then every person can match with every other
person including themselves. So the number of matches are:
$$ k \times k = k^2 $$
So if Germany has 100 people, it contributes $$ 100^2 = 10,000 $$

To re-write the query without a join, we can first count how many people are in each country, square that,
and then add the results up. i.e.

```sql
SELECT SUM(country_count * country_count) AS matching_pairs
FROM (
    SELECT country, COUNT(*) AS country_count
    FROM people_100k
    GROUP BY country
) AS country_counts;
```
Result from query:
```bash
 109946508
```
(Returns the same as the join on `people_100k`)

Results running on `people_big`
```bash
10983941260
(1 row)

Time: 1408.804 ms (00:01.409)
```
Comparing results

|  | Prediction (with join) | Actual (without join)|
|---|---|---|
| Result size | 10,984,890,150 | 10,983,941,260|
| Execution time | 1,047,025ms | 1,408.804 ms |

Our predicted result size is quite close to the actual, however, the execution time is much shorter 
in the case of the query without a join.


**Step 4 — Discussion (submit in writing).** Considering **scalability** and
**efficiency**, which approaches and/or optimizations can be applied to improve
this kind of query in a real system? Discuss at least:

- what the rewrite in Step 3 tells you about adding more hardware or an index;

Answer:

The initial query performs a self-join, and then counts the number of records in that join.

However, we can calculate that same figure by counting how many people are in each country, squaring that number, 
and then adding up the results. This is much more efficient because we never generate the potentially enormous set of matching pairs.

It shows that adding more hardware or an index isn't necessarily the best solution.

An index on `country` might improve the join in some situations, but it does not eliminate the fundamental problem: the join 
can produce a huge number of matching pairs. The rewrite changes the algorithm itself, reducing the amount of work required.

- what you would do if the business actually needed **the pairs themselves**
  (not just their count) — would a bigger machine or a cluster help, and how
  much?

Answer:
  
If the business needed every pair of people from the same country, we cannot avoid generating those pairs.
```bash
1,000 people → 1,000,000 pairs
10,000 people → 100,000,000 pairs
100,000 people → 10,000,000,000 pairs
1,000,000 people → 1,000,000,000,000 pairs
```
As you can see, the number of pairs grows quadratically. A bigger machine can help to some extent
through more CPU, RAM and faster storage, but it doesn't change the quadratic growth.

A cluster or distributed processing system would split the work across multiple machines.


- the limits of an **OLTP database** for this workload, especially in a
  **large-scale cloud environment**.

An OLTP database is designed primarily for transactional workloads such as:

- inserting and updating individual records;
- processing customer transactions;
- looking up individual customers;
- maintaining data consistency;
- handling many concurrent users.

In a large-scale cloud environment, running a large self-join query that produces billions of rows  directly against a production 
OLTP database can cause problems such as:

- high CPU and memory consumption;
- increased disk I/O;
- contention with normal transactions;
- slower response times for application users;
- large temporary/intermediate results;
- potentially very high cloud infrastructure costs.

An index can help selective queries, but it cannot solve the fundamental problem of generating an enormous join result. 
For some large analytical queries, PostgreSQL may even choose a sequential scan instead of an index because accessing 
a large proportion of the table through an index isn't cheaper.

<!-- > **Optional:** support your answer with a diagram, SQL or code. -->

## Submission

Commit to the Github repository and add the link to the Github solution in the Moodle submission:

- **2.1:** your `CREATE TABLE` statement and the SQL for A–D with their
  results, and your answer to E;
- **2.2:** the table from Step 1 with your prediction, your observations from
  Steps 2–3 (including your rewrite query), and the written discussion from
  Step 4.

Be ready to shortly present your solutions (5–8 minutes) in the next exercise
session (05.11.2026).
