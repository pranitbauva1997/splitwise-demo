package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/pranitbauva1997/splitwise-demo/pkg/store"
	"io"
	"net/http"
	"strconv"
)

// UnsettledBalances is map[UserId]Amount
type UnsettledBalances map[int64]int

func CalculateUnsettledBalances(transactions []store.UserTransaction, userId int64) UnsettledBalances {
	balances := make(UnsettledBalances)
	for _, t := range transactions {
		if t.OwedTo == userId {
			outstandingBalance, ok := balances[t.Owes]
			if !ok {
				balances[t.Owes] = t.Amount
			} else {
				balances[t.Owes] = outstandingBalance + t.Amount
			}
		} else if t.Owes == userId {
			outstandingBalance, ok := balances[t.OwedTo]
			if !ok {
				balances[t.OwedTo] = -t.Amount
			} else {
				balances[t.OwedTo] = outstandingBalance - t.Amount
			}
		}
	}

	return balances
}

func GetAllUserIds(b store.AddBillRequest) []int64 {
	userIds := make([]int64, 0)
	userIds = append(userIds, b.CreatedBy)
	for _, t := range b.Transactions {
		userIds = append(userIds, t.OwedTo, t.Owes)
	}

	return userIds
}

func CheckUsers(userIds []int64, app *Application) error {
	for _, v := range userIds {
		valid, err := app.StorageClient.IsUserValid(v)
		if err != nil || !valid {
			return fmt.Errorf("couldn't find the user with user_id=%d: %s", v, err)
		}
	}

	return nil
}

func addBillPost(w http.ResponseWriter, r *http.Request, app *Application) {
	defer r.Body.Close()
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	var newBill store.AddBillRequest
	err = json.Unmarshal(buf, &newBill)
	if err != nil {
		app.Log.err.Printf("couldn't parse the request body: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	valid, err := newBill.Validate()
	if err != nil || !valid {
		app.Log.err.Printf("couldn't validate the bill: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if err = CheckUsers(GetAllUserIds(newBill), app); err != nil {
		app.Log.err.Printf("couldn't validate the bill: %s", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = app.StorageClient.AddBill(newBill)
	if err != nil {
		app.Log.err.Printf("couldn't insert the bill in db: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type Status struct {
		Message string `json:"message"`
	}
	w.Header().Set(ContentType, ContentType_ApplicationJson)
	buf, _ = json.Marshal(Status{Message: "success"})
	_, _ = w.Write(buf)
}

type UnsettledBalance struct {
	UserId int64 `json:"user_id"`
	Amount int   `json:"amount"`
}

type UnsettledBalancesResponse []UnsettledBalance

func transformToUnsettledBalancesResponse(balances UnsettledBalances) UnsettledBalancesResponse {
	response := make(UnsettledBalancesResponse, 0)
	for k, v := range balances {
		response = append(response, UnsettledBalance{
			UserId: k,
			Amount: v,
		})
	}
	return response
}

func summaryGet(w http.ResponseWriter, r *http.Request, app *Application) {
	userIdRaw := r.URL.Query().Get("user_id")
	userId, err := strconv.ParseInt(userIdRaw, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	allTransactions, err := app.StorageClient.UserTransactions(userId)
	if err != nil {
		app.Log.err.Printf("couldn't query the DB: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	balances := CalculateUnsettledBalances(allTransactions, userId)
	response := transformToUnsettledBalancesResponse(balances)

	buf, _ := json.Marshal(response)
	w.Header().Set(ContentType, ContentType_ApplicationJson)
	_, _ = w.Write(buf)
}
