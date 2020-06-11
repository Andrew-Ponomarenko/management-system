package main

import (
    "fmt"
    "github.com/omc-college/management-system/pkg/ims/repository/postgres"
    "github.com/omc-college/management-system/pkg/ims/models"
    "log"
    "net/http"
)

func main() {
    rows,err:=dbquery.DbQuerier("roles")//passing which table you need to query from into the querier
    rls := make([]*models.Roles, 0)
    for rows.Next() {
        rl:= new(models.Roles)
        err := rows.Scan(rl.ID, rl.Name, rl.Entries)
        if err != nil {
            log.Fatal(err)
        }
        rls = append(rls,rl)
    }
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }

    for _,rl := range rls {
        fmt.Printf("%d, %s",rl.ID,rl.Name,rl.Entries)
    }

    log.Fatal(http.ListenAndServe(":8080", nil))
}
