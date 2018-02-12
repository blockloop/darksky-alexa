# Tea

Tea is a utility library intended to help improve the flow of your Go HTTP servers. Tea goes well with [Chi](https://github.com/go-chi/chi), but works well with the standard lib.

## Install

```bash
go get -u github.com/blockloop/tea
```

## Docs

[Godoc](https://godoc.org/github.com/blockloop/tea)

## Example

```go
package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/blockloop/tea"
)

func main() {
	r := chi.NewRouter()
	// use the JSON handler provided by tea
	r.NotFound(tea.NotFound)
	r.Get("/brewery/{id}", tea.Handler(GetBrewery))
	r.Post("/brewery", tea.Handler(PostBrewery))

	http.ListenAndServe(":3000", r)
}

func GetBrewery(w http.ResponseWriter, r *http.Request) (int, interface{}) {
        id, err := tea.URLInt(r, "id", "required,gt=0")
        if err != nil {
                return tea.StatusError(404)
        }

        brewery, err := db.GetBrewery(id)
        if err != nil {
                log.WithError(err).Error("failed to get brewery from db")
                return tea.StatusError(500)
        }

        return 200, brewery
}

type PostBreweryRequest struct {
        Name string `json:"name" validate:"required"`
        City string `json:"city" validate:"required"`
}

func PostBrewery(w http.ResponseWriter, r *http.Request) (int, interface{}) {
        var req PostBreweryRequest
        // parse the request body as JSON and validate it's struct fields
        err := tea.Body(r, &req)
        if err != nil {
                return tea.Error(400, err)
        }

        // create with db
        if err != nil {
                log.WithError(err).Error("failed to create brewery")
                return tea.StatusError(500)
        }

        return 200, brewery
}
```


