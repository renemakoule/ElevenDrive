package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"autoschool/backend/internal/database"
	"autoschool/backend/internal/models"
)

type GuichetRepository struct {
	db *database.DB
}

func NewGuichetRepository(db *database.DB) *GuichetRepository {
	return &GuichetRepository{db: db}
}

// EnsureSchemaCompat met à jour les contraintes SQL héritées (catégories de formules & types de pièces).
func (r *GuichetRepository) EnsureSchemaCompat(ctx context.Context) {
	_, _ = r.db.ExecContext(ctx, `ALTER TABLE formules_permis DROP CONSTRAINT IF EXISTS formules_permis_categorie_permis_check`)
	_, _ = r.db.ExecContext(ctx, `ALTER TABLE formules_permis ALTER COLUMN categorie_permis TYPE VARCHAR(20)`)
	_, _ = r.db.ExecContext(ctx, `ALTER TABLE candidats DROP CONSTRAINT IF EXISTS candidats_type_piece_id_check`)
	_, _ = r.db.ExecContext(ctx, `ALTER TABLE candidats ADD CONSTRAINT candidats_type_piece_id_check CHECK (type_piece_id IN ('CNI', 'PASSEPORT', 'ATTESTATION', 'CONSULAIRE', 'CARTE_CONSULAIRE', 'PERMIS', 'AUTRE'))`)
}

// GetFormules : Récupère le catalogue officiel de l'école
func (r *GuichetRepository) GetFormules(ctx context.Context, autoEcoleID string) ([]models.FormulePermis, error) {
	var res []models.FormulePermis
	err := database.Retry(ctx, func() error {
		return database.Safe(func() error {
			var e error
			res, e = r.getFormulesOnce(ctx, autoEcoleID)
			return e
		})
	})
	return res, err
}

func (r *GuichetRepository) getFormulesOnce(ctx context.Context, autoEcoleID string) ([]models.FormulePermis, error) {
	query := `
		SELECT id, auto_ecole_id, categorie_permis, intitule, description,
		       heures_theoriques, heures_pratiques, prix_total, acompte_minimum,
		       est_actif, created_at, updated_at
		FROM formules_permis
		WHERE auto_ecole_id = $1 AND est_actif = true
		ORDER BY categorie_permis, prix_total
	`
	rows, err := r.db.QueryContext(ctx, query, autoEcoleID)
	if err != nil {
		return nil, fmt.Errorf("erreur récupération formules: %w", err)
	}
	defer rows.Close()

	var formules []models.FormulePermis
	for rows.Next() {
		var f models.FormulePermis
		if err := rows.Scan(
			&f.ID, &f.AutoEcoleID, &f.CategoriePermis, &f.Intitule, &f.Description,
			&f.HeuresTheoriques, &f.HeuresPratiques, &f.PrixTotal, &f.AcompteMinimum,
			&f.EstActif, &f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		formules = append(formules, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return formules, nil
}

// CreerFormule : Enregistre une nouvelle formule et fixe son prix initial (Directeur ou Secrétaire)
func (r *GuichetRepository) CreerFormule(ctx context.Context, autoEcoleID string, input models.CreerFormuleInput) (*models.FormulePermis, error) {
	query := `
		INSERT INTO formules_permis (
			id, auto_ecole_id, categorie_permis, intitule, description,
			heures_theoriques, heures_pratiques, prix_total, acompte_minimum, est_actif
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, true
		)
		RETURNING id, auto_ecole_id, categorie_permis, intitule, description,
		          heures_theoriques, heures_pratiques, prix_total, acompte_minimum,
		          est_actif, created_at, updated_at
	`
	var f models.FormulePermis
	err := r.db.QueryRowContext(
		ctx, query,
		autoEcoleID, input.CategoriePermis, input.Intitule, input.Description,
		input.HeuresTheoriques, input.HeuresPratiques, input.PrixTotal, input.AcompteMinimum,
	).Scan(
		&f.ID, &f.AutoEcoleID, &f.CategoriePermis, &f.Intitule, &f.Description,
		&f.HeuresTheoriques, &f.HeuresPratiques, &f.PrixTotal, &f.AcompteMinimum,
		&f.EstActif, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("erreur création formule: %w", err)
	}
	return &f, nil
}

// ModifierPrixFormule : Modification d'un prix existant (Réservé au Directeur)
func (r *GuichetRepository) ModifierPrixFormule(ctx context.Context, autoEcoleID string, input models.ModifierPrixFormuleInput) (*models.FormulePermis, error) {
	var query string
	var args []interface{}

	if input.AcompteMinimum != nil {
		query = `
			UPDATE formules_permis
			SET prix_total = $1, acompte_minimum = $2, updated_at = NOW()
			WHERE id = $3 AND auto_ecole_id = $4
			RETURNING id, auto_ecole_id, categorie_permis, intitule, description,
			          heures_theoriques, heures_pratiques, prix_total, acompte_minimum,
			          est_actif, created_at, updated_at
		`
		args = []interface{}{input.NouveauPrix, *input.AcompteMinimum, input.FormuleID, autoEcoleID}
	} else {
		query = `
			UPDATE formules_permis
			SET prix_total = $1, updated_at = NOW()
			WHERE id = $2 AND auto_ecole_id = $3
			RETURNING id, auto_ecole_id, categorie_permis, intitule, description,
			          heures_theoriques, heures_pratiques, prix_total, acompte_minimum,
			          est_actif, created_at, updated_at
		`
		args = []interface{}{input.NouveauPrix, input.FormuleID, autoEcoleID}
	}

	var f models.FormulePermis
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&f.ID, &f.AutoEcoleID, &f.CategoriePermis, &f.Intitule, &f.Description,
		&f.HeuresTheoriques, &f.HeuresPratiques, &f.PrixTotal, &f.AcompteMinimum,
		&f.EstActif, &f.CreatedAt, &f.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("formule introuvable ou non autorisée pour cette auto-école")
		}
		return nil, fmt.Errorf("erreur modification prix formule: %w", err)
	}
	return &f, nil
}

// GetGuichetStats : Statistiques temps réel de caisse et élèves
func (r *GuichetRepository) GetGuichetStats(ctx context.Context, autoEcoleID, agenceID string) (*models.GuichetStats, error) {
	stats := &models.GuichetStats{
		ModeConnectivite: "CLOUD_POSTGRES_CONNECTED",
	}

	queryCandidats := `SELECT COUNT(*) FROM candidats WHERE auto_ecole_id = $1`
	if err := r.db.QueryRowContext(ctx, queryCandidats, autoEcoleID).Scan(&stats.TotalCandidats); err != nil {
		stats.TotalCandidats = 0
	}

	queryCaisse := `
		SELECT COALESCE(fond_de_caisse, 0), COALESCE(total_especes_theorique, 0)
		FROM sessions_caisse
		WHERE agence_id = $1 AND statut = 'OUVERTE'
		ORDER BY date_ouverture DESC
		LIMIT 1
	`
	if errSess := r.db.QueryRowContext(ctx, queryCaisse, agenceID).Scan(&stats.FondDeCaisse, &stats.TotalEspecesTheorique); errSess != nil {
		nowUTC := time.Now().UTC()
		startOfDay := time.Date(nowUTC.Year(), nowUTC.Month(), nowUTC.Day(), 0, 0, 0, 0, time.UTC)
		_ = r.db.QueryRowContext(ctx, `
			SELECT COALESCE(SUM(montant), 0)
			FROM paiements
			WHERE auto_ecole_id = $1
			  AND ($2 = '' OR agence_id::text = $2)
			  AND mode_reglement = 'ESPECES'
			  AND COALESCE(est_annule, false) = false
			  AND COALESCE(reference_transaction, '') <> 'REPRISE_HISTORIQUE'
			  AND date_paiement >= $3
		`, autoEcoleID, agenceID, startOfDay).Scan(&stats.TotalEspecesTheorique)
	}

	return stats, nil
}

// SynchroniserLotEvenements : Appelle process_sync_batch dans PostgreSQL ou insère dans journal_sync_events
func (r *GuichetRepository) SynchroniserLotEvenements(
	ctx context.Context,
	autoEcoleID, agenceID string,
	events []models.SyncEventPayload,
) (*models.SyncBatchResponse, error) {
	nowISO := time.Now().UTC().Format(time.RFC3339)
	if len(events) == 0 {
		return &models.SyncBatchResponse{
			Statut:            "SUCCES",
			EvenementsTraites: 0,
			DoublonsIgnores:   0,
			Erreurs:           0,
			SyncedAt:          nowISO,
			Message:           "Aucun événement en attente à synchroniser",
		}, nil
	}

	if agenceID == "" {
		_ = r.db.QueryRowContext(ctx, `SELECT id FROM agences WHERE auto_ecole_id = $1 ORDER BY est_siege DESC LIMIT 1`, autoEcoleID).Scan(&agenceID)
	}

	eventsJSON, err := json.Marshal(events)
	if err != nil {
		return nil, fmt.Errorf("erreur sérialisation lot de synchronisation: %w", err)
	}

	var rawResult []byte
	errRPC := r.db.QueryRowContext(ctx, `SELECT process_sync_batch(NULLIF($1, '')::uuid, $2::jsonb)`, agenceID, string(eventsJSON)).Scan(&rawResult)
	if errRPC == nil && len(rawResult) > 0 {
		var resp models.SyncBatchResponse
		if json.Unmarshal(rawResult, &resp) == nil {
			if resp.SyncedAt == "" {
				resp.SyncedAt = nowISO
			}
			resp.Message = fmt.Sprintf("%d événement(s) synchronisé(s) via process_sync_batch()", resp.EvenementsTraites)
			return &resp, nil
		}
	}

	return r.ReplaySyncEvents(ctx, autoEcoleID, agenceID, events)
}
