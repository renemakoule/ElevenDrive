<script setup lang="ts">
import type { SeanceConduite } from "../../../types";
import {
  extraireCategorie,
  getStatutBadge,
  getProgressPercent,
} from "./planningConduiteUtils";
import { Car, Clock, User, MessageSquare, Trash2 } from "lucide-vue-next";

defineProps<{
  seances: SeanceConduite[];
  isCancellingId: string | null;
}>();

const emit = defineEmits<{
  (e: "whatsapp", seance: SeanceConduite): void;
  (e: "cancel", seance: SeanceConduite): void;
}>();
</script>

<template>
  <div class="rounded-[4px] bg-card border border-border overflow-hidden">
    <div
      v-if="seances.length === 0"
      class="p-12 text-center text-xs text-muted-foreground space-y-2"
    >
      <Car class="w-8 h-8 text-muted-foreground/40 mx-auto" />
      <p>Aucune séance de conduite correspondant aux critères de recherche.</p>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="w-full text-left text-xs">
        <thead
          class="bg-muted text-[11px] text-muted-foreground uppercase tracking-wider border-b border-border"
        >
          <tr>
            <th class="py-2.5 px-4 font-bold">Horaire</th>
            <th class="py-2.5 px-4 font-bold">Apprenant / Candidat</th>
            <th class="py-2.5 px-4 font-bold">Moniteur Référent</th>
            <th class="py-2.5 px-4 font-bold">Véhicule École</th>
            <th class="py-2.5 px-4 font-bold">Avancement</th>
            <th class="py-2.5 px-4 font-bold">Statut</th>
            <th class="py-2.5 px-4 font-bold text-right">Actions</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-border">
          <tr
            v-for="s in seances"
            :key="s.id"
            class="hover:bg-muted/20 transition"
          >
            <td
              class="py-3 px-4 font-mono font-bold text-brand whitespace-nowrap"
            >
              <div class="flex items-center gap-1.5">
                <Clock class="w-3.5 h-3.5 text-brand" />
                <span>{{ s.heure_debut }} – {{ s.heure_fin }}</span>
              </div>
            </td>

            <td class="py-3 px-4 whitespace-nowrap">
              <div class="flex items-center space-x-2.5">
                <div
                  class="w-7 h-7 rounded-full bg-brand text-white font-black flex items-center justify-center text-xs"
                >
                  {{ (s.candidat_nom || "C").charAt(0) }}
                </div>
                <div>
                  <div class="flex items-center gap-1.5">
                    <span class="font-bold text-foreground">{{
                      s.candidat_nom
                    }}</span>
                    <span
                      class="text-[9px] font-mono font-bold px-1.5 py-0.2 rounded bg-sky-500/10 text-sky-600 border border-sky-500/20"
                    >
                      Permis {{ extraireCategorie(s.formule_nom) }}
                    </span>
                  </div>
                  <span class="text-[10px] text-muted-foreground block"
                    >📱 {{ s.telephone_whatsapp || "Non renseigné" }}</span
                  >
                </div>
              </div>
            </td>

            <td class="py-3 px-4 whitespace-nowrap">
              <div
                class="flex items-center gap-1.5 text-foreground font-medium"
              >
                <User class="w-3.5 h-3.5 text-muted-foreground" />
                <span>{{ s.moniteur_nom || "Moniteur Assigné" }}</span>
              </div>
            </td>

            <td class="py-3 px-4 whitespace-nowrap">
              <div>
                <span class="font-bold font-mono text-foreground">{{
                  s.vehicule_immatriculation || "En rotation"
                }}</span>
                <span
                  v-if="s.vehicule_modele"
                  class="text-[10px] text-muted-foreground block"
                >
                  {{ s.vehicule_modele }}
                </span>
              </div>
            </td>

            <td class="py-3 px-4 whitespace-nowrap">
              <div class="w-28 space-y-1">
                <div
                  class="flex items-center justify-between text-[10px] font-mono text-muted-foreground"
                >
                  <span>{{ s.heures_deja_effectuees || 0 }}h</span>
                  <span>{{ s.heures_forfait_total || 20 }}h</span>
                </div>
                <div class="w-full h-1.5 bg-muted rounded-full overflow-hidden">
                  <div
                    class="h-full bg-brand rounded-full"
                    :style="{
                      width: `${getProgressPercent(s.heures_deja_effectuees, s.heures_forfait_total)}%`,
                    }"
                  />
                </div>
              </div>
            </td>

            <td class="py-3 px-4 whitespace-nowrap">
              <span
                class="text-[10px] font-bold px-2 py-0.5 rounded border uppercase tracking-wider"
                :class="getStatutBadge(s.statut).class"
              >
                {{ getStatutBadge(s.statut).label }}
              </span>
            </td>

            <td class="py-3 px-4 whitespace-nowrap text-right space-x-1.5">
              <button
                type="button"
                @click="emit('whatsapp', s)"
                class="p-1.5 rounded-[4px] bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-600 transition cursor-pointer"
                title="WhatsApp"
              >
                <MessageSquare class="w-3.5 h-3.5" />
              </button>

              <button
                type="button"
                @click="emit('cancel', s)"
                :disabled="isCancellingId === s.id"
                class="p-1.5 rounded-[4px] text-muted-foreground hover:text-red-500 hover:bg-red-500/10 transition cursor-pointer"
                title="Supprimer la réservation"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
