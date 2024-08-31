## sqlc-joins-gen

> a code generator to make relational queries easier for [sqlc](https://sqlc.dev/).
> very similar to [drizzle-orm's query api](https://orm.drizzle.team/docs/rqb).

### install

```sh
go install github.com/lqr471814/sqlc-joins-gen@latest
```

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

here's an example query config and `schema.sql` file.

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
    // Method - Define a relational query
    {
        // name of the function generated for the method
        name: "getAuthors", 
        // whether to return first or many results
        return: "first",
        // the SQL table to target at the top level
        table: "Author",
        // Query - Define options for the select statement
        query: {
            columns: {
                // only include the `name` column + primary keys
                // this will exclude `age` & `alive`
                name: true 
                // primary keys are always included in the underlying select statement
                // due to implementation details
            },
            // define which tables to join with
            with: {
                // Query - Define a joined table and options for its select statement
                Book: { 
                    // common select options
                    orderBy: {
                        name: "dsc"
                    },
                    limit: 5,
                    // you could keep nesting Queries "Book" had another relation as well.
                    with: {}
                },
            },
            // common select options
            // note that interpolation of arguments in where must take the form of
            // $<argname>:<int | str | float | bool>(?)([])
            // ex: $age:int -> func (..., age int) ...
            // ex: $name:str? -> func (..., name sql.NullString) ...
            // ex: $emails:str[] -> func (..., name []string) ...
            // ex: $null_emails:str?[] -> func (..., name []sql.NullString) ...
            // passing an empty slice or `nil` into an arg with an array type
            // will be treated as `null`
            where: "Author.age > $age:int",
            orderBy: {
                age: "asc",
            },
            offset: 4,
        },
    },
]
```

this would produce a select statement that looks like this:

```sql
select
"Author"."name" as "Author_name",
"Author"."id" as "Author_id",
"Book"."id" as "Book_id",
"Book"."authorId" as "Book_authorId",
"Book"."name" as "Book_name"
from "Author"
inner join (select * from "Book"  order by "Book"."name" dsc limit 5) as "Book" on "Book"."authorId" = "Author"."id"
where Author.age > $age:int
order by "Author"."age" asc
offset 4
```

as well resulting go code that looks like this:

```go
// Table: Author
type GetAuthors struct {
	Name string
	Id   int
	Book []GetAuthors0
}

// Table: Book
type GetAuthors0 struct {
	Id       int
	AuthorId int
	Name     string
}

// the actual query method is appended to the "Queries" struct generated by sqlc
func (q *Queries) GetAuthors(ctx context.Context, age int) (GetAuthors, error) {
    ...
}
```

