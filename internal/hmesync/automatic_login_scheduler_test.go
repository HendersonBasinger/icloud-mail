package hmesync

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"icloud-api/internal/apple"
	"icloud-api/internal/autologincode"
	"icloud-api/internal/domain"
	"icloud-api/internal/secure"
	"icloud-api/internal/store"
)

type automaticSchedulerRepository struct {
	*store.Store
	lists chan struct{}
}

func (r *automaticSchedulerRepository) ListEnabledAccountAutoLoginConfigs(ctx context.Context) ([]domain.AccountAutoLoginConfig, error) {
	configs, err := r.Store.ListEnabledAccountAutoLoginConfigs(ctx)
	select {
	case r.lists <- struct{}{}:
	default:
	}
	return configs, err
}

type automaticSchedulerCodeFetcher struct {
	fetch func(context.Context, string) (autologincode.Result, error)
}

func (f automaticSchedulerCodeFetcher) Fetch(ctx context.Context, url string) (autologincode.Result, error) {
	return f.fetch(ctx, url)
}

func automaticSchedulerAccount(t *testing.T, db *store.Store, index int, enabled bool) domain.Account {
	t.Helper()
	cipher, err := secure.NewCipher([]byte("01234567890123456789012345678901"))
	if err != nil {
		t.Fatal(err)
	}
	email := fmt.Sprintf("scheduler-%d@icloud.com", index)
	ciphertext, err := cipher.EncryptAccountAutoLogin(email, secure.AccountAutoLoginCredentials{
		Password: "fixture-password", CodeURL: fmt.Sprintf("https://codes.example.test/%d", index),
	})
	if err != nil {
		t.Fatal(err)
	}
	account, err := db.CreateAccount(context.Background(), domain.Account{
		Name: "scheduler fixture", Email: email, IMAPHost: "imap.mail.me.com", IMAPPort: 993,
		IMAPUsername: email, PasswordCiphertext: "fixture-imap", Enabled: true,
		AutoLogin: &domain.AccountAutoLoginConfig{Enabled: enabled, Region: "global", CredentialsCiphertext: ciphertext},
	})
	if err != nil {
		t.Fatal(err)
	}
	return account
}

func automaticSchedulerWait[T any](t *testing.T, events <-chan T, description string) T {
	t.Helper()
	select {
	case event := <-events:
		return event
	case <-time.After(3 * time.Second):
		t.Fatalf("timed out waiting for %s", description)
		var zero T
		return zero
	}
}

func automaticSchedulerStart(t *testing.T, service *Service, repo *automaticSchedulerRepository) chan<- time.Time {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	ticks := make(chan time.Time)
	done := make(chan struct{})
	go func() {
		defer close(done)
		service.runAutomaticLogin(ctx, nil, ticks)
	}()
	t.Cleanup(func() {
		cancel()
		automaticSchedulerWait(t, done, "scheduler shutdown and worker cancellation")
	})
	automaticSchedulerWait(t, repo.lists, "initial scheduler refresh")
	return ticks
}

func automaticSchedulerTick(t *testing.T, ticks chan<- time.Time, repo *automaticSchedulerRepository, now time.Time) {
	t.Helper()
	select {
	case ticks <- now:
	case <-time.After(3 * time.Second):
		t.Fatal("scheduler did not accept the next tick")
	}
	automaticSchedulerWait(t, repo.lists, "scheduler refresh")
}

func TestAutomaticLoginSchedulerDeduplicatesAcrossTicksWithoutBlockingOtherAccounts(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "scheduler.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	automaticSchedulerAccount(t, db, 1, true)
	automaticSchedulerAccount(t, db, 2, true)
	repo := &automaticSchedulerRepository{Store: db, lists: make(chan struct{}, 16)}
	var clock atomic.Int64
	clock.Store(time.Date(2026, 10, 22, 1, 0, 0, 0, time.UTC).UnixNano())
	now := func() time.Time { return time.Unix(0, clock.Load()).UTC() }
	service := newTestService(t, repo, &fakeAppleClient{}, &fakeLocker{}, now)
	started := make(chan string, 8)
	var fetches atomic.Int32
	service.autoLoginCodeFetcher = automaticSchedulerCodeFetcher{fetch: func(ctx context.Context, url string) (autologincode.Result, error) {
		fetches.Add(1)
		started <- url
		<-ctx.Done()
		return autologincode.Result{}, ctx.Err()
	}}
	ticks := automaticSchedulerStart(t, service, repo)
	first := automaticSchedulerWait(t, started, "first account code baseline")
	second := automaticSchedulerWait(t, started, "second account while first remains blocked")
	if first == second {
		t.Fatalf("one account started twice while another was blocked: %q", first)
	}
	// Both reserved deadlines are now in the past. The in-flight account gate
	// must still prevent another worker on each subsequent scheduler tick.
	clock.Add(int64(time.Hour))
	for i := 0; i < 5; i++ {
		automaticSchedulerTick(t, ticks, repo, now())
	}
	if got := fetches.Load(); got != 2 {
		t.Fatalf("cross-tick duplicate worker reached code provider: fetches=%d want=2", got)
	}
}

func TestAutomaticLoginSchedulerDisabledConfigurationAndHealthyChannelsAvoidNetwork(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "inactive.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	disabled := automaticSchedulerAccount(t, db, 1, false)
	healthy := automaticSchedulerAccount(t, db, 2, true)
	repo := &automaticSchedulerRepository{Store: db, lists: make(chan struct{}, 16)}
	now := time.Date(2026, 10, 22, 2, 0, 0, 0, time.UTC)
	var appleCalls, codeCalls atomic.Int32
	client := &accountTestClient{
		accountSignIn: func(session apple.AccountSession) (apple.AccountSession, bool, error) {
			appleCalls.Add(1)
			return session, false, errors.New("unexpected fixture Apple request")
		},
		fakeAppleClient: fakeAppleClient{signIn: func(context.Context, string, string, apple.Region, *apple.Session) (apple.Session, bool, error) {
			appleCalls.Add(1)
			return apple.Session{}, false, errors.New("unexpected fixture Apple request")
		}},
	}
	service := newTestService(t, repo, client, &fakeLocker{}, func() time.Time { return now })
	service.autoLoginCodeFetcher = automaticSchedulerCodeFetcher{fetch: func(context.Context, string) (autologincode.Result, error) {
		codeCalls.Add(1)
		return autologincode.Result{}, autologincode.ErrNoCode
	}}
	managed := apple.AccountSession{AppleID: healthy.Email, Region: apple.RegionGlobal,
		APIKey: "fixture-key", AuthenticatedAt: now, ExpiresAt: now.Add(time.Hour)}
	if err := service.persistAccountManagementSession(context.Background(), healthy.ID, identityOf(healthy), managed); err != nil {
		t.Fatal(err)
	}
	if _, err := service.saveSession(context.Background(), healthy.ID, apple.Session{AppleID: healthy.Email,
		Region: apple.RegionGlobal, SessionToken: "fixture-web-token", Cookies: []apple.PersistentCookie{
			{Name: "X-APPLE-WEBAUTH-TOKEN", Value: "fixture-cookie", Expires: now.Add(time.Hour)},
		}}); err != nil {
		t.Fatal(err)
	}
	if info, err := service.GetAccountSession(context.Background(), healthy.ID); err != nil || info.Status != StatusAuthenticated {
		t.Fatalf("management fixture not healthy: %#v %v", info, err)
	}
	if info, err := service.GetSession(context.Background(), healthy.ID); err != nil || info.Status != StatusAuthenticated {
		t.Fatalf("web fixture not healthy: %#v %v", info, err)
	}
	for _, account := range []domain.Account{disabled, healthy} {
		config, err := db.GetAccountAutoLoginConfig(context.Background(), account.ID)
		if err != nil {
			t.Fatal(err)
		}
		if service.automaticLoginDue(context.Background(), config) {
			t.Fatalf("disabled configuration or healthy channel scheduled a login for account %d", account.ID)
		}
	}
	ticks := automaticSchedulerStart(t, service, repo)
	for i := 0; i < 4; i++ {
		automaticSchedulerTick(t, ticks, repo, now)
	}
	if got := appleCalls.Load(); got != 0 {
		t.Fatalf("disabled/healthy accounts reached Apple: %d calls", got)
	}
	if got := codeCalls.Load(); got != 0 {
		t.Fatalf("disabled/healthy accounts polled verification provider: %d calls", got)
	}
}

func TestAutomaticLoginSchedulerRestartHonorsReservedRetryDeadlines(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "restart.db")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	account := automaticSchedulerAccount(t, db, 1, true)
	config, err := db.GetAccountAutoLoginConfig(ctx, account.ID)
	if err != nil {
		t.Fatal(err)
	}
	var clock atomic.Int64
	clock.Store(time.Date(2026, 10, 22, 3, 0, 0, 0, time.UTC).UnixNano())
	now := func() time.Time { return time.Unix(0, clock.Load()).UTC() }
	service := newTestService(t, db, &fakeAppleClient{}, &fakeLocker{}, now)
	var latestDeadline time.Time
	for _, channel := range []string{domain.AutoLoginChannelAppleAccount, domain.AutoLoginChannelICloudWeb} {
		attempt, err := service.prepareAutomaticLogin(ctx, config, channel)
		if err != nil || attempt == nil {
			t.Fatalf("reserve %s attempt: %v", channel, err)
		}
		state, err := db.GetAccountAutoLoginState(ctx, account.ID, channel)
		if err != nil || state.NextAttemptAt == nil || state.Status != "logging_in" || state.Failures != 1 {
			t.Fatalf("missing durable reservation for %s: %#v %v", channel, state, err)
		}
		if state.NextAttemptAt.After(latestDeadline) {
			latestDeadline = *state.NextAttemptAt
		}
		attempt.credentials = secure.AccountAutoLoginCredentials{}
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	db, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	repo := &automaticSchedulerRepository{Store: db, lists: make(chan struct{}, 16)}
	service = newTestService(t, repo, &fakeAppleClient{}, &fakeLocker{}, now)
	if service.automaticLoginDue(ctx, config) {
		t.Fatal("reopened scheduler ignored an in-progress attempt's durable retry deadline")
	}
	started := make(chan struct{}, 8)
	var fetches atomic.Int32
	service.autoLoginCodeFetcher = automaticSchedulerCodeFetcher{fetch: func(ctx context.Context, _ string) (autologincode.Result, error) {
		fetches.Add(1)
		started <- struct{}{}
		<-ctx.Done()
		return autologincode.Result{}, ctx.Err()
	}}
	ticks := automaticSchedulerStart(t, service, repo)
	for i := 0; i < 3; i++ {
		automaticSchedulerTick(t, ticks, repo, now())
	}
	if got := fetches.Load(); got != 0 {
		t.Fatalf("restart ignored reserved deadline: %d verification-provider requests", got)
	}
	clock.Store(latestDeadline.Add(time.Second).UnixNano())
	automaticSchedulerTick(t, ticks, repo, now())
	automaticSchedulerWait(t, started, "retry after the durable deadline")
	if got := fetches.Load(); got != 1 {
		t.Fatalf("deadline expiry started %d attempts, want 1", got)
	}
}
