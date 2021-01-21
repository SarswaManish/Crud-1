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
	Slice := strings.Split(d1, "/")

	newDate := Slice[2] + "/" + Slice[1] + "/" + Slice[0]
	myDate, ok := time.Parse("2006/01/02", newDate)

	if ok != nil {
		panic(ok)
	}

	return int(time.Now().Unix() - myDate.Unix())
}


func getCustomer(w http.ResponseWriter,r *http.Request){

	db,ok:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
	defer db.Close()
	if ok!=nil{
		panic(ok)
	}

	var info[] interface {}
	slice:=strings.Split(r.RequestURI,"=")
	var Name string
	if len(slice)==1{
		Name=""
	}else{
		Name=slice[1]
	}


	query:=`select * from customer inner join address on customer.id=address.cid`
	if len(Name)!=0{
		query+=` where customer.name=?`
		info=append(info,Name)
	}
	rows,ok:=db.Query(query,info...)
	if ok!=nil {
		panic(ok)
	}

	var response []Customer

	defer rows.Close()

	for rows.Next() {
		var detail Customer
		ok = rows.Scan(&detail.Id,&detail.Name,&detail.Dob,&detail.Add.Id,&detail.Add.StreetName,&detail.Add.City,&detail.Add.State,&detail.Add.CustomerId)
		response = append(response, detail)
	}
	if response==nil{
		w.WriteHeader(http.StatusNoContent)
		return
	}
	json.NewEncoder(w).Encode(response)


}


func getCustomerById(w http.ResponseWriter,r *http.Request){
	db,ok:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
	if ok!=nil{
		panic(ok)
	}

	defer db.Close()
	vars:=mux.Vars(r)
	id:=vars["id"]

	query:="select * from customer inner join address on customer.id=address.cid"
	if id=="0"{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var info[] interface{}
	if id!="0" {
		query=`select * from customer inner join address on customer.id=address.cid where customer.id=?`
		x,_:=strconv.Atoi(id)
		info=append(info,x)
	}
	row,_:=db.Query(query,info...)
	defer row.Close()
	var response Customer
	for row.Next(){
		row.Scan(&response.Id,&response.Name,&response.Dob,&response.Add.Id,&response.Add.StreetName,&response.Add.City,&response.Add.State,&response.Add.CustomerId)
	}
		if response.Id==0{
			w.WriteHeader(http.StatusNoContent)
			return
		}
		json.NewEncoder(w).Encode(response)
}


func createCustomer(w http.ResponseWriter,r *http.Request){
	db,ok:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
	defer db.Close()
	if ok!=nil{
			panic(ok)
		}
	var info[] interface{}
	var customer Customer
	query:=`insert into customer (name,dob) values(?,?)`
	body,_:=ioutil.ReadAll(r.Body)
	json.Unmarshal(body,&customer)
	if customer.Name=="" || customer.Dob==""{
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	info=append(info,&customer.Name)
	info=append(info,&customer.Dob)
	age:=dateInSeconds(customer.Dob)
	if age/(365*24*3600) <18 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w,"not eligible")
		return
	}
	row,_:=db.Exec(query,info...)
	query=`insert into address (street_name,city,state,cid) values(?,?,?,?)`
	var addr[] interface{}
	if customer.Add.StreetName=="" || customer.Add.City=="" || customer.Add.State==""{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	addr=append(addr,&customer.Add.StreetName)
	addr=append(addr,&customer.Add.City)
	addr=append(addr,&customer.Add.State)

	id,ok1:=row.LastInsertId()
	if ok1!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	addr=append(addr,id)
	_,ok=db.Exec(query,addr...)
	if ok!=nil{
		panic(ok)
	}
	query=`select * from customer inner join address on customer.id=address.cid where customer.id=?`


	newRow,_:=db.Query(query,id)
	var detail Customer
	for newRow.Next() {
		newRow.Scan(&detail.Id, &detail.Name, &detail.Dob, &detail.Add.Id, &detail.Add.StreetName, &detail.Add.City, &detail.Add.State, &detail.Add.CustomerId)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(detail)
}


//update customer will update the details of customer other than id and dateOfBirth(dob)
func updateCustomer(w http.ResponseWriter,r *http.Request){
		db,ok:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
		defer db.Close()
		if ok!=nil{
			panic(ok)
		}
		vars:=mux.Vars(r)
		id:=vars["id"]
		query:=`update customer set`
		var info [] interface{}

		var customer Customer
		body,_:=ioutil.ReadAll(r.Body)
		json.Unmarshal(body,&customer)
		if customer.Name!=""{
			query+=" name=?"
			info=append(info,customer.Name)
		}

		query+=" where customer.id=?"
		info=append(info,id)
		_,er:=db.Exec(query,info...)

		if er!=nil{
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query=`update address set `
		var idd []interface{}
		if customer.Add.StreetName!="" {
			idd=append(idd,customer.Add.StreetName)
			query+=" street_name=? "
		}

		if customer.Add.City!=""{
			idd=append(idd,customer.Add.City)
			query+=", city=? "
		}

		if customer.Add.State!=""{
			idd=append(idd,customer.Add.State)
			query+=", state=? "
		}

		query+="where address.cid=?"
		idd=append(idd,id)
		_,ok1:=db.Exec(query,idd...)

		if ok1!=nil{
			panic(ok1)
		}

		query=`select * from customer inner join address on customer.id=address.cid where customer.id=?`
		rows,_:=db.Query(query,id)
		var detail Customer
		for rows.Next(){
			rows.Scan(&detail.Id,&detail.Name,&detail.Dob,&detail.Add.Id,&detail.Add.StreetName,&detail.Add.City,&detail.Add.State,&detail.Add.CustomerId)
		}
		if detail.Id==0{
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(detail)

}
func DeleteCustomer(w http.ResponseWriter,r *http.Request){
		db,ok:=sql.Open("mysql","root:Manish@123Sharma@/Customer_services")
		defer db.Close()
		if ok!=nil{
			w.WriteHeader(http.StatusBadRequest)
		}
		vars:=mux.Vars(r)
	    id,_:=strconv.Atoi(vars["id"])
	    var info[] interface{}
	    info=append(info,id)

	    query := `delete from customer where id=?`
	    _, ok = db.Exec(query, info...)
	    if ok!=nil{
	    	w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusNoContent)
	    fmt.Println("successfully deleted")
}


func main(){
	r:=mux.NewRouter()
	r.HandleFunc("/customer",getCustomer).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id:[0-9]+}",getCustomerById).Methods(http.MethodGet)
	r.HandleFunc("/customer/",createCustomer).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id:[0-9]+}",updateCustomer).Methods(http.MethodPut)
	r.HandleFunc("/customer/{id:[0-9]+}",DeleteCustomer).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":3003",r))
}
