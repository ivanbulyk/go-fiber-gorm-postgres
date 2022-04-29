# A simple API that  returns slowest connected database queries.
Uses Fiber and PostgreSQL.
Features:
- Supports pagination(`http://localhost:8080/api/pages?page=2&page_size=10`)
- Support filterting by SELECT,INSERT,UPDATE,DELETE(`http://localhost:8080/api/sql/SELECT`)
- Support sorting by time spent(`http://localhost:8080/api/queries_asc`, `http://localhost:8080/api/queries`(desc by default))
- 80%+ test coverage
Databbase interaction:
- Create(POST `http://localhost:8080/api/create_books`)
- Delete(DELETE `http://localhost:8080/api/delete_books/34`)
- Get all(GET `http://localhost:8080/api/books`)
- Get by ID(GET`http://localhost:8080/api/get_books/10`)

To run the app: `make` in project root