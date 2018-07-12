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
	prefix      string
	loadSeconds float64
	totalLoaded int64
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
		time.Sleep(15 * time.Second)
		return GetBTCBalance(address)
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
	total := big.NewFloat(0)
	for _, v := range allWatching {
		if v.Balance == "" {
			v.Balance = "0"
		}
		bal := big.NewFloat(0)
		bal.SetString(v.Balance)
		total.Add(total, bal)
		allOut = append(allOut, fmt.Sprintf("%vbtc_balance{name=\"%v\",address=\"%v\"} %v", prefix, v.Name, v.Address, v.Balance))
	}
	allOut = append(allOut, fmt.Sprintf("%vbtc_balance_total %0.8f", prefix, total))
	allOut = append(allOut, fmt.Sprintf("%vbtc_load_seconds %0.2f", prefix, loadSeconds))
	allOut = append(allOut, fmt.Sprintf("%vbtc_loaded_addresses %v", prefix, totalLoaded))
	allOut = append(allOut, fmt.Sprintf("%vbtc_total_addresses %v", prefix, len(allWatching)))
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
	prefix = os.Getenv("PREFIX")
	err := OpenAddresses("addresses.txt")
	if err != nil {
		panic(err)
	}

	fmt.Printf("BTC Exporter started on port %v, http://0.0.0.0:%v/metrics\n", port, port)

	// check address balances
	go func() {
		for {
			totalLoaded = 0
			t1 := time.Now()
			fmt.Printf("Scanning %v addresses\n", len(allWatching))
			for _, v := range allWatching {
				v.Balance = GetBTCBalance(v.Address).String()
				totalLoaded++
			}
			t2 := time.Now()
			loadSeconds = t2.Sub(t1).Seconds()
			fmt.Printf("Completed Scanning %v addresses in %v seconds, sleeping for 60 seconds\n", len(allWatching), loadSeconds)
			time.Sleep(60 * time.Second)
		}
	}()

	fmt.Printf("BTCexporter has started on port %v\n", port)
	http.HandleFunc("/metrics", MetricsHttp)
	panic(http.ListenAndServe("0.0.0.0:"+port, nil))
}
