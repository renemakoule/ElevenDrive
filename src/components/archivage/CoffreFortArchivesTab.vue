<script setup lang="ts">
import { useArchivageStore } from "../../stores/archivageStore";
import type { DossierArchiveItem, MotifArchivage } from "../../types";
import { FolderArchive, Search, Eye, RotateCcw, Trash2 } from "lucide-vue-next";

interface Props {
  isManager: boolean;
}

defineProps<Props>();

const emit = defineEmits<{
  (e: "view-archive", archive: DossierArchiveItem): void;
  (e: "desarchiver", archive: DossierArchiveItem): void;
  (e: "delete-archive", archive: DossierArchiveItem): void;
}>();

const archivageStore = useArchivageStore();

const fcfa = (n: number) => (n || 0).toLocaleString("fr-FR");

const motifLabels: Record<MotifArchivage, string> = {
  PERMIS_OBTENU: "Permis obtenu / Admis",
  DOSSIER_SOLDE: "Formation achevée & Solde 100%",
  ABANDON_DESISTEMENT: "Abandon / Désistement",
  PRESCRIPTION_LEGALE: "Prescription légale (Délai expiré)",
  TRANSFERT_AUTRE_ECOLE: "Transfert vers autre auto-école",
  AUTRE: "Autre motif administratif",
};
</script>

<template>
  <div class="space-y-4">
    <!-- Barre de recherche et filtres d'archives -->
    <div
      class="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-3 bg-card border border-border p-3 rounded-[4px]"
    >
      <div class="relative flex-1">
        <Search
          class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground"
        />
        <input
          v-model="archivageStore.searchQuery"
          type="text"
          placeholder="Rechercher par nom, prénom, N° dossier, boîte physique, permis..."
          class="w-full pl-9 pr-3 py-1.5 text-xs rounded-[4px] border border-border bg-background text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-1 focus:ring-brand"
        />
      </div>

      <div class="flex items-center gap-2 overflow-x-auto pb-1 sm:pb-0">
        <!-- Filtre motif -->
        <select
          v-model="archivageStore.filterMotif"
          class="px-2.5 py-1.5 rounded-[4px] text-xs border border-border bg-background text-foreground focus:outline-none"
        >
          <option value="ALL">Tous les motifs</option>
          <option value="PERMIS_OBTENU">Permis Obtenu / Admis</option>
          <option value="DOSSIER_SOLDE">Formation Achevée</option>
          <option value="ABANDON_DESISTEMENT">Abandon / Désistement</option>
          <option value="PRESCRIPTION_LEGALE">Prescription Légale</option>
          <option value="TRANSFERT_AUTRE_ECOLE">Transfert</option>
        </select>

        <!-- Filtre année -->
        <select
          v-model="archivageStore.filterAnnee"
          class="px-2.5 py-1.5 rounded-[4px] text-xs border border-border bg-background text-foreground focus:outline-none"
        >
          <option value="ALL">Toutes les années</option>
          <option
            v-for="an in archivageStore.anneesDisponibles"
            :key="an"
            :value="an"
          >
            Année {{ an }}
          </option>
        </select>
      </div>
    </div>

    <!-- Liste des archives -->
    <div class="rounded-[4px] border border-border bg-card overflow-hidden">
      <div
        v-if="archivageStore.filteredArchives.length === 0"
        class="p-12 text-center"
      >
        <FolderArchive
          class="h-10 w-10 text-muted-foreground/40 mx-auto mb-3"
        />
        <p class="text-sm font-medium text-foreground">
          Aucun dossier archivé trouvé
        </p>
        <p class="text-xs text-muted-foreground mt-1 max-w-sm mx-auto">
          Vous pouvez transférer les dossiers terminés ou soldés depuis l'onglet
          « Dossiers Éligibles ».
        </p>
      </div>

      <div v-else class="divide-y divide-border">
        <div
          v-for="item in archivageStore.filteredArchives"
          :key="item.id"
          class="p-4 transition hover:bg-muted/30 flex flex-col md:flex-row md:items-center justify-between gap-4"
        >
          <div class="flex items-start gap-3.5 min-w-0">
            <div
              class="h-10 w-10 rounded-[4px] shrink-0 flex items-center justify-center font-bold text-white bg-brand-dark"
            >
              <FolderArchive class="h-5 w-5" />
            </div>

            <div class="min-w-0">
              <div class="flex flex-wrap items-center gap-2 mb-1">
                <h3 class="text-sm font-bold text-foreground">
                  {{ item.candidat_nom }} {{ item.candidat_prenoms }}
                </h3>
                <span
                  class="font-mono text-xs px-2 py-0.5 rounded bg-muted text-brand-dark font-semibold border border-border"
                >
                  N° {{ item.numero_dossier }}
                </span>
                <span
                  class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-medium border"
                  :class="
                    item.motif_archivage === 'PERMIS_OBTENU'
                      ? 'bg-emerald-500/10 text-emerald-700 border-emerald-500/30 dark:text-emerald-300'
                      : item.motif_archivage === 'DOSSIER_SOLDE'
                        ? 'bg-blue-500/10 text-blue-700 border-blue-500/30 dark:text-blue-300'
                        : 'bg-muted text-muted-foreground border-border'
                  "
                >
                  {{
                    motifLabels[item.motif_archivage] || item.motif_archivage
                  }}
                </span>
              </div>

              <p
                class="text-xs text-muted-foreground flex flex-wrap items-center gap-x-3 gap-y-1"
              >
                <span
                  >Permis {{ item.categorie_permis }} ({{
                    item.formule_intitule
                  }})</span
                >
                <span>•</span>
                <span
                  >Total payé :
                  <strong class="text-foreground"
                    >{{ fcfa(item.total_paye) }} FCFA</strong
                  ></span
                >
                <span>•</span>
                <span
                  >Boîte :
                  <code class="font-mono text-foreground">{{
                    item.reference_boite_physique || "Non classé"
                  }}</code></span
                >
                <span>•</span>
                <span>{{ item.documents_count }} document(s) scellé(s)</span>
              </p>

              <p class="text-[11px] text-muted-foreground mt-1">
                Archivé le
                {{ new Date(item.date_archivage).toLocaleDateString("fr-FR") }}
                par {{ item.archive_par_nom }} ({{ item.archive_par_role }})
                <span v-if="item.observations" class="italic ml-2"
                  >« {{ item.observations }} »</span
                >
              </p>
            </div>
          </div>

          <!-- Actions d'archive -->
          <div class="flex items-center gap-2 shrink-0 self-end md:self-center">
            <button
              type="button"
              @click="emit('view-archive', item)"
              class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-foreground hover:bg-muted transition"
              title="Consulter les pièces du dossier"
            >
              <Eye class="h-3.5 w-3.5 text-brand" />
              Consulter
            </button>

            <button
              type="button"
              @click="emit('desarchiver', item)"
              class="inline-flex items-center gap-1.5 px-2.5 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-amber-700 dark:text-amber-400 hover:bg-amber-500/10 transition"
              title="Restaurer le dossier dans les dossiers actifs"
            >
              <RotateCcw class="h-3.5 w-3.5" />
              Désarchiver
            </button>

            <button
              v-if="isManager"
              type="button"
              @click="emit('delete-archive', item)"
              class="inline-flex items-center gap-1.5 px-2 py-1.5 rounded-[4px] border border-border bg-card text-xs font-medium text-rose-600 hover:bg-rose-500/10 transition"
              title="Suppression définitive certifiée (Directeur)"
            >
              <Trash2 class="h-3.5 w-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
