package usecase

import (
	"github.com/sirupsen/logrus"
	"github.com/timickb/txserver/internal/domain"
	"github.com/timickb/txserver/internal/infrastructure/storage/dummy"
	"testing"
)

type txTestCase struct {
	senderId   string
	receiverId string
	amount     int64
}

func prepareService(accounts []domain.Account) (*TransactionUseCase, error) {
	repo := dummy.New()
	service := NewTransaction(repo, logrus.New())

	for _, acc := range accounts {
		if err := repo.CreateAccount(acc); err != nil {
			return nil, err
		}
	}

	return service, nil
}

func TestPrimaryTxValidSuccess(t *testing.T) {
	accounts := []domain.Account{
		{Id: "1", Balance: 1000},
		{Id: "2", Balance: 1000},
	}
	service, err := prepareService(accounts)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []txTestCase{
		{
			senderId:   accounts[0].Id,
			receiverId: accounts[1].Id,
			amount:     200,
		},
		{
			senderId:   accounts[1].Id,
			receiverId: accounts[0].Id,
			amount:     200,
		},
	}

	for _, testCase := range testCases {
		err := service.primaryTxValid(testCase.senderId, testCase.receiverId, testCase.amount)
		if err != nil {
			t.Fatal(err, "case = ", testCase)
		}
	}
}

func TestPrimaryTxValidFail(t *testing.T) {
	accounts := []domain.Account{
		{Id: "1", Balance: 1000},
		{Id: "2", Balance: 1000},
	}
	service, err := prepareService(accounts)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []txTestCase{
		// Transaction to user itself
		{
			senderId:   accounts[1].Id,
			receiverId: accounts[1].Id,
			amount:     200,
		},
		// Zero amount
		{
			senderId:   accounts[0].Id,
			receiverId: accounts[1].Id,
			amount:     0,
		},
		// Negative amount
		{
			senderId:   accounts[0].Id,
			receiverId: accounts[1].Id,
			amount:     -1000,
		},
		// Sender not found
		{
			senderId:   "broken",
			receiverId: accounts[1].Id,
			amount:     1,
		},
		// Receiver not found
		{
			senderId:   accounts[0].Id,
			receiverId: "broken",
			amount:     1,
		},
	}

	for _, testCase := range testCases {
		err := service.primaryTxValid(testCase.senderId, testCase.receiverId, testCase.amount)
		if err == nil {
			t.Fatal("expected error on case", testCase)
		}
	}
}
