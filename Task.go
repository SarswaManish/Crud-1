package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Address struct {
	Id int          `json:id`
	StreetName string `json:streetName`
	City string			`json:city`
	State string        `json:state`
	CustomerId int      `json:customerId`
}

type Customer struct{
	Id int         `json:id`
	Name string    `json:name`
	Dob string     `json:dob`
	Add Address
}


func dateInSeconds(d1 string) int {
	d1Slice := strings.Split(d1, "/")

	newDate := d1Slice[2] + "/" + d1Slice[1] + "/" + d1Slice[0]
	myDate, err := time.Parse("2006/01/02", newDate)

	if err != nil {
		panic(err)
	}

	return int(time.Now().Unix() - myDate.Unix())
}


func getCustomer(w http.ResponseWriter,r *http.Request){

	db,err:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
	if err!=nil{
		panic(err)
	}

	var idr[] interface {}
	x:=strings.Split(r.RequestURI,"=")
	var Name string
	if len(x)==1{
		Name=""
	}else{
		Name=x[1]
	}


	query:=`select * from customer inner join address on customer.id=address.cid`
	if len(Name)!=0{
		query+=` where customer.name=?`
		idr=append(idr,Name)
	}
	rows,err:=db.Query(query,idr...)
	if err!=nil {
		panic(err)
	}

	var result []Customer

	defer rows.Close()

	for rows.Next() {
		var cust Customer

		err = rows.Scan(&cust.Id,&cust.Name,&cust.Dob,&cust.Add.Id,&cust.Add.StreetName,&cust.Add.City,&cust.Add.State,&cust.Add.CustomerId)
		result = append(result, cust)


	}
	json.NewEncoder(w).Encode(result)

}


func getCustomerId(w http.ResponseWriter,r *http.Request){
	db,err:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
	if err!=nil{
	}

	defer db.Close()
	vars:=mux.Vars(r)
	id:=vars["id"]


	query:="select * from cust"
	var ids[] interface{}
	if(id!="0"){
		query=`select * from customer inner join address on customer.id=address.cid where customer.id=?`
		x,_:=strconv.Atoi(id)
		ids=append(ids,x)
	}
	rows,_:=db.Query(query,ids...)
	defer rows.Close()
	var c Customer

	for rows.Next(){

		rows.Scan(&c.Id,&c.Name,&c.Dob,&c.Add.Id,&c.Add.StreetName,&c.Add.City,&c.Add.State,&c.Add.CustomerId)

	}

		json.NewEncoder(w).Encode(c)
}


func createCustomer(w http.ResponseWriter,r *http.Request){
	db,err:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
		if err!=nil{
		}
	var idr[] interface{}
	var c Customer
	query:=`insert into customer (name,dob) values(?,?)`
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&c)
	idr=append(idr,&c.Name)
	idr=append(idr,&c.Dob)
	age:=dateInSeconds(c.Dob)
	if age/(365*24*3600) <18 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w,"not eligible")
		return
	}
	row,er:=db.Exec(query,idr...)
	if er!=nil{

	}
	query=`insert into address (street_name,city,state,cid) values(?,?,?,?)`
	var idd[] interface{}
	idd=append(idd,&c.Add.StreetName)
	idd=append(idd,&c.Add.City)
	idd=append(idd,&c.Add.State)

	id,err:=row.LastInsertId()
	idd=append(idd,id)
	_,err=db.Exec(query,idd...)
	if err!=nil{

	}
	query=`select * from customer inner join address on customer.id=address.cid where customer.id=?`


	rows,err:=db.Query(query,id)
	var cust Customer
	for rows.Next() {

		rows.Scan(&cust.Id, &cust.Name, &cust.Dob, &cust.Add.Id, &cust.Add.StreetName, &cust.Add.City, &cust.Add.State, &cust.Add.CustomerId)
	}
	json.NewEncoder(w).Encode(cust)
}



func updateCustomer(w http.ResponseWriter,r *http.Request){
		db,err:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
		if err!=nil{

		}
		vars:=mux.Vars(r)
		ide:=vars["id"]
		query:=`update customer set`
		var idr [] interface{}

		var c Customer
		body,_:=ioutil.ReadAll(r.Body)
		json.Unmarshal(body,&c)
		fmt.Println(c)
		if c.Name!=""{
			query+=" customer.name=?"
			idr=append(idr,c.Name)
		}

		query+=" where customer.id=?"

		idr=append(idr,ide)
		_,er:=db.Exec(query,idr...)

		if er!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Println("c addrs is ", c.Add)
		query=`update address set `
		var idd []interface{}
		if c.Add.StreetName!="" {
			idd=append(idd,c.Add.StreetName)
			query+="street_name=?, "
		}

		if c.Add.City!=""{
			idd=append(idd,c.Add.City)
			query+="city=?, "
		}

		if c.Add.State!=""{
			idd=append(idd,c.Add.State)
			query+="state=? "
		}

		query+="where address.cid=?"
		idd=append(idd,ide)
		_,err1:=db.Exec(query,idd...)

		fmt.Println("in ")
		if err1!=nil{

		}




		query=`select * from customer inner join address on customer.id=address.cid where customer.id=?`
		rows,_:=db.Query(query,ide)
	var cust Customer
		for rows.Next(){

			rows.Scan(&cust.Id,&cust.Name,&cust.Dob,&cust.Add.Id,&cust.Add.StreetName,&cust.Add.City,&cust.Add.State,&cust.Add.CustomerId)

		}
		json.NewEncoder(w).Encode(cust)

}
func DeleteCutomer(w http.ResponseWriter,r *http.Request){
		db,err:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
		defer db.Close()
		if err!=nil{

		}

		vars:=mux.Vars(r)
	    id,_:=strconv.Atoi(vars["id"])
	    var idr[] interface{}
	    idr=append(idr,id)
		var c Customer
		query:=`select * from customer inner join address on customer.id=address.cid where customer.id=?`
		rows,err:=db.Query(query,idr...)

		switch{
		case err!=nil:
			log.Fatal(err)
		default:
			for rows.Next() {
				rows.Scan(&c.Id, &c.Name, &c.Dob, &c.Add.Id,&c.Add.StreetName,&c.Add.City,&c.Add.State,&c.Add.CustomerId)

			}
			query = `delete from customer where id=?`
			_, err = db.Exec(query, idr...)

			json.NewEncoder(w).Encode(c)

		}

}


func main(){
	r:=mux.NewRouter()
	r.HandleFunc("/customer",getCustomer).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id:[0-9]+}",getCustomerId).Methods(http.MethodGet)
	r.HandleFunc("/customer/",createCustomer).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id:[0-9]+}",updateCustomer).Methods(http.MethodPut)
	r.HandleFunc("/customer/{id:[0-9]+}",DeleteCutomer).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":3003",r))
}
