# CHANGELOG FOR 2.X
---

### v2.0.0 (2020-10-14)

**This version is not backward compatible with v1.0.0.**

##### Details

* Dropped support form GORM v1
* Added support for GORM v2
* `Paginator` is not an interface and most of the methods return an error:

    ```go 
    type Paginator interface {
        SetPage(page int)
        Page() (int, error)
        Results(data interface{}) error
        Nums() (int64, error)
        HasPages() (bool, error)
        HasNext() (bool, error)
        PrevPage() (int, error)
        NextPage() (int, error)
        HasPrev() (bool, error)
        PageNums() (int, error)
    }
    ```
* All the methods of `Viewer`'s interface return an error
