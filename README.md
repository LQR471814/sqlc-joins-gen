## sqlc-joins-gen

> a code generator to make relational queries easier for [sqlc](https://sqlc.dev/).
> very similar to [drizzle-orm's query api](https://orm.drizzle.team/docs/rqb).

### usage

some terminology for the config:

1. a `Method` refers to a top level definition of a relational query.
    1. if you're coming from drizzle, it is like the `db.query.<table>.findMany()` part of the query.
1. a `Query` refers to the options used in the query.
    1. if you're coming from drizzle, it is like the
        ```js
        {
            with: {
                posts: true      
            },
        }
        ```
        part of the query.

`sqlc-joins-gen` reads your `sqlc.yaml` config file in the current directory to know where to look for the schema (if you have a custom config location you can use `-config` to specify it).

to configure relational queries, you must define a [JSON5](https://json5.org/) config file.

this config file must have the name of the `queries` with the `.sql` extension removed and `.joins.json5` appended at the end.

```yaml
sql:
  - engine: "sqlite"
    queries: "query.sql" # queries should be defined at `query.joins.json5`
```

here's an example query config and the results it produces.

```sql
-- schema.sql (sqlite flavored)
create table Author (
    id INT NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    age INT,
    alive INT NOT NULL
);

create table Book (
    id INT NOT NULL PRIMARY KEY,
    authorId INT NOT NULL,
    name TEXT NOT NULL,
    FOREIGN KEY (authorId) REFERENCES Author(id)
);
```

```js
// query.joins.json5
[
    { // Method
        name: "getAuthors", // name of the function generated for the method
        return: "first", // whether to return first or many results
        table: "Author", // the SQL table to target at the top level
        query: { // Query
            columns: {
                name: true // only include the `name` column + primary keys

                // primary keys are always included in the underlying select statement
                // due to implementation details
            },
            with: {
                Book: {},
            },
            where: "Author.age > ?",
            orderBy: {
                age: "asc",
            }
        },
    },
]
```

