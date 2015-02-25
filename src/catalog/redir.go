/*
 Package catalog implements a middleware to redirect to thumbnail image, given paytm sku and product id
*/
package catalog

import (
  "net/http"
  "strconv"
  "log"
  "fmt"
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
    var name , imagesize string
    var product_id_index, imagesize_index int

    product_id_index = 3
    imagesize_index = 3
    imagesize = "210x210"

    fields := strings.Split(r.URL.Path,"/")

    if len(fields) < 4 {
      http.Error(w, "Bad Path", http.StatusBadRequest)
      return
    }

    sku := fields[2]

    if len(fields) == 5 {
      product_id_index = product_id_index + 1
      imagesize = fields[imagesize_index]
    }

    product_id := strings.TrimSuffix(fields[product_id_index],".jpg")

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
	if len(name) > 0 {
invalid_url := fmt.Sprintf("http://assetscdn.paytm.com/images/catalog/brand/%s", name)
http.Redirect(w,r,invalid_url,301)
return
} else {

      http.Error(w, "Bad Product Id", http.StatusNotFound)
      return
    }
}

    url := fmt.Sprintf("http://%s/images/catalog/product/%s/%s/%s/%s/%s", "assets.paytm.com", sku[0:1], sku[0:2], sku, imagesize, name)
    log.Println("Redirecting to ",url)
    http.Redirect(w,r,url,302)
  }
}
