// Client
package main

import
(
    "fmt"	
	"net"
	"net/rpc/jsonrpc"
	"os"
	"strings"
	"strconv"
)

type Reply struct {

    TradeId int
    Stocks []string
	Price []float64
    Quantity []int
    UninvestedAmount float64
	Profit_Loss [5]string
	CurrentMarketValue float64
	
}

type TradeDetails struct { 

    Trade_Number int64	
}

type Trade struct { 

  Stocks[] string
  Quantity[] string 
  
}

type Client struct { 
 
   Stocks[] string 
   Percentage[] float64
  // Price  []float64
  // Quantity []int
   Budget float64
}

func main() {
   
    var user *Client
	var id *TradeDetails
	
	// *************************************** GET DATA from Command Line
	
	// fmt.Println("Usage ",os.Args[1])

// var stock_names []string
//	var percentage []float64

	// var budget float64
	
	OS_Args := strings.Replace(os.Args[1], "%", "", -1)
	var Details []string = strings.Split(OS_Args,",")
    var Values []string
	
	Total_Percentage := 0.0
	
	// fmt.Println(" Length :",len(Details))
	
	Budget ,err := strconv.ParseFloat(Details[len(Details)-1],64)
	
	// fmt.Println("Budget ",Budget)
	
	Percentage  := make([]float64,len(Details))
    Stock_names := make([]string,len(Details))
	
	for i:=0; i<len(Details); i++  {
	 
	 	// fmt.Println(details[i])
		
	       Values = strings.Split(Details[i],":")
		   
		  //fmt.Println(len(values))
		   
		   for j:=0; j<len(Values)-1; j++ { 
		       
		        Stock_names[i] = Values[0]
				
			   //  percentage[j], err = strconv.ParseFloat(values[1],64)
                    
					if s, err := strconv.ParseFloat(Values[1], 64); err == nil {
				        //fmt.Println(err)
		        		Percentage[i] = s
			            Total_Percentage = Total_Percentage + s
					// fmt.Print(Stock_names[j],"\t",Percentage[j])
		     			
	          }
			   // fmt.Println(" Stock Name ",j,"\t",Stock_names[j])
	
			 //  fmt.Print("\t",values[1])    
		  }   
		 
		
	 }
	/* fmt.Println(Stock_names,"\n",Percentage)
	 fmt.Println() */
	
	/* for k:=0;k<len(Stock_names);k++ {
		
		fmt.Println(Stock_names[k],"\t",Percentage[k],"\n")
	} 
	
	/*  ***************** Trade ID and Details  **********/  
	
	if Total_Percentage > 100 { 
	  fmt.Println("Total Percentage greater than 100")
	  return
	}
	
	var reply Reply
	
	conn, err := net.Dial("tcp", "localhost:8222")

    if err != nil {
        panic(err)
    }
    defer conn.Close()

    c := jsonrpc.NewClient(conn)
	
	if len(Details) == 1 { 
	
	 
	
	 T_id ,err_int := strconv.ParseInt(Details[0],10,64)   
	 if err_int!= nil {   fmt.Println(err_int)  }
	  
	id = &TradeDetails{T_id}
	err = c.Call("Arith.GetTradeDetails",id,&reply)
	
	fmt.Println("Trade ID :",reply.TradeId)
	fmt.Print("Stocks :")
	
	for x:=0;x<len(reply.Stocks);x++ {
	   
	   if reply.Stocks[x] == "" { 
	     
		  break 
	    
		}
      
	    fmt.Print("(",reply.Stocks[x])
	    fmt.Print(":",reply.Quantity[x])
		fmt.Print()
	    fmt.Print(":",reply.Profit_Loss[x],"$",reply.Price[x],")")

	}
	
	    fmt.Println()
	    fmt.Println("Current Market Value : ",reply.CurrentMarketValue)
	    fmt.Println("Un-Invested Amount :",reply.UninvestedAmount)
	
	
	} else { 

    
	
	
	// *********************************************************
	
	user = &Client{Stock_names,Percentage,Budget}
	
	/*****************   EXECUTE TRADE DETAILS   *********************/
	
	/* */
	
	err = c.Call("Arith.ExecuteTrade",user,&reply)

	 fmt.Println("Trade ID :",reply.TradeId)
	 fmt.Print("Stocks :")
	
	for x:=0;x<len(reply.Stocks);x++ {
	   
	   if reply.Stocks[x] == "" { 
	     
		  break 
	   
	   }
      
	  fmt.Print("(",reply.Stocks[x])
	  fmt.Print(":",reply.Quantity[x])
	  fmt.Print(":$",reply.Price[x],")")
	
	}
	
	fmt.Println()
	
	fmt.Println("Un-Invested Amount :",reply.UninvestedAmount) 
	
	}	
}
