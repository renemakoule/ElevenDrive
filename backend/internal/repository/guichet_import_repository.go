package repository

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"

	"autoschool/backend/internal/models"
)

// ImporterCandidatsBatch enrôle un lot de candidats dans une transaction atomique
// avec détection des doublons et numérotation séquentielle sans collision.
func (r *GuichetRepository) ImporterCandidatsBatch(
	ctx context.Context,
	autoEcoleID, agenceID, caissierID string,
	input *models.ImporterCandidatsBatchInput,
) (*models.ImporterCandidatsBatchResult, error) {
	result := &models.ImporterCandidatsBatchResult{
		TotalRecus:    len(input.Candidats),
		DossiersCrees: []models.ImportDossierCree{},
		Erreurs:       []models.ImportErreurLigne{},
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Pré-charger les formules actives de l'auto-école
	formulesMap := make(map[string]models.FormulePermis)
	fRows, err := tx.QueryContext(ctx, `
		SELECT id, auto_ecole_id, categorie_permis, intitule, prix_total, acompte_minimum, heures_theoriques, heures_pratiques
		FROM formules_permis WHERE auto_ecole_id = $1
	`, autoEcoleID)
	if err != nil {
		return nil, fmt.Errorf("erreur chargement formules: %w", err)
	}
	for fRows.Next() {
		var f models.FormulePermis
		if scanErr := fRows.Scan(
			&f.ID, &f.AutoEcoleID, &f.CategoriePermis, &f.Intitule, &f.PrixTotal,
			&f.AcompteMinimum, &f.HeuresTheoriques, &f.HeuresPratiques,
		); scanErr == nil {
			formulesMap[f.ID] = f
		}
	}
	fRows.Close()

	// Pré-charger les numéros de pièces existants pour bloquer les doublons
	existingPieces := make(map[string]bool)
	pRows, err := tx.QueryContext(ctx, `
		SELECT LOWER(TRIM(numero_piece_id)) FROM candidats WHERE auto_ecole_id = $1 AND numero_piece_id IS NOT NULL
	`, autoEcoleID)
	if err == nil {
		for pRows.Next() {
			var num string
			if pRows.Scan(&num) == nil && num != "" {
				existingPieces[num] = true
			}
		}
		pRows.Close()
	}

	var countDossiers int
	_ = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM dossiers_inscription WHERE auto_ecole_id = $1`, autoEcoleID).Scan(&countDossiers)

	var countPaiements int
	_ = tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM paiements WHERE auto_ecole_id = $1`, autoEcoleID).Scan(&countPaiements)

	year := time.Now().Year()
	month := int(time.Now().Month())

	for idx, item := range input.Candidats {
		nom := strings.ToUpper(strings.TrimSpace(item.Nom))
		prenoms := strings.TrimSpace(item.Prenoms)
		pieceNorm := strings.ToLower(strings.TrimSpace(item.NumeroPieceID))
		tel := strings.TrimSpace(item.TelephoneWhatsApp)

		if nom == "" || prenoms == "" || tel == "" || pieceNorm == "" || strings.TrimSpace(item.FormuleID) == "" {
			result.TotalErreurs++
			result.Erreurs = append(result.Erreurs, models.ImportErreurLigne{
				Index: idx + 1, Nom: nom, Prenoms: prenoms, NumeroPieceID: item.NumeroPieceID,
				Raison: "Champs obligatoires manquants (nom, prénoms, WhatsApp, pièce ou formule)",
			})
			continue
		}

		f, ok := formulesMap[item.FormuleID]
		if !ok {
			result.TotalErreurs++
			result.Erreurs = append(result.Erreurs, models.ImportErreurLigne{
				Index: idx + 1, Nom: nom, Prenoms: prenoms, NumeroPieceID: item.NumeroPieceID,
				Raison: "Formule de permis introuvable",
			})
			continue
		}

		if existingPieces[pieceNorm] {
			if input.IgnorerDoublons {
				result.TotalIgnores++
				result.Erreurs = append(result.Erreurs, models.ImportErreurLigne{
					Index: idx + 1, Nom: nom, Prenoms: prenoms, NumeroPieceID: item.NumeroPieceID,
					Raison: "Doublon ignoré : numéro de pièce déjà existant",
				})
				continue
			}
			result.TotalErreurs++
			result.Erreurs = append(result.Erreurs, models.ImportErreurLigne{
				Index: idx + 1, Nom: nom, Prenoms: prenoms, NumeroPieceID: item.NumeroPieceID,
				Raison: "Doublon : numéro de pièce déjà existant",
			})
			continue
		}

		if _, spErr := tx.ExecContext(ctx, "SAVEPOINT sp_import_row"); spErr != nil {
			return nil, spErr
		}

		candID := uuid.New().String()
		dosID := uuid.New().String()
		countDossiers++
		numDossier := fmt.Sprintf("DOS-%d-%04d", year, countDossiers)

		medStatut := "VALIDE_SANS_RESERVE"
		if item.VisiteMedicaleStatut != nil && strings.TrimSpace(*item.VisiteMedicaleStatut) != "" {
			medStatut = strings.TrimSpace(*item.VisiteMedicaleStatut)
		}
		genre := strings.ToUpper(strings.TrimSpace(item.Genre))
		if genre != "F" {
			genre = "M"
		}
		typePiece := normalizeTypePieceID(item.TypePieceID)
		lieu := strings.TrimSpace(item.LieuNaissance)
		if lieu == "" {
			lieu = "Abidjan"
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO candidats (
				id, auto_ecole_id, agence_id, nom, prenoms, genre, date_naissance,
				lieu_naissance, telephone_whatsapp, type_piece_id, numero_piece_id,
				visite_medicale_statut, photo_url
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`, candID, autoEcoleID, agenceID, nom, prenoms, genre, item.DateNaissance,
			lieu, tel, typePiece, strings.TrimSpace(item.NumeroPieceID), medStatut, item.PhotoURL)
		if err != nil {
			_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_import_row")
			countDossiers--
			result.TotalErreurs++
			result.Erreurs = append(result.Erreurs, models.ImportErreurLigne{
				Index: idx + 1, Nom: nom, Prenoms: prenoms, NumeroPieceID: item.NumeroPieceID,
				Raison: fmt.Sprintf("Erreur insertion candidat : %v", err),
			})
			continue
		}

		acompte := math.Max(0, item.AcompteInitial)
		if acompte > f.PrixTotal {
			acompte = f.PrixTotal
		}
		soldeRestant := math.Max(0, f.PrixTotal-acompte)
		modeFormation := "PRESENTIEL"
		if item.ModeFormationTheorique != nil && strings.TrimSpace(*item.ModeFormationTheorique) != "" {
			modeFormation = strings.TrimSpace(*item.ModeFormationTheorique)
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO dossiers_inscription (
				id, auto_ecole_id, agence_id, candidat_id, formule_id, numero_dossier,
				montant_total, remise_accordee, solde_restant, statut_dossier, mode_formation_theorique
			) VALUES ($1, $2, $3, $4, $5, $6, $7, 0, $8, 'EN_COURS', $9)
		`, dosID, autoEcoleID, agenceID, candID, f.ID, numDossier, f.PrixTotal, soldeRestant, modeFormation)
		if err != nil {
			_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_import_row")
			countDossiers--
			result.TotalErreurs++
			result.Erreurs = append(result.Erreurs, models.ImportErreurLigne{
				Index: idx + 1, Nom: nom, Prenoms: prenoms, NumeroPieceID: item.NumeroPieceID,
				Raison: fmt.Sprintf("Erreur insertion dossier : %v", err),
			})
			continue
		}

		if acompte > 0 {
			countPaiements++
			payID := uuid.New().String()
			recuNum := fmt.Sprintf("REC-%d%02d-%05d", year, month, countPaiements)
			modeReg := strings.ToUpper(strings.TrimSpace(item.ModeReglement))
			if modeReg == "" {
				modeReg = "ESPECES"
			}
			var refTx *string
			if !input.ImpacterCaisseJour {
				hist := "REPRISE_HISTORIQUE"
				refTx = &hist
			}

			_, err = tx.ExecContext(ctx, `
				INSERT INTO paiements (
					id, auto_ecole_id, agence_id, dossier_id, montant, mode_reglement,
					reference_transaction, reference_recu, solde_apres_paiement, encaisse_par_id
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`, payID, autoEcoleID, agenceID, dosID, acompte, modeReg, refTx, recuNum, soldeRestant, caissierID)
			if err != nil {
				_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT sp_import_row")
				countDossiers--
				countPaiements--
				result.TotalErreurs++
				result.Erreurs = append(result.Erreurs, models.ImportErreurLigne{
					Index: idx + 1, Nom: nom, Prenoms: prenoms, NumeroPieceID: item.NumeroPieceID,
					Raison: fmt.Sprintf("Erreur paiement initial : %v", err),
				})
				continue
			}

			if input.ImpacterCaisseJour {
				recordReceiptAndCashTx(ctx, tx, payID, recuNum, agenceID, modeReg, acompte)
			} else {
				secHash := fmt.Sprintf("SECHASH-%s-%s", payID[:8], recuNum)
				_, _ = tx.ExecContext(ctx, `
					INSERT INTO recus_caisse (id, paiement_id, numero_recu, code_securite_hash)
					VALUES ($1, $2, $3, $4)
				`, uuid.New().String(), payID, recuNum, secHash)
			}
		}

		_, _ = tx.ExecContext(ctx, "RELEASE SAVEPOINT sp_import_row")
		existingPieces[pieceNorm] = true
		result.TotalCrees++
		result.DossiersCrees = append(result.DossiersCrees, models.ImportDossierCree{
			CandidatID:    candID,
			DossierID:     dosID,
			NumeroDossier: numDossier,
			NomComplet:    fmt.Sprintf("%s %s", nom, prenoms),
			SoldeRestant:  soldeRestant,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return result, nil
}

func normalizeTypePieceID(raw string) string {
	val := strings.ToUpper(strings.TrimSpace(raw))
	switch {
	case val == "":
		return "CNI"
	case strings.Contains(val, "CONSUL"):
		return "CONSULAIRE"
	case strings.Contains(val, "PASS"):
		return "PASSEPORT"
	case strings.Contains(val, "ATTEST"):
		return "ATTESTATION"
	case val == "CNI" || val == "PERMIS" || val == "AUTRE":
		return val
	default:
		return "AUTRE"
	}
}
