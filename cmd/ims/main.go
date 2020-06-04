package main

import (
    "fmt"
    dbquery "github.com/omc-college/management-system/cmd/ims/dbq"
    "github.com/omc-college/management-system/pkg/rbac/models"
    "log"
    "net/http"
)

func main() {
    rows,err:=dbquery.DbQuerier("roles")//passing which table you need to query from into the querier
    rls := make([]*models.Role, 0)
    for rows.Next() {
        rl:= new(models.Role)
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