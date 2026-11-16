<script setup lang="ts">
import type { SeanceConduite } from "../../../types";
import {
  extraireCategorie,
  getCreneauIcon,
  getDureeCreneauLabel,
  getStatutBadge,
  getProgressPercent,
} from "./planningConduiteUtils";
import { Plus, MessageSquare, Trash2, Clock } from "lucide-vue-next";

export interface CreneauGroupe {
  id: string;
  label: string;
  debut: string;
  fin: string;
  icon?: string;
  seances: SeanceConduite[];
}

defineProps<{
  creneauxGroupes: CreneauGroupe[];
  seancesHorsCreneaux: SeanceConduite[];
  isCancellingId: string | null;
}>();

const emit = defineEmits<{
  (e: "assign", creneau: { debut: string; fin: string }): void;
  (e: "whatsapp", seance: SeanceConduite): void;
  (e: "cancel", seance: SeanceConduite): void;
}>();
</script>

<template>
  <div class="space-y-4">
    <div
      v-for="creneau in creneauxGroupes"
      :key="creneau.id"
      class="rounded-[4px] bg-card border border-border overflow-hidden"
    >
      <!-- En-tête du créneau -->
      <div
        class="px-4 py-2.5 bg-muted/30 border-b border-border flex items-center justify-between"
      >
        <div class="flex items-center space-x-2.5">
          <span class="text-sm">{{
            creneau.icon || getCreneauIcon(creneau.debut)
          }}</span>
          <div>
            <span class="text-xs font-bold text-foreground">{{
              creneau.label
            }}</span>
            <span class="text-xs font-mono font-bold text-primary ml-2">
              {{ creneau.debut }} – {{ creneau.fin }}
            </span>
            <span class="text-[10px] text-muted-foreground ml-2">{{
              getDureeCreneauLabel(creneau.debut, creneau.fin)
            }}</span>
          </div>
        </div>

        <div class="flex items-center space-x-2">
          <span
            class="text-[10px] font-semibold px-2 py-0.5 rounded-full"
            :class="
              creneau.seances.length > 0
                ? 'bg-brand text-white font-bold'
                : 'bg-muted text-muted-foreground'
            "
          >
            {{ creneau.seances.length }} leçon(s)
          </span>

          <button
            type="button"
            @click="emit('assign', { debut: creneau.debut, fin: creneau.fin })"
            class="px-2 py-1 rounded-[4px] bg-brand-light hover:bg-brand-light/80 text-primary text-[11px] font-bold border border-brand-border flex items-center gap-1 transition cursor-pointer"
          >
            <Plus class="w-3 h-3" />
            <span>Affecter</span>
          </button>
        </div>
      </div>

      <!-- Corps : Liste des séances sur ce créneau -->
      <div class="p-3">
        <div
          v-if="creneau.seances.length === 0"
          class="py-4 text-center text-xs text-muted-foreground/60 italic border border-dashed border-border rounded-[4px]"
        >
          Aucune leçon programmée sur cette plage horaire.
        </div>

        <div
          v-else
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3"
        >
          <div
            v-for="s in creneau.seances"
            :key="s.id"
            class="p-3 rounded-[4px] bg-background border border-border hover:border-brand-border transition flex flex-col justify-between space-y-2.5"
          >
            <div class="space-y-1.5">
              <div class="flex items-start justify-between gap-2">
                <div class="flex items-center space-x-2 min-w-0">
                  <div
                    class="w-7 h-7 rounded-full bg-brand text-white font-black flex items-center justify-center text-xs shrink-0"
                  >
                    {{ (s.candidat_nom || "E").charAt(0) }}
                  </div>
                  <div class="min-w-0">
                    <p
                      class="text-xs font-bold text-foreground truncate"
                      :title="s.candidat_nom"
                    >
                      {{ s.candidat_nom }}
                    </p>
                    <div
                      class="text-[10px] text-muted-foreground truncate flex items-center gap-1.5 mt-0.5"
                    >
                      <span
                        class="text-[9px] font-mono font-bold px-1.5 py-0.2 rounded bg-sky-500/10 text-sky-600 border border-sky-500/20"
                      >
                        Permis {{ extraireCategorie(s.formule_nom) }}
                      </span>
                      <span class="truncate">{{
                        s.formule_nom || "Formation pratique"
                      }}</span>
                    </div>
                  </div>
                </div>

                <span
                  class="text-[9px] font-bold px-1.5 py-0.5 rounded border shrink-0 uppercase tracking-wider"
                  :class="getStatutBadge(s.statut).class"
                >
                  {{ getStatutBadge(s.statut).label }}
                </span>
              </div>

              <!-- Progression heures -->
              <div class="space-y-1 pt-1">
                <div
                  class="flex items-center justify-between text-[10px] text-muted-foreground"
                >
                  <span>Heures effectuées :</span>
                  <span class="font-bold font-mono text-foreground">
                    {{ s.heures_deja_effectuees || 0 }}h /
                    {{ s.heures_forfait_total || 20 }}h
                  </span>
                </div>
                <div class="w-full h-1.5 bg-muted rounded-full overflow-hidden">
                  <div
                    class="h-full bg-brand rounded-full transition-all duration-300"
                    :style="{
                      width: `${getProgressPercent(s.heures_deja_effectuees, s.heures_forfait_total)}%`,
                    }"
                  />
                </div>
              </div>

              <!-- Moniteur & Véhicule -->
              <div
                class="grid grid-cols-2 gap-1.5 pt-1 text-[10px] text-muted-foreground border-t border-border"
              >
                <div class="truncate">
                  <span class="font-medium text-foreground block"
                    >👤 Moniteur :</span
                  >
                  <span class="truncate">{{
                    s.moniteur_nom || "Non assigné"
                  }}</span>
                </div>
                <div class="truncate">
                  <span class="font-medium text-foreground block"
                    >🚗 Véhicule :</span
                  >
                  <span
                    class="font-mono truncate font-semibold text-foreground"
                  >
                    {{
                      s.vehicule_immatriculation ||
                      s.vehicule_modele ||
                      "En rotation"
                    }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Actions directes -->
            <div
              class="pt-2 border-t border-border flex items-center justify-between gap-2"
            >
              <button
                type="button"
                @click="emit('whatsapp', s)"
                class="flex-1 flex items-center justify-center space-x-1 px-2 py-1 rounded-[4px] bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-600 font-semibold text-[11px] transition cursor-pointer"
                title="Confirmer la séance par WhatsApp"
              >
                <MessageSquare class="w-3 h-3" />
                <span>WhatsApp</span>
              </button>

              <button
                type="button"
                @click="emit('cancel', s)"
                :disabled="isCancellingId === s.id"
                class="p-1 rounded-[4px] text-muted-foreground hover:text-red-500 hover:bg-red-500/10 transition cursor-pointer"
                title="Annuler / Libérer ce créneau"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bloc séances hors créneaux si existantes -->
    <div
      v-if="seancesHorsCreneaux.length > 0"
      class="rounded-[4px] bg-card border border-border overflow-hidden"
    >
      <div
        class="px-4 py-2.5 bg-amber-500/10 border-b border-border flex items-center justify-between"
      >
        <div class="flex items-center space-x-2">
          <Clock class="w-4 h-4 text-amber-600" />
          <span class="text-xs font-bold text-foreground"
            >Autres Créneaux Horaires Spécifiques</span
          >
        </div>
        <span
          class="text-[10px] font-semibold px-2 py-0.5 rounded-full bg-amber-500/20 text-amber-700"
        >
          {{ seancesHorsCreneaux.length }} leçon(s)
        </span>
      </div>

      <div class="p-3 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
        <div
          v-for="s in seancesHorsCreneaux"
          :key="s.id"
          class="p-3 rounded-[4px] bg-background border border-border space-y-2"
        >
          <div class="flex items-center justify-between text-xs">
            <span class="font-bold text-foreground">{{ s.candidat_nom }}</span>
            <span class="font-mono text-primary font-bold"
              >{{ s.heure_debut }} - {{ s.heure_fin }}</span
            >
          </div>
          <p class="text-[10px] text-muted-foreground">
            Moniteur : {{ s.moniteur_nom }} · Véhicule :
            {{ s.vehicule_immatriculation }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

