package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	allWatching []*Watching
	port        string
	updates     string
)

type Watching struct {
	Name    string
	Address string
	Balance string
}

//
// Fetch BTC balance from blockchain.info
func GetBTCBalance(address string) *big.Float {
	balance := big.NewFloat(0)
	url := fmt.Sprintf("https://blockchain.info/q/addressbalance/%v", address)
	response, err := http.Get(url)
	if err != nil {
		return big.NewFloat(0)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return big.NewFloat(0)
		}
		balance.SetString(string(contents))
	}
	balance.Mul(balance, big.NewFloat(0.00000001))
	return balance
}

//
// HTTP response handler for /metrics
func MetricsHttp(w http.ResponseWriter, r *http.Request) {
	var allOut []string
	for _, v := range allWatching {
		allOut = append(allOut, fmt.Sprintf("btc_balance{name=\"%v\",address=\"%v\"} %v", v.Name, v.Address, v.Balance))
	}
	fmt.Fprintln(w, strings.Join(allOut, "\n"))
}

//
// Open the addresses.txt file (name:address)
func OpenAddresses(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		object := strings.Split(scanner.Text(), ":")
		w := &Watching{
			Name:    object[0],
			Address: object[1],
		}
		allWatching = append(allWatching, w)
	}
	return err
}

func main() {
	port = os.Getenv("PORT")
	err := OpenAddresses("addresses.txt")
	if err != nil {
		panic(err)
	}

	// check address balances
	go func() {
		for {
			for _, v := range allWatching {
				v.Balance = GetBTCBalance(v.Address).String()
			}
			time.Sleep(60 * time.Second)
		}
	}()

	fmt.Printf("BTCexporter has started on port %v\n", port)
	http.HandleFunc("/metrics", MetricsHttp)
	panic(http.ListenAndServe("0.0.0.0:"+port, nil))
}
