package store_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"icloud-api/internal/domain"
)

func TestAppleDirectoryObservationPreservesInventory(t *testing.T) {
	ctx := context.Background()
	db := openTestStore(t)
	account := createAccount(t, ctx, db, "directory counts", "directory-counts@icloud.com")
	var before []domain.Alias
	for _, address := range []string{"present@icloud.com", "missing@icloud.com"} {
		a, err := db.CreateAlias(ctx, domain.Alias{AccountID: account.ID, Address: address, APIKeyHash: []byte(address), APIKeyPrefix: "fixture", Enabled: true})
		if err != nil {
			t.Fatal(err)
		}
		before = append(before, a)
	}
	observed := time.Now().UTC()
	if err := db.RecordAppleDirectoryObservation(ctx, account, "icloud_web", observed, []domain.AliasImportCandidate{{Address: before[0].Address, Active: true}}); err != nil {
		t.Fatal(err)
	}
	got, err := db.GetAccount(ctx, account.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.AliasCount != 2 || got.AppleDirectory == nil || got.AppleDirectory.Total != 1 || got.AppleDirectory.LocalMissing != 1 || got.AppleDirectory.Active != 1 {
		t.Fatalf("snapshot = %+v", got.AppleDirectory)
	}
	for _, a := range before {
		after, err := db.GetAlias(ctx, a.ID)
		if err != nil || !reflect.DeepEqual(a, after) {
			t.Fatalf("directory changed alias %d: %v", a.ID, err)
		}
	}
	if err := db.RecordAppleDirectoryObservation(ctx, account, "apple_account", observed.Add(-time.Second), nil); err != nil {
		t.Fatal(err)
	}
	got, err = db.GetAccount(ctx, account.ID)
	if err != nil || got.AppleDirectory.Total != 1 {
		t.Fatal("older observation replaced newer snapshot")
	}
	wrongIdentity := account
	wrongIdentity.Email = "different@icloud.com"
	if err := db.RecordAppleDirectoryObservation(ctx, wrongIdentity, "icloud_web", observed.Add(time.Second), nil); err == nil {
		t.Fatal("wrong identity accepted")
	}
	if err := db.RecordAppleDirectoryObservation(ctx, account, "icloud_web", observed.Add(time.Second), nil); err != nil {
		t.Fatal(err)
	}
	got, err = db.GetAccount(ctx, account.ID)
	if err != nil || got.AppleDirectory.Total != 0 || got.AppleDirectory.LocalMissing != 2 || got.AliasCount != 2 {
		t.Fatal("explicit empty directory corrupted local counts")
	}
	if err := db.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	got, err = db.GetAccount(ctx, account.ID)
	if err != nil || got.AppleDirectory.Total != 0 {
		t.Fatal("repeat migration lost observation")
	}
	changed := account
	changed.IMAPUsername = "different-login@icloud.com"
	if _, err := db.UpdateAccount(ctx, changed); err != nil {
		t.Fatal(err)
	}
	got, err = db.GetAccount(ctx, account.ID)
	if err != nil || got.AppleDirectory != nil {
		t.Fatal("old identity snapshot exposed after account update")
	}
	if err := db.RecordAppleDirectoryObservation(ctx, account, "icloud_web", observed.Add(2*time.Second), nil); err == nil {
		t.Fatal("stale identity observation accepted")
	}
}
