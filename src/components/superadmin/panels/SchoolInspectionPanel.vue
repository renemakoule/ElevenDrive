<script setup lang="ts">
import {
  Power,
  ExternalLink,
  X,
  Edit3,
  ShieldCheck,
  Layers,
  Trash2,
  Monitor,
  Calendar,
  MapPin,
  Building2,
  Phone,
} from "lucide-vue-next";
import type { SchoolItemSummary } from "../../../types";

defineProps<{
  school: SchoolItemSummary;
}>();

const emit = defineEmits<{
  (e: "toggle-status", school: SchoolItemSummary): void;
  (e: "edit", school: SchoolItemSummary): void;
  (e: "manage-licence", school: SchoolItemSummary): void;
  (e: "manage-modules", school: SchoolItemSummary): void;
  (e: "delete", school: SchoolItemSummary): void;
  (e: "close"): void;
}>();
</script>

<template>
  <div
    class="w-80 bg-card border-l border-border p-4 flex flex-col justify-between shrink-0 overflow-y-auto"
  >
    <div class="space-y-4">
      <!-- En-tête avec bouton fermer -->
      <div class="border-b border-border pb-3">
        <div class="flex items-center justify-between">
          <span
            class="text-[10px] uppercase font-bold text-muted-foreground tracking-wider"
          >
            Fiche d'Inspection
          </span>
          <button
            type="button"
            @click="emit('close')"
            class="text-muted-foreground hover:text-foreground p-1 rounded-[4px] transition-colors cursor-pointer"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        </div>
        <h3 class="text-base font-bold text-foreground mt-1">
          {{ school.nom_commercial }}
        </h3>
        <p class="text-xs text-muted-foreground font-mono">ID: {{ school.id }}</p>
      </div>

      <!-- Badge Licence & Quotas Machines -->
      <div
        class="p-3 rounded-[4px] bg-muted/50 border border-border space-y-2"
      >
        <div class="flex items-center justify-between">
          <span
            class="text-[10px] uppercase font-bold text-muted-foreground tracking-wider"
            >Statut Licence</span
          >
          <span
            :class="[
              'px-2 py-0.5 rounded-[3px] text-[10px] font-mono font-bold',
              school.statut_licence === 'ACTIVE'
                ? 'bg-brand/10 text-brand border border-brand/30'
                : 'bg-destructive/10 text-destructive border border-destructive/30',
            ]"
          >
            {{
              school.statut_licence ||
              (school.est_actif ? "ACTIVE" : "SUSPENDUE")
            }}
          </span>
        </div>

        <div class="grid grid-cols-2 gap-2 text-xs pt-1">
          <div class="flex items-center gap-1.5 text-foreground">
            <Monitor class="w-3.5 h-3.5 text-brand shrink-0" />
            <span class="font-mono text-[11px]"
              >{{ school.machines_actives || 0 }} /
              {{ school.max_machines || 3 }} postes</span
            >
          </div>

          <div class="flex items-center gap-1.5 text-foreground">
            <Calendar class="w-3.5 h-3.5 text-sky-500 shrink-0" />
            <span class="font-mono text-brand text-[11px]"
              >{{ school.jours_restants ?? 365 }}j restants</span
            >
          </div>
        </div>

        <div class="text-[10px] font-mono text-muted-foreground">
          Échéance :
          {{
            school.date_expiration_licence
              ? new Date(school.date_expiration_licence).toLocaleDateString(
                  "fr-FR",
                )
              : "31/12/2026"
          }}
        </div>
      </div>

      <!-- Cartes de métriques rapides -->
      <div class="grid grid-cols-2 gap-2 text-xs">
        <div class="p-2.5 rounded-[4px] bg-muted/50 border border-border">
          <span class="text-[10px] text-muted-foreground">Agences</span>
          <div class="text-base font-bold font-mono text-foreground mt-0.5">
            {{ school.agences_count }}
          </div>
        </div>
        <div class="p-2.5 rounded-[4px] bg-muted/50 border border-border">
          <span class="text-[10px] text-muted-foreground">Collaborateurs</span>
          <div class="text-base font-bold font-mono text-foreground mt-0.5">
            {{ school.users_count }}
          </div>
        </div>
        <div class="p-2.5 rounded-[4px] bg-muted/50 border border-border">
          <span class="text-[10px] text-muted-foreground">Candidats</span>
          <div class="text-base font-bold font-mono text-brand mt-0.5">
            {{ school.candidats_count }}
          </div>
        </div>
        <div class="p-2.5 rounded-[4px] bg-muted/50 border border-border">
          <span class="text-[10px] text-muted-foreground">Devise</span>
          <div class="text-base font-bold font-mono text-foreground mt-0.5">
            {{ school.devise }}
          </div>
        </div>
      </div>

      <!-- Informations de Contact & Registre -->
      <div class="space-y-2 text-xs">
        <div>
          <span class="text-[10px] text-muted-foreground uppercase font-mono"
            >Agrément officiel</span
          >
          <p class="text-foreground font-mono">{{ school.numero_agrement }}</p>
        </div>
        <div>
          <span class="text-[10px] text-muted-foreground uppercase font-mono"
            >Lien Web Dédié</span
          >
          <p class="text-brand font-mono break-all">
            /{{ school.code_unique }}
          </p>
        </div>
        <div>
          <span class="text-[10px] text-muted-foreground uppercase font-mono"
            >Téléphone & WhatsApp</span
          >
          <p class="text-foreground font-mono">{{ school.telephone }}</p>
        </div>
        <div>
          <span class="text-[10px] text-muted-foreground uppercase font-mono flex items-center gap-1"
            ><MapPin class="w-3 h-3 text-brand" /> Ville / Pays</span
          >
          <p class="text-foreground font-mono">
            {{ school.ville || "Douala" }}, {{ school.pays }}
          </p>
        </div>

        <!-- Liste des Agences avec Localisation Exacte -->
        <div class="pt-2 border-t border-border space-y-2">
          <div class="flex items-center justify-between">
            <span class="text-[10px] text-muted-foreground uppercase font-mono font-bold flex items-center gap-1">
              <Building2 class="w-3 h-3 text-sky-500" />
              Agences & Succursales ({{ school.agences?.length || school.agences_count || 1 }})
            </span>
          </div>

          <div v-if="school.agences && school.agences.length > 0" class="space-y-1.5 max-h-44 overflow-y-auto pr-1">
            <div
              v-for="ag in school.agences"
              :key="ag.id"
              class="p-2 rounded-[4px] bg-muted/40 border border-border text-[11px] space-y-1"
            >
              <div class="flex items-center justify-between font-semibold text-foreground">
                <span class="truncate">{{ ag.nom }}</span>
                <span
                  v-if="ag.est_siege"
                  class="px-1 py-0.2 rounded text-[9px] bg-brand-light text-brand-dark font-bold border border-brand-border"
                >
                  Siège
                </span>
              </div>
              <div class="text-[10px] text-muted-foreground font-mono leading-tight flex items-start gap-1">
                <MapPin class="w-3 h-3 text-muted-foreground shrink-0 mt-0.5" />
                <span>
                  {{ ag.adresse ? `${ag.adresse}, ` : '' }}{{ ag.commune ? `${ag.commune}, ` : '' }}{{ ag.ville || school.ville || 'Douala' }}
                </span>
              </div>
              <div v-if="ag.telephone" class="text-[10px] text-muted-foreground font-mono flex items-center gap-1">
                <Phone class="w-3 h-3 text-muted-foreground shrink-0" />
                <span>{{ ag.telephone }}</span>
              </div>
            </div>
          </div>
          <div v-else class="text-[11px] text-muted-foreground italic font-mono">
            Siège social : {{ school.ville || "Douala" }}
          </div>
        </div>
      </div>
    </div>

    <!-- Actions Super-Admin Complètes -->
    <div class="pt-4 border-t border-border space-y-2">
      <!-- 1. Modifier l'établissement -->
      <button
        type="button"
        @click="emit('edit', school)"
        class="w-full py-2 px-3 rounded-[4px] text-xs font-semibold bg-muted hover:bg-muted/80 text-foreground border border-border flex items-center justify-center gap-2 transition cursor-pointer"
      >
        <Edit3 class="w-3.5 h-3.5 text-sky-500" />
        <span>Modifier l'établissement</span>
      </button>

      <!-- 2. Gérer Licence & Postes Machines -->
      <button
        type="button"
        @click="emit('manage-licence', school)"
        class="w-full py-2 px-3 rounded-[4px] text-xs font-semibold bg-muted hover:bg-muted/80 text-foreground border border-border flex items-center justify-center gap-2 transition cursor-pointer"
      >
        <ShieldCheck class="w-3.5 h-3.5 text-brand" />
        <span>Définir Postes & Durée Licence</span>
      </button>

      <!-- 3. Droits d'Accès & Modules -->
      <button
        type="button"
        @click="emit('manage-modules', school)"
        class="w-full py-2 px-3 rounded-[4px] text-xs font-semibold bg-muted hover:bg-muted/80 text-foreground border border-border flex items-center justify-center gap-2 transition cursor-pointer"
      >
        <Layers class="w-3.5 h-3.5 text-purple-500" />
        <span>Droits d'Accès & Modules</span>
      </button>

      <!-- 4. Suspendre ou Activer -->
      <button
        type="button"
        @click="emit('toggle-status', school)"
        :class="[
          'w-full py-2 px-3 rounded-[4px] text-xs font-semibold flex items-center justify-center gap-2 transition-colors cursor-pointer',
          school.est_actif
            ? 'bg-destructive/10 text-destructive border border-destructive/30 hover:bg-destructive/20'
            : 'bg-brand/10 text-brand border border-brand/30 hover:bg-brand/20',
        ]"
      >
        <Power class="w-3.5 h-3.5" />
        <span>
          {{ school.est_actif ? "Suspendre l'accès" : "Réactiver l'accès" }}
        </span>
      </button>

      <div class="flex items-center gap-2 pt-1">
        <!-- 5. Voir portail public -->
        <a
          :href="'/' + school.code_unique"
          target="_blank"
          class="flex-1 py-1.5 px-2 rounded-[4px] text-xs font-medium text-muted-foreground hover:text-foreground border border-border hover:bg-muted flex items-center justify-center gap-1.5 transition-colors"
        >
          <span>Portail Web</span>
          <ExternalLink class="w-3 h-3" />
        </a>

        <!-- 6. Supprimer définitivement -->
        <button
          type="button"
          @click="emit('delete', school)"
          title="Supprimer définitivement l'auto-école"
          class="p-2 rounded-[4px] text-destructive hover:bg-destructive/10 border border-destructive/30 transition cursor-pointer"
        >
          <Trash2 class="w-3.5 h-3.5" />
        </button>
      </div>
    </div>
  </div>
</template>
