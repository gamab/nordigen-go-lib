package nordigen

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccountMetadata struct {
	Id            string `json:"id,omitempty"`
	Created       string `json:"created,omitempty"`
	LastAccessed  string `json:"last_accessed,omitempty"`
	Iban          string `json:"iban,omitempty"`
	InstitutionId string `json:"institution_id,omitempty"`
	// There is an issue in the api, the status is still a string
	// like in v1
	Status string `json:"status,omitempty"`
	//Status        []string `json:"status"`
}

type AccountBalanceAmount struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type AccountBalance struct {
	BalanceAmount      AccountBalanceAmount `json:"balanceAmount,omitempty"`
	BalanceType        string               `json:"balanceType,omitempty"`
	ReferenceDate      string               `json:"referenceDate,omitempty"`
	LastChangeDateTime string               `json:"lastChangeDateTime,omitempty"`
}

type AccountBalances struct {
	Balances []AccountBalance `json:"balances,omitempty"`
}

type AccountDetails struct {
	Account struct {
		ResourceId string `json:"resourceId,omitempty"`
		Iban       string `json:"iban,omitempty"`
		Currency   string `json:"currency,omitempty"`
		OwnerName  string `json:"ownerName,omitempty"`
		Product    string `json:"product,omitempty,"`
		Status     string `json:"status,omitempty"`
		Name       string `json:"name,omitempty"`
	} `json:"account"`
}

type AccountTransactions struct {
	Transactions struct {
		Booked []struct {
			TransactionId     string `json:"transactionId,omitempty"`
			EntryReference    string `json:"entryReference,omitempty"`
			BookingDate       string `json:"bookingDate,omitempty"`
			ValueDate         string `json:"valueDate,omitempty"`
			TransactionAmount struct {
				Amount   string `json:"amount,omitempty"`
				Currency string `json:"currency,omitempty"`
			} `json:"transactionAmount,omitempty"`
			CreditorName    string `json:"creditorName,omitempty"`
			CreditorAccount struct {
				Iban string `json:"iban,omitempty"`
			} `json:"creditorAccount,omitempty"`
			UltimateCreditor string `json:"ultimateCreditor,omitempty"`
			DebtorName       string `json:"debtorName,omitempty"`
			DebtorAccount    struct {
				Iban string `json:"iban,omitempty"`
			} `json:"debtorAccount,omitempty"`
			UltimateDebtor                         string   `json:"ultimateDebtor,omitempty"`
			RemittanceInformationUnstructured      string   `json:"remittanceInformationUnstructured,omitempty"`
			RemittanceInformationUnstructuredArray []string `json:"remittanceInformationUnstructuredArray,omitempty"`
			BankTransactionCode                    string   `json:"bankTransactionCode,omitempty"`
		} `json:"booked,omitempty"`
		Pending []struct {
			TransactionAmount struct {
				Amount   string `json:"amount,omitempty"`
				Currency string `json:"currency,omitempty"`
			} `json:"transactionAmount"`
			ValueDate                              string   `json:"valueDate,omitempty"`
			RemittanceInformationUnstructured      string   `json:"remittanceInformationUnstructured,omitempty"`
			RemittanceInformationUnstructuredArray []string `json:"remittanceInformationUnstructuredArray,omitempty"`
		} `json:"pending,omitempty"`
	} `json:"transactions,omitempty"`
}

const accountPath = "accounts"
const balancesPath = "balances"
const detailsPath = "details"
const transactionsPath = "transactions"

func (c *Client) GetAccountMetadata(ctx context.Context, id string) (AccountMetadata, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{accountPath, id, ""}, "/"),
		},
	}

	req = req.WithContext(ctx)

	resp, err := c.c.Do(req)
	if err != nil {
		return AccountMetadata{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountMetadata{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountMetadata{}, &APIError{resp.StatusCode, string(body), err}
	}
	accMtdt := AccountMetadata{}
	err = json.Unmarshal(body, &accMtdt)

	if err != nil {
		return AccountMetadata{}, err
	}

	return accMtdt, nil
}

func (c *Client) GetAccountBalances(ctx context.Context, id string) (AccountBalances, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{accountPath, id, balancesPath, ""}, "/"),
		},
	}

	req = req.WithContext(ctx)
	resp, err := c.c.Do(req)

	if err != nil {
		return AccountBalances{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountBalances{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountBalances{}, &APIError{resp.StatusCode, string(body), err}
	}
	accBlnc := AccountBalances{}
	err = json.Unmarshal(body, &accBlnc)

	if err != nil {
		return AccountBalances{}, err
	}

	return accBlnc, nil
}

func (c *Client) GetAccountDetails(ctx context.Context, id string) (AccountDetails, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path: strings.Join([]string{accountPath, id, detailsPath, ""}, "/"),
		},
	}

	req = req.WithContext(ctx)
	resp, err := c.c.Do(req)

	if err != nil {
		return AccountDetails{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountDetails{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountDetails{}, &APIError{resp.StatusCode, string(body), err}
	}
	accDtl := AccountDetails{}
	err = json.Unmarshal(body, &accDtl)

	if err != nil {
		return AccountDetails{}, err
	}

	return accDtl, nil
}

func (c *Client) GetAccountTransactions(ctx context.Context, id string, from, to *time.Time) (AccountTransactions, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Path:     strings.Join([]string{accountPath, id, transactionsPath, ""}, "/"),
			RawQuery: dateParams(from, to),
		},
	}

	req = req.WithContext(ctx)
	resp, err := c.c.Do(req)

	if err != nil {
		return AccountTransactions{}, err
	}
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return AccountTransactions{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return AccountTransactions{}, &APIError{resp.StatusCode, string(body), err}
	}
	accTxns := AccountTransactions{}
	err = json.Unmarshal(body, &accTxns)

	if err != nil {
		return AccountTransactions{}, err
	}

	return accTxns, nil
}

func dateParams(from, to *time.Time) string {
	params := url.Values{}
	if from != nil {
		params.Add("date_from", from.Format("2006-01-02"))
	}
	if to != nil {
		params.Add("date_to", to.Format("2006-01-02"))
	}
	return params.Encode()
}
