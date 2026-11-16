<script setup lang="ts">
import type { UserTeamMember } from "../../../types";
import {
  Users,
  Building2,
  Phone,
  Mail,
  Shield,
  Smartphone,
  Monitor,
  Share2,
  Pencil,
  Trash2,
} from "lucide-vue-next";

defineProps<{
  member: UserTeamMember;
  isCurrentUser: boolean;
}>();

const emit = defineEmits<{
  (e: "edit", member: UserTeamMember): void;
  (e: "delete", member: UserTeamMember): void;
}>();

const roleBadges: Record<
  string,
  { label: string; box: string; dot: string; icon: any }
> = {
  DIRECTEUR: {
    label: "Directeur",
    box: "bg-indigo-500/10 text-indigo-700 dark:text-indigo-400 border-indigo-500/20",
    dot: "bg-indigo-500",
    icon: Shield,
  },
  SECRETAIRE: {
    label: "Secrétaire",
    box: "bg-emerald-500/10 text-emerald-700 dark:text-emerald-400 border-emerald-500/20",
    dot: "bg-emerald-500",
    icon: Monitor,
  },
  CAISSIER: {
    label: "Caissière",
    box: "bg-amber-500/10 text-amber-700 dark:text-amber-400 border-amber-500/20",
    dot: "bg-amber-500",
    icon: Monitor,
  },
  MONITEUR: {
    label: "Moniteur",
    box: "bg-sky-500/10 text-sky-700 dark:text-sky-400 border-sky-500/20",
    dot: "bg-sky-500",
    icon: Smartphone,
  },
  ELEVE: {
    label: "Apprenant",
    box: "bg-purple-500/10 text-purple-700 dark:text-purple-400 border-purple-500/20",
    dot: "bg-purple-500",
    icon: Smartphone,
  },
};

const getRoleConfig = (role: string) =>
  roleBadges[role] || {
    label: role,
    box: "bg-muted text-muted-foreground border-border",
    dot: "bg-muted-foreground",
    icon: Users,
  };
</script>

<template>
  <div
    class="flex flex-col justify-between rounded-[5px] border border-border bg-card p-4 transition hover:border-brand-border hover:shadow-xs"
  >
    <div>
      <div class="flex items-start justify-between gap-2">
        <div class="flex items-center gap-3 min-w-0">
          <div
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-brand-light font-bold text-sm text-brand-dark border border-brand-border"
          >
            {{ (member.prenoms?.[0] || "U").toUpperCase()
            }}{{ (member.nom?.[0] || "").toUpperCase() }}
          </div>
          <div class="min-w-0">
            <h3 class="truncate text-sm font-semibold text-foreground">
              {{ member.prenoms }} {{ member.nom }}
            </h3>
            <div class="flex items-center gap-1.5 text-xs text-muted-foreground mt-0.5">
              <Building2 class="h-3 w-3 shrink-0 text-brand" />
              <span class="truncate">{{
                member.agence_nom || "Agence Principale"
              }}</span>
            </div>
          </div>
        </div>

        <span
          class="inline-flex shrink-0 items-center gap-1 rounded-full border px-2 py-0.5 text-[10px] font-medium"
          :class="getRoleConfig(member.role).box"
        >
          <span
            class="h-1.5 w-1.5 rounded-full"
            :class="getRoleConfig(member.role).dot"
          ></span>
          {{ getRoleConfig(member.role).label }}
        </span>
      </div>

      <div class="mt-4 space-y-1.5 border-t border-border/50 pt-3 text-xs">
        <div class="flex items-center justify-between text-muted-foreground">
          <span class="flex items-center gap-1.5">
            <Phone class="h-3 w-3" />
            <span class="font-mono text-foreground font-medium">{{
              member.telephone
            }}</span>
          </span>
          <span
            class="inline-flex items-center gap-1 rounded px-1.5 py-0.5 text-[10px] font-medium"
            :class="
              member.statut === 'ACTIF'
                ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
                : 'bg-rose-500/10 text-rose-600 dark:text-rose-400'
            "
          >
            {{ member.statut }}
          </span>
        </div>
        <div
          v-if="member.email"
          class="flex items-center gap-1.5 text-muted-foreground truncate"
        >
          <Mail class="h-3 w-3 shrink-0" />
          <span class="truncate">{{ member.email }}</span>
        </div>
      </div>
    </div>

    <div
      class="mt-4 flex items-center justify-between gap-1.5 border-t border-border/50 pt-3"
    >
      <a
        :href="'tel:' + member.telephone"
        class="inline-flex flex-1 items-center justify-center gap-1 rounded-lg border border-border bg-background py-1.5 text-xs font-medium text-foreground transition hover:bg-muted"
        title="Appeler"
      >
        <Phone class="h-3 w-3 text-brand" />
        <span>Appeler</span>
      </a>
      <a
        :href="'https://wa.me/' + member.telephone.replace(/[^0-9]/g, '')"
        target="_blank"
        class="inline-flex flex-1 items-center justify-center gap-1 rounded-lg border border-emerald-600/30 bg-emerald-500/10 py-1.5 text-xs font-medium text-emerald-700 transition hover:bg-emerald-500/20 dark:text-emerald-400"
        title="WhatsApp"
      >
        <Share2 class="h-3 w-3" />
        <span>WhatsApp</span>
      </a>
      <button
        type="button"
        @click="emit('edit', member)"
        class="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-border bg-background text-foreground transition hover:bg-muted hover:border-brand/50"
        title="Modifier ce collaborateur"
      >
        <Pencil
          class="h-3.5 w-3.5 text-muted-foreground hover:text-foreground"
        />
      </button>
      <button
        v-if="!isCurrentUser"
        type="button"
        @click="emit('delete', member)"
        class="inline-flex h-8 w-8 items-center justify-center rounded-lg border border-rose-500/30 bg-rose-500/10 text-rose-600 transition hover:bg-rose-500/20 dark:text-rose-400"
        title="Supprimer ce collaborateur"
      >
        <Trash2 class="h-3.5 w-3.5" />
      </button>
    </div>
  </div>
</template>

