package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"icloud-api/internal/domain"
)

func migrateAppleDirectoryObservations(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS apple_directory_observations (
		account_id BIGINT PRIMARY KEY REFERENCES accounts(id) ON DELETE CASCADE,
		source TEXT NOT NULL CHECK(source IN ('icloud_web','apple_account')),
		identity_email TEXT NOT NULL,
		identity_imap_username TEXT NOT NULL,
		identity_mailbox_version BIGINT NOT NULL,
		observed_at BIGINT NOT NULL,
		total BIGINT NOT NULL CHECK(total >= 0),
		active BIGINT NOT NULL CHECK(active >= 0),
		inactive BIGINT NOT NULL CHECK(inactive >= 0),
		local_missing BIGINT NOT NULL CHECK(local_missing >= 0))`)
	return err
}

// RecordAppleDirectoryObservation never disables or deletes missing aliases.
// It publishes counts only after a validated sync and rechecks account identity.
func (s *Store) RecordAppleDirectoryObservation(ctx context.Context, expected domain.Account, source string, startedAt time.Time, candidates []domain.AliasImportCandidate) error {
	accountID := expected.ID
	if accountID < 1 || startedAt.IsZero() || (source != "icloud_web" && source != "apple_account") {
		return fmt.Errorf("invalid Apple directory observation")
	}
	remote := make(map[string]bool, len(candidates))
	active := 0
	for _, c := range candidates {
		address := domain.NormalizeEmail(c.Address)
		if address == "" {
			return fmt.Errorf("empty Apple directory address")
		}
		if _, exists := remote[address]; exists {
			return fmt.Errorf("duplicate Apple directory address")
		}
		remote[address] = c.Active
		if c.Active {
			active++
		}
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := s.lockAccountVersionForUpdate(ctx, tx, accountID); err != nil {
		return err
	}
	state, err := s.readAccountMailboxStateTx(ctx, tx, accountID)
	if err != nil {
		return err
	}
	var configVersion int64
	if err := s.txQueryRowContext(ctx, tx, `SELECT COALESCE((SELECT updated_at FROM account_mailbox_settings WHERE account_id = ?),0)`, accountID).Scan(&configVersion); err != nil {
		return err
	}
	if state.MailboxType != domain.MailboxTypeICloud || domain.NormalizeEmail(state.Email) != domain.NormalizeEmail(expected.Email) ||
		state.IMAPUsername != expected.IMAPUsername || configVersion != expected.MailboxConfigVersion {
		return ErrICloudMailboxRequired
	}
	rows, err := s.txQueryContext(ctx, tx, `SELECT address FROM aliases WHERE account_id = ? AND created_at <= ?`, accountID, timestamp(startedAt))
	if err != nil {
		return err
	}
	missing := 0
	for rows.Next() {
		var address string
		if err := rows.Scan(&address); err != nil {
			rows.Close()
			return err
		}
		if _, exists := remote[domain.NormalizeEmail(address)]; !exists {
			missing++
		}
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return err
	}
	_, err = s.txExecContext(ctx, tx, `INSERT INTO apple_directory_observations(account_id,source,identity_email,identity_imap_username,identity_mailbox_version,observed_at,total,active,inactive,local_missing)
		VALUES (?,?,?,?,?,?,?,?,?,?) ON CONFLICT(account_id) DO UPDATE SET
		identity_email=excluded.identity_email,identity_imap_username=excluded.identity_imap_username,identity_mailbox_version=excluded.identity_mailbox_version,
		source=excluded.source,observed_at=excluded.observed_at,total=excluded.total,
		active=excluded.active,inactive=excluded.inactive,local_missing=excluded.local_missing
		WHERE apple_directory_observations.observed_at <= excluded.observed_at`,
		accountID, source, expected.Email, expected.IMAPUsername, expected.MailboxConfigVersion, timestamp(startedAt), len(remote), active, len(remote)-active, missing)
	if err != nil {
		return err
	}
	return tx.Commit()
}
