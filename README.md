
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/dd935c6857884fd4ad0ae2a4886a9872)](https://app.codacy.com/gh/BlindGarret/echorend/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_coverage)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/BlindGarret/echorend/ci.yaml)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/dd935c6857884fd4ad0ae2a4886a9872)](https://app.codacy.com/gh/BlindGarret/echorend/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
![GitHub language count](https://img.shields.io/github/languages/count/BlindGarret/echorend)
![GitHub top language](https://img.shields.io/github/languages/top/BlindGarret/echorend)
![GitHub License](https://img.shields.io/github/license/BlindGarret/echorend)
[![GoDoc](https://godoc.org/github.com/BlindGarret/echorend?status.svg)](https://godoc.org/github.com/BlindGarret/echorend)

<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/BlindGarret/echorend">
    <img src="images/logo.png" alt="Logo" width="50">
  </a>

  <p align="center">
    EchoRend
    <br />
    <a href="https://github.com/BlindGarret/echorend/issues">Report Bug</a> |
    <a href="https://github.com/BlindGarret/echorend/issues">Request Feature</a> |
    <a href="https://pkg.go.dev/github.com/BlindGarret/echorend">Documentation</a>
  </p>
</p>

### Built With

* [Go](https://golang.org/)
* [Echo](https://echo.labstack.com/)
* [Raymond](https://github.com/aymerick/raymond) 

## Description

This is a simple library which wraps the functionality of gathering, parsing, and rendering Echo templates.

It's designed to make projects like server-side rendering with Go easier to manage.


## Why?

I've found myself writing a lot of small and medium sized projects recently which use Echo as an API framework and server side rendered HTML templates from go. The problem I found is, while it's fairly easy to start a project and write a few templates, the scaling complexity gets out of control fairly rapidly to handle the registration manually. While it was a relatively simple problem to dynamically load in the templates, there were just enough edgecases that I would find subtle bugs in implementations.

With a more robust approach to the code, it quickly became enough code/tests that it made no sense to keep copying around between projects, so I decided to make a library out of it.

Is it a perfect solution? No. But it's a solution that works for me and I hope it works for you too. At least for smaller/medium projects, where keeping the parsed templates in memory is reasonable, it should be more than enough to handle the complexity of the templates.

## Usage

### Installation

1. go get github.com/BlindGarret/echorend

### Example

The following is a simple example of usage with the Handlebars renderer and the glob gatherer.  This example assumes you have a directory structure like the following:

```
main.go
views/
  index.hbs
partials/
  test_component.hbs
```
It sets up a simple echo server which renders index and test_component templates.

```go   
package main

import (
        "fmt"
        "net/http"

        "github.com/BlindGarret/echorend/gatherers/glob"
        "github.com/BlindGarret/echorend/renderers/handlebars"
        "github.com/labstack/echo/v4"
)

func main() {
        viewDir := "views"
        partialsDir := "partials"
        viewGatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
                TemplateDir:     &viewDir,
                IncludeTLDInKey: false,
                Extensions:      []string{".hbs"},
        })

        partialGatherer := glob.NewGlobGatherer(glob.GlobGathererConfig{
                TemplateDir:     &partialsDir,
                IncludeTLDInKey: false,
                Extensions:      []string{".hbs"},
        })

        renderer := handlebars.NewHandlebarsRenderer(viewGatherer, partialGatherer)
        renderer.MustSetup()
        errs := renderer.CheckRenders()
        if len(errs) > 0 {
                for _, err := range errs {
                        fmt.Println(err)
                }
                return
        }

        e := echo.New()
        e.Debug = true
        e.Renderer = renderer
        e.GET("/", func(c echo.Context) error {
                return c.Render(http.StatusOK, "index", map[string]interface{}{
                        "Title": "Hello, World!",
                })
        })

        e.GET("/partial", func(c echo.Context) error {
                return c.Render(http.StatusOK, "test_component", map[string]interface{}{
                        "Title": "Just the partial now",
                })
        })

        e.Logger.Fatal(e.Start(":1323"))
}

```

## Handlebars

### Partials
Partials are parsed and registered with the Raymond Library. This is stored statically, which leads to a couple of interesting caveats.

1. Partial Registration is global. Even though the renderer is instantiated as a class.
    - This means if you register a partial with the same name in two different renderers, the second renderer will throw an error.
    - This is particularly important for testing, where renderers are setup constantly.
2. Partials are also registered as view.
    - This is a convience issue, as there are often times you want to define a "component like" partial where you reuse it multiple places, but you also may want to render it by itself for something like an AJAX request.


