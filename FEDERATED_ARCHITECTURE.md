# Federated Table Architecture

## Overview

The system has been refactored from raw MySQL packet proxying to using MySQL **Federated Tables**. This design allows users to modify a selected table locally while transparently querying all other tables from the remote database.

## Previous Architecture (Packet Proxying)

```
Client
  ↓ (raw MySQL packets)
Proxy Server (intercepts protocol)
  ├─ Routes to Local DB (selected table with data)
  └─ Routes to Remote DB (all other tables)
```

**Issues with previous approach:**
- Complex packet-level handling
- Manual data replication for all tables
- Difficult to maintain consistency
- Size limitations on data that could be replicated

## New Architecture (Federated Tables)

```
Client
  ↓ (SQL queries)
Local MySQL Instance
  ├─ Local Table (non-federated, editable)
  │   └── User can INSERT/UPDATE/DELETE
  └─ Federated Tables (point to remote)
      └── Transparent passthrough to remote server
          
Remote MySQL Server
  └── Provides read-only data for federated tables
```

## How It Works

### 1. Setup Phase

When the application starts:

```go
// User selects schema and table
Selected Schema: customer_db
Selected Table: orders (will be local, editable)

// Application:
1. Creates fresh local database (customer_db)
2. Creates 'orders' table locally with same schema
3. Creates all other tables as FEDERATED tables:
   - FEDERATED -> customers (remote)
   - FEDERATED -> products (remote)
   - FEDERATED -> invoices (remote)
   etc.
```

### 2. Local Table (Selected)

The selected table is created as a **normal InnoDB table**:

```sql
-- Created locally in the Docker container
CREATE TABLE orders (
  id INT AUTO_INCREMENT PRIMARY KEY,
  customer_id INT,
  order_date DATETIME,
  total_amount DECIMAL(10,2),
  -- all columns from the remote schema
) ENGINE=InnoDB;
```

**User permissions:**
- ✅ SELECT - Read local data
- ✅ INSERT - Add new records
- ✅ UPDATE - Modify existing records
- ✅ DELETE - Remove records

### 3. Federated Tables (All Others)

Other tables are created as **Federated tables** pointing to the remote server:

```sql
-- Created locally but reads from remote
CREATE TABLE customers (
  -- Same schema as remote
) ENGINE=FEDERATED
CONNECTION='mysql://remote_user:remote_pass@remote_host:3306/customer_db/customers';
```

**User permissions:**
- ✅ SELECT - Read from remote
- ❌ INSERT/UPDATE/DELETE - Not supported through federation

This is intentional - users can only modify the selected local table.

### 4. User Query Execution

When a user executes queries, they all go to the **local MySQL instance**:

#### Reading from Federated Table
```sql
-- User query
SELECT * FROM products WHERE id = 100;

-- Automatically:
Local MySQL -> (FEDERATED) -> Remote MySQL
              <- Returns data <-
```

#### Modifying Local Table
```sql
-- User query
INSERT INTO orders (customer_id, order_date, total_amount)
VALUES (5, NOW(), 99.99);

-- This data is stored ONLY locally
-- Remote server is unaffected
```

#### Joining Local and Remote Data
```sql
-- This works seamlessly!
SELECT o.id, o.total_amount, c.name
FROM orders o
FEDERATED JOIN customers c ON o.customer_id = c.id;

-- Execution:
1. Fetch 'orders' from local table
2. Fetch 'customers' from remote (via FEDERATED)
3. Join in local MySQL instance
```

## Benefits

1. **No Data Replication** - No need to copy entire tables locally
2. **Scalability** - Works with large remote databases
3. **Transparent Access** - Users see one unified database
4. **Local Modifications** - Can edit one table without affecting remote
5. **Consistency** - Remote data is always current
6. **Simplicity** - Standard MySQL queries work unchanged
7. **No Packet Handling** - No complex proxy logic needed

## Implementation Changes

### Files Modified

#### 1. `query.go`
Added:
- `FEDERATED_ENGINE_CREATE` - Query template for creating Federated tables
- `GET_ALL_TABLES_EXCEPT` - Query to get all tables except selected one

#### 2. `database.go`
Added:
- `CreateLocalTableSchema()` - Creates the selected table locally
- `CreateFederatedTable()` - Creates Federated tables for other tables

#### 3. `main.go`
Refactored setup flow:
- Removed: Complex topological sort and data replication logic
- Removed: Table size estimation checks
- Added: Simple Federated table creation loop

**Old flow (many steps):**
```
Topological sort → Calculate sizes → Fetch all data → Insert locally → Start proxy
```

**New flow (simple):**
```
Create local table → Create Federated tables → Start proxy
```

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                    User Application                      │
│              (Connects to localhost:5432)                │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ↓
         ┌─────────────────────────────┐
         │   Local MySQL (Docker)      │
         │                             │
         │  Database: customer_db      │
         │  ├─ orders (InnoDB)         │◄─── User can modify
         │  │  - id                    │
         │  │  - customer_id           │
         │  │  - total_amount          │
         │  │                          │
         │  └─ FEDERATED TABLES:       │
         │     ├─ customers            │
         │     ├─ products             │
         │     └─ invoices             │
         │        ↓ (tunneled)         │
         └─────────┼──────────────────┘
                   │
            FEDERATED tunnels
                   │
         ┌─────────↓──────────────────┐
         │  Remote MySQL Server       │
         │                            │
         │  Database: customer_db     │
         │  ├─ orders (read-only)     │
         │  ├─ customers              │
         │  ├─ products               │
         │  └─ invoices               │
         └────────────────────────────┘
```

## Querying Examples

### Example 1: Read from Remote via Federated
```go
// Application setup
selectedTable := "orders"
remoteDB := "customer_db"

// User connects to local MySQL and executes:
// SELECT * FROM customers WHERE id = 1;

// Behind the scenes:
// Local MySQL detects 'customers' is FEDERATED
// Opens connection to remote MySQL
// Executes query on remote
// Returns results to user
// User never knows it came from remote!
```

### Example 2: Modify Local Table
```go
// User connects to local MySQL and executes:
// INSERT INTO orders VALUES (NULL, 5, NOW(), 100.00);

// Behind the scenes:
// Local MySQL inserts into local 'orders' table
// Data is stored only locally
// Remote 'orders' table is NOT affected
// Next time user queries remote orders, they won't see this insert
// But they WILL see it in local orders
```

### Example 3: Complex Query with Join
```go
// User executes:
// SELECT o.id, o.total_amount, c.name, c.email
// FROM orders o
// JOIN customers c ON o.customer_id = c.id;

// MySQL execution plan:
// 1. Fetch all 'orders' from local table
// 2. For each order, fetch customer from FEDERATED 'customers'
// 3. Join locally and return results
```

## Configuration

Federated tables are created with connection string:
```sql
CONNECTION='mysql://user:pass@host:port/database/table'
```

Uses environment variables:
- `REMOTE_DB_HOST` - Remote server hostname
- `REMOTE_DB_PORT` - Remote server port
- `REMOTE_DB_USER` - Remote MySQL user
- `REMOTE_DB_PASS` - Remote MySQL password  
- `REMOTE_DB_NAME` - Remote database name

## Security Notes

1. **Federated Connection Credentials** - Stored in MySQL (in `mysql.servers` table)
2. **Network Security** - Federated tables use TCP connection to remote server
3. **User Isolation** - Each Federated table connection uses remote credentials configured
4. **Local Data Privacy** - Local modifications are isolated from remote database

## Known Limitations

1. **No FOREIGN KEY constraints** between local and Federated tables (MySQL limitation)
2. **Transaction scope** - Transactions don't span local and remote tables atomically
3. **Full-text search** - Limited support with Federated tables
4. **Full outer joins** - Not supported with Federated tables

## Future Enhancements

1. **Sync Module** - Optionally sync local changes back to remote
2. **Conflict Resolution** - Handle conflicts if remote data changes
3. **Partitioned Tables** - Selectively federate by partition
4. **Caching Layer** - Cache Federated data locally for performance
5. **Change Data Capture** - Track changes to the local table for auditing

## References

- [MySQL Federated Storage Engine](https://dev.mysql.com/doc/refman/8.0/en/federated-storage-engine.html)
- [Creating Federated Tables](https://dev.mysql.com/doc/refman/8.0/en/create-table.html)
