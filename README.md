# GoFur
Web Framework, close to Laravel and Django.

> Based on [slick](https://github.com/anthdm/slick)


## Components
- Router: [httprouter](https://github.com/julienschmidt/httprouter)
- Session Management: scs
- Migrations: Goose
- SQL quries: sqlc (Maybe change to GORM)
- Temlates: templ
- CSS: TailwindCSS
- HTMX
- Console/CLI: cobra
- Default DB: Postgres (sqlc and goose is bad with other DBs)
- Auth: No candidate yet
- Events: no candidate
- Mails: no candidate
- Form/Validations: no candidate
- Static Assets
- Logging


For now you will need to install templ and air manually (working on making this work with `gofur install`)

- [https://github.com/a-h/templ/](https://github.com/a-h/templ/)
- [https://github.com/cosmtrek/air](https://github.com/cosmtrek/air)

Install the gofur cli
```
go install "github.com/anthdm/gofur/gofur@latest"
```

Create new gofur project
```
gofur new myapp
```

Install the project
```
cd myapp && gofur install
```

Start the project
```
gofur run
```

Run application in watch mode using [air](https://github.com/cosmtrek/air)
```
air
```
