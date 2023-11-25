package main

import (
	"Backend_Golang/internal/routers"
	"Backend_Golang/pkg"
	"log"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
	// fmt.Println("Init Berjalan")
}
func main() {
	database, err := pkg.PostgreSQLDB()
	if err != nil {
		log.Fatal(err)
	}
	routers := routers.New(database)
	server := pkg.Server(routers)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
