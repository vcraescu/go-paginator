# Paginator

[![Go Report Card](https://goreportcard.com/badge/github.com/vcraescu/go-paginator?kill_cache=2)](https://goreportcard.com/report/github.com/vcraescu/go-paginator) 
[![Build Status](https://travis-ci.com/vcraescu/go-paginator.svg?branch=master&kill_cache=2)](https://travis-ci.com/vcraescu/go-paginator) 
[![Coverage Status](https://coveralls.io/repos/github/vcraescu/go-paginator/badge.svg?branch=master&kill_cache=2)](https://coveralls.io/github/vcraescu/go-paginator?branch=master)

A simple way to implement pagination in Golang.

## Usage

```go
package main

import (
    "fmt"
    "github.com/vcraescu/go-paginator"
    "github.com/vcraescu/go-paginator/adapter"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "time"
)

type Post struct {
	ID          uint `gorm:"primary_key"`
	Title       string
	Body        string
	PublishedAt time.Time
}

func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connection error: %s", err))
	}

	if err := db.AutoMigrate(&Post{}); err != nil {
        panic(err)
    }
	
	var posts []Post
	
	q := db.Model(Post{}).Where("published_at > ?", time.Now())
	p := paginator.New(adapter.NewGORMAdapter(q), 10)
	p.SetPage(2)
	
	if err = p.Results(&posts); err != nil {
		panic(err)
	}
	
	for _, post := range posts {
		fmt.Println(post.Title)
	}
}
```

Some of other methods available:

```go
p.HasNext()
p.HasPrev()
p.HasPages()
p.PageNums()
```

## Adapters

An adapter must implement the `Adapter` interface which has 2 methods: 

* **Nums** - must return the number of results;
* **Slice** - must retrieve a slice for an offset and length.

This way you can create your own adapter for any kind of data source you want to paginate. 

```golang 
type Adapter interface {
	Nums() (int64, error)
	Slice(offset, length int, data interface{}) error
}
```

### GORM Adapter

To paginate a **GORM** query builder.

```go
q := db.Model(Post{}).Where("published_at > ?", time.Now())
p := paginator.New(adapter.NewGORMAdapter(q), 10)
```

### Slice adapter

To paginate a slice.

```go
var pages []int
p := paginator.New(adapter.NewSliceAdapter(pages), 10)
```

## Views

View models contains all necessary logic to render the paginator inside a template.

### DefaultView

Use it if you want to render a paginator similar to the one from Google search.

**< Prev** 2 3 4 5 6 **7** 8 9 10 11 **Next >**

```go
func main() {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connection error: %s", err))
	}

	db.AutoMigrate(&Post{})
	
	var posts []Post
	
	q := db.Model(Post{}).Where("published_at > ?", time.Now())
	p := paginator.New(adapter.NewGORMAdapter(q), 10)
	p.SetPage(7)
	
	view := view.New(&p)
	
    pages, _ := view.Pages()
	fmt.Println(pages) // [2 3 4 5 6 7 8 9 10 11]
    
    next, _ := view.Next()
	fmt.Println(next) // 8
    
	prev, _ := view.Prev()
	fmt.Println(prev) // 6
    
    current, _ := view.Current()
	fmt.Println(current) // 7
}
```

## Changelog

* [v2.0.0](https://github.com/vcraescu/go-paginator/blob/v2.0.0/CHANGELOG-2.0.md)

## TODO

* More adapters

## License

Paginator is licensed under the [MIT License](LICENSE).
