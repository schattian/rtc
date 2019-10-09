# Change

Any change could be reduced to one of the following mutually-exclusive types:

|   | Create | Retrieve | Update | Delete |
|---|---|---|---|---|
| Table | Must | Must | Must | Must |
| Column | Must | Must | Must |  No |
| Value | Must | No | Must | No  |
| ID | No | Can | Must | Must |

## Non-imperative

There are three values which *Can* be filled. Those are:

### Retrieve

- ID: if **blank**, retrieves **all the entities under the table**. Otherwise, the specified entity_id.
- Column: if **blank**, retrieves **all the possibly columns**. Otherwise, the specified column value.

> Note: **take in consideration that if, ID is blank, column MUST be ZERO**

### Create

- Columns: if **blank**, the push will create an entity with **multiple fields at one time** (as it only accepts embeddings, it requires that the **value needs to be bytes or json**). Otherwise, the push will send any column different from declared as blank.
