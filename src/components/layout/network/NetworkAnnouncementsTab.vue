<script setup lang="ts">
import type { NetworkAnnouncement } from "../../../types";
import { Megaphone, RotateCcw } from "lucide-vue-next";
import NetworkAnnouncementCard from "./NetworkAnnouncementCard.vue";

defineProps<{
  announcements: NetworkAnnouncement[];
  dismissedCount: number;
}>();

const emit = defineEmits<{
  (e: "preview", ann: NetworkAnnouncement): void;
  (e: "dismiss", id: string): void;
  (e: "restore"): void;
}>();
</script>

<template>
  <div>
    <!-- État vide si tous les panneaux sont fermés/masqués -->
    <div
      v-if="announcements.length === 0"
      class="py-14 text-center text-muted-foreground space-y-2"
    >
      <Megaphone class="w-10 h-10 text-muted-foreground mx-auto opacity-40" />
      <div class="font-bold text-foreground text-sm">
        Aucun panneau d'annonce actif
      </div>
      <p class="text-[11px] max-w-sm mx-auto">
        Toutes les annonces ont été consultées ou fermées.
      </p>
      <button
        v-if="dismissedCount > 0"
        type="button"
        @click="emit('restore')"
        class="mt-2 inline-flex items-center gap-1.5 text-[11px] text-brand hover:underline font-semibold cursor-pointer"
      >
        <RotateCcw class="w-3.5 h-3.5" />
        <span>Réafficher les {{ dismissedCount }} annonces masquées</span>
      </button>
    </div>

    <!-- Grille responsive des panneaux Ads -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3.5 items-stretch">
      <NetworkAnnouncementCard
        v-for="ann in announcements"
        :key="ann.id"
        :announcement="ann"
        @preview="emit('preview', $event)"
        @dismiss="emit('dismiss', $event)"
      />
    </div>
  </div>
</template>
