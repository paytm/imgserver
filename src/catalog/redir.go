/*
 Package catalog implements a middleware to redirect to thumbnail image, given paytm sku and product id
*/
package catalog

import (
  "net/http"
  "strconv"
  "log"
  "strings"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type HandlerFunc func(rw http.ResponseWriter, r *http.Request)
func ImageRedir(dsn string) (HandlerFunc) {
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    log.Fatal("db error ", err.Error())
  }

  db.SetMaxIdleConns(5)

  return func (w http.ResponseWriter, r* http.Request) {
    var name string
    fields := strings.Split(r.URL.Path,"/")

    if len(fields) < 4 {
      http.Error(w, "Bad Path", http.StatusBadRequest)
      return
    }

    sku := fields[2]
    product_id := strings.TrimSuffix(fields[3],".jpg")

    _,err := strconv.Atoi(product_id)
    if err != nil {
      http.Error(w, "Bad Product Id", http.StatusBadRequest)
      return
    }

    log.Println("fetching for sku,product ",sku,product_id)

    // fetch the thumbnail by default
    //err = db.QueryRow("SELECT value from catalog_product_resource where product_id = ? and is_default = ?",product_id,2).Scan(&name)
    //err = db.QueryRow(`select paytm_sku,catalog_product_resource.value from catalog_product join catalog_product_resource
    //           on catalog_product.id = catalog_product_resource.product_id and is_default = ? and catalog_product.id = ?`,2,product_id).Scan(&sku,&name)
    err = db.QueryRow("select paytm_sku,thumbnail from catalog_product where id = ?",product_id).Scan(&sku,&name)
    if err != nil {
      log.Println(err.Error())
      http.Error(w, "Bad Product Id", http.StatusNotFound)
      return
    }

    url := "http://assets.paytm.com/images/catalog/product/" + sku[0:1] + "/" + sku[0:2] + "/" + sku + "/210x210/" + name
    log.Println("Redirecting to ",url)
    http.Redirect(w,r,url,302)
  }
}
