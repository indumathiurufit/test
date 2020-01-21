package main
import (
        "fmt"
        "database/sql"
        "strings"
        _ "github.com/go-sql-driver/mysql"
)
func main(){
 // fmt.Println("Connected")
  db,err:=sql.Open("mysql","root:Root@123@tcp(127.0.0.1:3306)/userinfo")
if err !=nil{
    panic(err.Error())
    defer db.Close()
  }
  var value,HR, HRV, Totalsleep, Litesleep, Deepsleep, Walk, Run,Bike,Calorieburnt,UserId,Date string
  fmt.Println("Enter Your UserId")
  e,_:=fmt.Scanln(&UserId)
  fmt.Print(e)
  fmt.Println("Enter Date")
  e10,_:=fmt.Scanln(&Date)
  fmt.Print(e10)
  fmt.Println(" To find the value of")
  fmt.Println("\n MAXHR or MAXHRV or TOTALSLEEP or LITESLEEP or DEEPSLEEP or MAXWALK or MAXRUN or MAXBIKE or MAXCALORIES_BURNT \n")
  fmt.Scanln(&value)
  var value1=strings.ToUpper(value)
  switch value1{

  case "MAXHR":
    rows1,err1:=db.Query("select Max(HR) from HR where date =?",Date)
    if err1 != nil{
                 panic(err1)
               }
               for rows1.Next(){
                 err1=rows1.Scan(&HR)
                 fmt.Println(HR)
               }
               defer db.Close()

   case "MAXHRV":
     rows2,err2:=db.Query("select Max(HRV) from HRV where date =?",Date)
     if err2 != nil{
       panic(err2)
     }
     for rows2.Next(){
       err2=rows2.Scan(&HRV)
       fmt.Println(HRV)
     }
      defer db.Close()

case "TOTALSLEEP":
              rows3,err3:=db.Query("select Total_time from Sleep where date =?",Date)
            if err3 != nil{
                 panic(err3)
               }
             for rows3.Next(){
             err3=rows3.Scan(&Totalsleep)
             fmt.Println(Totalsleep)
           }
             defer db.Close()

 case "LITESLEEP":
              rows4,err4:=db.Query("select Lite_time from Sleep where date =?",Date)
            if err4 != nil{
                 panic(err4)
               }
             for rows4.Next(){
             err4=rows4.Scan(&Litesleep)
             fmt.Println(Litesleep)
           }
             defer db.Close()

case "DEEPSLEEP":
              rows5,err5:=db.Query("select Deep_time from Sleep where date =?",Date)
              if err5 != nil{
                 panic(err5)
               }
             for rows5.Next(){
             err5=rows5.Scan(&Deepsleep)
             fmt.Println(Deepsleep)
           }
             defer db.Close()

case "MAXWALK":
            rows6,err6:=db.Query("select Max(Walk) from step_count where date =?",Date)
            if err6 != nil{
                 panic(err6)
               }
             for rows6.Next(){
             err6=rows6.Scan(&Walk)
             fmt.Println(Walk)
           }
             defer db.Close()

case "MAXRUN":
            rows7,err7:=db.Query("select Max(Run) from step_count where date =?",Date)
            if err7 != nil{
                 panic(err7)
               }
             for rows7.Next(){
             err7=rows7.Scan(&Run)
             fmt.Println(Run)
           }
             defer db.Close()

case "MAXBIKE":
            rows8,err8:=db.Query("select Max(Bike) from step_count where date =?",Date)
            if err8 != nil{
                 panic(err8)
               }
             for rows8.Next(){
             err=rows8.Scan(&Bike)
             fmt.Println(Bike)
           }
             defer db.Close()

case "MAXCALORIES_BURNT":
            rows9,err9:=db.Query("select  Max(Calories_burnt) from step_count where date =?",Date)
            if err9 != nil{
                 panic(err9)
               }
             for rows9.Next(){
             err=rows9.Scan(&Calorieburnt)
             fmt.Println(Calorieburnt)
           }
             defer db.Close()

default :
fmt.Println("Incorrect")
}
}
