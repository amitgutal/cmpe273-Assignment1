// Server
package main

import 
(
 	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	"strings"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"log"
	"encoding/json"
	"strconv"	
)

type Trade struct { 

  Stocks [5]string
  Quantity [5]int
  Price [5]float64
  TradeID int
  
  Remaining_budget float64
}

type Client struct { 
 
   Stocks []string 
   Percentage []float64
 //  Price []float64
 //  Quantity []int
   Budget float64 
}

type TradeDetails struct { 
  
   Trade_Number int
}

type Message struct {
		
		name, price string
	
	}
	
type Reply struct {
	
   TradeId int
   Stocks [5]string
   Price [5]float64
   Quantity [5]int
   UninvestedAmount float64
   CurrentMarketValue float64
   Profit_Loss [5]string
 }

type Arith struct { }  

func Display(str []string) { 

  fmt.Println(str)
}

var T []Trade 

var Trade_Counter int
var Trade_Id int

func getQuotes(c chan string) { 
 
     for i:=0; ; i++ { 
	
	 // fmt.Println("Inside getQuotes")
     // http://finance.yahoo.com/webservice/v1/symbols/GOOG/quote?format=json
	// http://finance.yahoo.com/d/quotes.csv?s=AAPL+GOOG+MSFT+YHOO&f=sl1
	
  //    resp, err := http.Get("http://finance.yahoo.com/webservice/v1/symbols/GOOG/quote?format=json")
      
    //   robots, err1 := ioutil.ReadAll(resp.Body)
	 //  const jsonStream = robots 
	
	  
     // str_quotes := GetStockDetails("GOOG") + GetStockDetails("YHOO") + GetStockDetails("MSFT") + GetStockDetails("AAPL")
	  
	//  str_quotes := strings.Join(GetStockDetails("GOOG")," ")
	
	  str_quotes := fmt.Sprint(GetStockDetails("GOOG"),GetStockDetails("YHOO"),GetStockDetails("MSFT"),GetStockDetails("AAPL"))
      
	 // str_quotes := strings.Join(str_quotes)
	
	  // quotes := GetStockDetails("MSFT")
	
	  // fmt.Println(str_quotes)
	 
	 /* for j:=0;j<=len(str_quotes);j++ { 
	         fmt.Println(str_quotes[i])
	      } */
	  // fmt.Println(GetStockDetails("MSFT"))
	  // fmt.Println(str_quotes)
	
// *************** JSON Parsing 	  
	  
// ***************** JSON Parsing

      // time.Sleep(time.Second * 2)
	     
	  c <- str_quotes
	 
   }
} 

func GetStockDetails(stock_code string) []string { 

type meta struct  { 
    Type string
    Start int 
    Count int
}

type fields struct { 

                               Name string
                               Price string
                               Symbol string
                               Ts string
                               Type string
                               Utctime string
                               Volume string

                   }
				   
type resource struct { 

           Classname string
           Fields fields
} 

type resources []struct { 
     Resource resource

} 

type list struct { 

   Meta meta 
   Resources resources
}	

type Parse struct { 

   List list
}

// var stock_details []string
Stock_details :=make([]string,3)

var parsed Parse
	
     // fmt.Println("Inside getQuotes")
     // http://finance.yahoo.com/webservice/v1/symbols/GOOG/quote?format=json
	// http://finance.yahoo.com/d/quotes.csv?s=AAPL+GOOG+MSFT+YHOO&f=sl1
	
      url := "http://finance.yahoo.com/webservice/v1/symbols/"+stock_code+"/quote?format=json"
	  	 
      resp, err := http.Get(url)
      
      robots, err1 := ioutil.ReadAll(resp.Body)
	
     // var s MsgSlice
	if err1 != nil || err != nil { 
	
	  fmt.Println(err1)
	  fmt.Println(err)
	}
/*  json.Unmarshal(robots,&s)
	  fmt.Println(s) */
	
	 //  const jsonStream = robots 
	 
	   
	//  str_quotes := string(robots[:])
	
// ******//********* JSON Parsing 	

err2 := json.Unmarshal(robots, &parsed)

        if err2 != nil { 
                
				fmt.Println(err2)  
		}
	        
		Stock_details[0] =""
		Stock_details[1] =""
        Stock_details[2] ="" 
		
		Stock_details[0] = parsed.List.Resources[0].Resource.Fields.Name
		Stock_details[1] = parsed.List.Resources[0].Resource.Fields.Symbol
		Stock_details[2] = parsed.List.Resources[0].Resource.Fields.Price
		resp.Body.Close()
		
		// fmt.Println(Stock_details[:])
	
	    // details := strings.Join(Stock_details," ")
	 
	    return Stock_details  
}

func GetPrice(stock string) string  { 
       Stock_details := GetStockDetails(stock) 
	   
	   return Stock_details[2]   
  }

func GetStockBudget(TotalBudget float64, StockPercentage float64) float64 { 

    Stock_budget :=  (TotalBudget * StockPercentage) / 100
	
    return Stock_budget
}

func (t *Arith) GetTradeDetails(args *TradeDetails, reply *Reply) error { 

      // args.Trade_Number 
	   
	   // int T_Counter =0
	   
	   for counter:=0;counter< len(T);counter++ { 
	       
		   if T[counter].TradeID == args.Trade_Number  {
		    
			   
		        
				for i:=0; i < len(T[counter].Stocks) ;i++ {
				 
			    reply.Profit_Loss[i] ="" 
				if T[counter].Stocks[i] == "" { 
				     break
				}
				reply.TradeId = T[counter].TradeID
				reply.Stocks[i] = T[counter].Stocks[i]
				reply.Quantity[i] = T[counter].Quantity[i]
				
				Stock_price,err  :=  strconv.ParseFloat(GetPrice(T[counter].Stocks[i]),64)
			    	
				if Stock_price > T[counter].Price[i] { 
				
				    reply.Profit_Loss[i] = "+"
				
				} else if Stock_price < T[counter].Price[i] { 
				   
				    reply.Profit_Loss[i] = "-"    
					
				}
				 
				reply.Price[i] = Stock_price
				
				if err!=nil { fmt.Println(err)}
				reply.CurrentMarketValue = reply.CurrentMarketValue + (float64(T[counter].Quantity[i]) * Stock_price)
			  } 
		      reply.UninvestedAmount = T[counter].Remaining_budget
		   }
		   
	   }
	
   return nil
}
func (t *Arith) ExecuteTrade(args *Client, reply *Reply) error {
	
    // fmt.Printf("\n Respective Percentage :  %f\n", args.Budget)
	//  fmt.Println(args)
	// var Temp_Quantity float64
	
	var Rounded_Quantity int
	executed := false

	Tmp_budget :=0.00
	for i:=0 ; i<len(args.Stocks)-1;i++ { 
	     	 
	     // fmt.Println(GetStockDetails(args.Stocks[i]))
		  
		Stock_price,err  :=  strconv.ParseFloat(GetPrice(args.Stocks[i]),64)
		 
		if err!=nil { fmt.Println(err)  }
		 
		Stock_budget := GetStockBudget(args.Budget,args.Percentage[i])
		
		// fmt.Println("Price        --->",Stock_price)
		
		// fmt.Println("Stock Budget --->",Stock_budget)
		
		Stock_quantity := Stock_budget/Stock_price
		
		// fmt.Println(" Stock Quantity : ",Stock_quantity)	 
		
		if Stock_quantity < 1 {  
		
			  fmt.Println("Less Budget for buying ",args.Stocks[i])
			  
	   
		   }  else  { 
		          			   
				   executed = true
				   Rounded_Quantity = int(Stock_quantity)
				   Temp_Qty := float64(Rounded_Quantity)
				   Tmp_budget = Tmp_budget + (Temp_Qty * Stock_price)
				   
				   Remaining_budget := args.Budget-Tmp_budget
			       
				   //Remaining_budget = Remaining_budget + Tmp_budget
				// fmt.Println(" Rounded Quantity %d",Rounded_Quantity," Budget : ",Remaining_budget)
				   
				   T[Trade_Counter].Stocks[i] = args.Stocks[i]
				   reply.Stocks[i] = args.Stocks[i]
				
				   T[Trade_Counter].Quantity[i] = Rounded_Quantity
				   reply.Quantity[i] = Rounded_Quantity
				
				   T[Trade_Counter].Price[i] = Stock_price
				   reply.Price[i] = Stock_price
				
				   T[Trade_Counter].TradeID = Trade_Id
				   reply.TradeId = Trade_Id
				   
				   T[Trade_Counter].Remaining_budget = Remaining_budget
				
				   reply.UninvestedAmount = Remaining_budget
					
			//	   T[Trade_Counter].Remaining_budget = Remaining_budget
			//	   fmt.Println(args.Stocks[i],"\t",Stock_price,"\t",Rounded_Quantity,"\t",Trade_Id)	   
	         //      fmt.Println("Remaining Budget ",Remaining_budget)
				} 
									
    }
    if executed == true { 
	     		 
            Trade_Counter++
			Trade_Id++
	    
		}  else {  
	
	            }			
					
   	return nil
}

func displayQuotes(c chan string) { 
   
for 
    { 
	  
	  ticker_info :=  <- c
     // fmt.Println("Inside Display Quotes") 
	  
	   // stocks := split(ticker_info,"\n")
	//    stocks := ticker_info.split()
	   
	  var stocks []string = strings.Split(ticker_info,"\n")
	  	  
	  fmt.Printf("%s", stocks)
	
	  fmt.Println()	  
	  time.Sleep(time.Second * 3)
	} 
} 

func startServer() {
    arith := new(Arith)

    server := rpc.NewServer()
	
	server.Register(arith)
	rpc.HandleHTTP()

    l, e := net.Listen("tcp", ":8222")
    if e != nil {
        log.Fatal("listen error:", e)
    }

    for {
        conn, err := l.Accept()
		
        if err != nil {
            log.Fatal(err)
        }

        go server.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}

func main()  {

T=make([]Trade,10)

Trade_Id = 1000

Trade_Counter = 0

go startServer()

//	var user *client 
//	user = new(client)
	
var stream chan string = make(chan string)

  fmt.Println(" ***************    YAHOO FINANCE API TICKER  	 *************")
  fmt.Println()
 
  go getQuotes(stream)
  go displayQuotes(stream)

 var input string

 fmt.Scanln(&input)

} 
