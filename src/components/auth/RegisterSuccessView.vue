<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import {
  Sparkles, QrCode, Download, Copy, CheckCircle2, Share2, Maximize2, ArrowRight
} from 'lucide-vue-next';
import type { AuthSessionData } from '../../types';
import { useGuichetStore } from '../../stores/guichetStore';
import { generateQrDataUrl, downloadQrCode } from '../../utils/qrGenerator';

interface AccessLinkItem {
  key: string;
  title: string;
  caption: string;
  shareTitle: string;
  url: string;
  dot: string;
}

interface Props {
  data: AuthSessionData;
  accessLinks: AccessLinkItem[];
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'open-qr', link: AccessLinkItem): void;
  (e: 'finish'): void;
}>();

const guichetStore = useGuichetStore();

const successViewMode = ref<'qrcodes' | 'links'>('qrcodes');
const qrDataUrls = reactive<Record<string, string>>({});
const copiedField = ref<string | null>(null);

const primaryBtn =
  'inline-flex h-10 w-full items-center justify-center gap-2 rounded-lg bg-brand text-sm font-semibold text-white ' +
  'transition hover:bg-brand-dark focus:outline-none focus-visible:ring-2 focus-visible:ring-brand/40 ' +
  'focus-visible:ring-offset-2 focus-visible:ring-offset-background active:scale-[0.99]';

// Génération automatique des QR Codes
watch(
  () => props.accessLinks,
  async (links) => {
    if (links && links.length > 0) {
      for (const link of links) {
        if (link.url && !qrDataUrls[link.key]) {
          qrDataUrls[link.key] = await generateQrDataUrl(link.url, {
            width: 320,
            margin: 2,
            darkColor: '#064e3b',
            lightColor: '#ffffff',
          });
        }
      }
    }
  },
  { immediate: true, deep: true }
);

const handleDownloadQr = (link: AccessLinkItem) => {
  const dataUrl = qrDataUrls[link.key];
  if (dataUrl) {
    const slug = props.data.auto_ecole.code_unique || 'autoecole';
    downloadQrCode(dataUrl, `qrcode-${link.key}-${slug}.png`);
    guichetStore.showToast(`QR Code ${link.title} téléchargé avec succès !`, 'success');
  }
};

const copyToClipboard = async (text: string, fieldName: string) => {
  try {
    await navigator.clipboard.writeText(text);
    copiedField.value = fieldName;
    guichetStore.showToast('Lien copié dans le presse-papiers !', 'success');
    setTimeout(() => {
      if (copiedField.value === fieldName) copiedField.value = null;
    }, 2500);
  } catch {
    const input = document.createElement('input');
    input.value = text;
    document.body.appendChild(input);
    input.select();
    document.execCommand('copy');
    document.body.removeChild(input);
    copiedField.value = fieldName;
    setTimeout(() => {
      if (copiedField.value === fieldName) copiedField.value = null;
    }, 2500);
  }
};

const shareOnWhatsApp = (url: string, title: string) => {
  const text = encodeURIComponent(`Bonjour, voici votre lien d'accès à ${title} : ${url}`);
  window.open(`https://wa.me/?text=${text}`, '_blank');
};
</script>

<template>
  <div class="animate-in fade-in zoom-in-95 duration-200">
    <div class="text-center">
      <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full border border-brand-border bg-brand-light text-brand">
        <Sparkles class="h-5 w-5" />
      </div>
      <h1 class="text-2xl font-semibold tracking-tight text-foreground">
        Votre auto-école est prête
      </h1>
      <p class="mx-auto mt-2 max-w-xs text-sm text-muted-foreground">
        Félicitations {{ data.user.prenom }}. Vos accès sont configurés. Transmettez ces liens à votre équipe et apprenants.
      </p>
    </div>

    <!-- Récapitulatif de l'établissement -->
    <div class="mt-8 flex items-center justify-between gap-3 rounded-xl border border-border bg-muted/40 p-3.5">
      <div class="flex min-w-0 items-center gap-3">
        <div class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-brand text-sm font-semibold text-white">
          {{ data.auto_ecole.nom_commercial.charAt(0) }}
        </div>
        <div class="min-w-0">
          <div class="truncate text-sm font-medium text-foreground">
            {{ data.auto_ecole.nom_commercial }}
          </div>
          <div class="truncate text-xs text-muted-foreground">
            Agrément <span class="font-mono text-foreground">{{ data.auto_ecole.numero_agrement }}</span>
            à {{ data.agence?.commune || 'Siège principal' }}
          </div>
        </div>
      </div>
      <div class="shrink-0 text-right">
        <span class="mb-0.5 block text-[11px] text-muted-foreground">Code établissement</span>
        <span class="rounded-md border border-brand-border bg-brand-light px-2 py-0.5 font-mono text-xs font-semibold text-brand-dark">
          {{ data.auto_ecole.code_unique }}
        </span>
      </div>
    </div>

    <!-- Sélecteur d'affichage : QR Codes vs Liens texte -->
    <div class="mt-6 flex items-center justify-between">
      <span class="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
        Accès mobiles configurés
      </span>
      <div class="inline-flex rounded-lg border border-border bg-muted/60 p-0.5">
        <button
          type="button"
          @click="successViewMode = 'qrcodes'"
          class="inline-flex items-center gap-1.5 rounded-md px-3 py-1 text-xs font-medium transition"
          :class="successViewMode === 'qrcodes' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'"
        >
          <QrCode class="h-3.5 w-3.5 text-brand" />
          <span>QR Codes</span>
        </button>
        <button
          type="button"
          @click="successViewMode = 'links'"
          class="inline-flex items-center gap-1.5 rounded-md px-3 py-1 text-xs font-medium transition"
          :class="successViewMode === 'links' ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'"
        >
          <span>Liens URL</span>
        </button>
      </div>
    </div>

    <!-- VUE 1 : GRILLE DES QR CODES SCANNABLES DIRECTEMENT -->
    <div v-if="successViewMode === 'qrcodes'" class="mt-4 grid gap-3 sm:grid-cols-2">
      <div
        v-for="link in accessLinks"
        :key="link.key"
        class="flex flex-col justify-between rounded-xl border border-border bg-card p-4 transition-all hover:border-brand-border hover:shadow-sm"
      >
        <div>
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2">
              <span class="h-2 w-2 rounded-full" :class="link.dot"></span>
              <span class="text-sm font-semibold text-foreground">{{ link.title }}</span>
            </div>
            <button
              type="button"
              @click="emit('open-qr', link)"
              title="Agrandir et imprimer"
              class="rounded p-1 text-muted-foreground hover:bg-muted hover:text-foreground transition"
            >
              <Maximize2 class="h-3.5 w-3.5" />
            </button>
          </div>
          <p class="mt-0.5 text-xs text-muted-foreground leading-tight">{{ link.caption }}</p>
        </div>

        <!-- Visuel du QR Code -->
        <div class="my-3 flex flex-col items-center">
          <div class="relative group rounded-xl border border-brand-border/60 bg-white p-2.5 shadow-sm">
            <img
              v-if="qrDataUrls[link.key]"
              :src="qrDataUrls[link.key]"
              :alt="'QR Code ' + link.title"
              class="h-36 w-36 object-contain rounded select-none cursor-pointer"
              @click="emit('open-qr', link)"
            />
            <div v-else class="flex h-36 w-36 items-center justify-center text-xs text-muted-foreground">
              <Sparkles class="h-4 w-4 animate-spin text-brand" />
            </div>
          </div>
          <span class="mt-2 text-[10px] font-medium text-muted-foreground">
            Visez avec l'appareil photo 📸
          </span>
        </div>

        <!-- Barre d'actions -->
        <div class="grid grid-cols-3 gap-1.5 pt-1 border-t border-border/60">
          <button
            type="button"
            @click="handleDownloadQr(link)"
            title="Télécharger l'image PNG"
            class="inline-flex h-8 items-center justify-center gap-1 rounded-md border border-border bg-background px-1.5 text-[11px] font-medium text-foreground transition hover:bg-muted"
          >
            <Download class="h-3 w-3 text-brand" />
            <span>Image</span>
          </button>
          <button
            type="button"
            @click="copyToClipboard(link.url, link.key)"
            class="inline-flex h-8 items-center justify-center gap-1 rounded-md border border-border bg-background px-1.5 text-[11px] font-medium text-foreground transition hover:bg-muted"
          >
            <CheckCircle2 v-if="copiedField === link.key" class="h-3 w-3 text-brand" />
            <Copy v-else class="h-3 w-3 text-muted-foreground" />
            <span>{{ copiedField === link.key ? 'Copié' : 'Lien' }}</span>
          </button>
          <button
            type="button"
            @click="shareOnWhatsApp(link.url, link.shareTitle)"
            title="Partager sur WhatsApp"
            class="inline-flex h-8 items-center justify-center gap-1 rounded-md border border-emerald-600/30 bg-emerald-500/10 px-1.5 text-[11px] font-medium text-emerald-700 transition hover:bg-emerald-500/20 dark:text-emerald-400"
          >
            <Share2 class="h-3 w-3" />
            <span>WhatsApp</span>
          </button>
        </div>
      </div>
    </div>

    <!-- VUE 2 : LIENS D'ACCÈS DÉDIÉS COMPACTS -->
    <div v-else class="mt-4 space-y-3">
      <div
        v-for="link in accessLinks"
        :key="link.key"
        class="rounded-xl border border-border bg-card p-3.5 transition-colors hover:border-brand-border"
      >
        <div class="mb-2.5 flex items-center justify-between">
          <div class="flex items-center gap-2">
            <span class="h-2 w-2 rounded-full" :class="link.dot"></span>
            <span class="text-sm font-medium text-foreground">{{ link.title }}</span>
          </div>
          <span class="text-xs text-muted-foreground">{{ link.caption }}</span>
        </div>
        <div class="flex items-center gap-2">
          <input
            type="text"
            readonly
            :value="link.url"
            class="h-9 min-w-0 flex-1 select-all rounded-lg border border-border bg-muted/40 px-2.5 font-mono text-xs text-foreground focus:outline-none focus:border-brand"
          />
          <button
            type="button"
            @click="emit('open-qr', link)"
            title="Afficher le QR Code"
            class="inline-flex h-9 items-center gap-1.5 rounded-lg border border-brand-border bg-brand-light/60 px-2.5 text-xs font-medium text-brand-dark transition hover:bg-brand-light focus:outline-none"
          >
            <QrCode class="h-3.5 w-3.5 text-brand" />
            <span>QR Code</span>
          </button>
          <button
            type="button"
            @click="copyToClipboard(link.url, link.key)"
            class="inline-flex h-9 items-center gap-1.5 rounded-lg border border-border px-2.5 text-xs font-medium text-foreground transition hover:bg-muted focus:outline-none"
          >
            <CheckCircle2 v-if="copiedField === link.key" class="h-3.5 w-3.5 text-brand" />
            <Copy v-else class="h-3.5 w-3.5 text-muted-foreground" />
            <span>{{ copiedField === link.key ? 'Copié' : 'Copier' }}</span>
          </button>
          <button
            type="button"
            @click="shareOnWhatsApp(link.url, link.shareTitle)"
            title="Envoyer par WhatsApp"
            class="inline-flex h-9 items-center rounded-lg border border-emerald-600/30 bg-emerald-500/10 px-2.5 text-xs font-medium text-emerald-700 transition hover:bg-emerald-500/20 dark:text-emerald-400"
          >
            WhatsApp
          </button>
        </div>
      </div>
    </div>

    <button type="button" @click="emit('finish')" :class="[primaryBtn, 'mt-6']">
      <span>Ouvrir mon espace direction et guichet</span>
      <ArrowRight class="h-4 w-4" />
    </button>
  </div>
</template>
