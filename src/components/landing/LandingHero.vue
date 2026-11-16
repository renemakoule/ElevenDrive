<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, computed } from "vue";
import {
  ArrowRight,
  Wallet,
  FileText,
  BookOpen,
  Car,
  Archive,
  CalendarCheck,
  Users,
  ShieldCheck,
  Receipt,
  WifiOff,
  BarChart3,
  Bell,
} from "lucide-vue-next";
import HeroDashboardPreview from "./HeroDashboardPreview.vue";
import { useLanguageStore } from "../../stores/languageStore";

const emit = defineEmits<{
  (e: "navigate-auth"): void;
  (e: "scroll-to", id: string): void;
}>();

const languageStore = useLanguageStore();
const t = (text: string) => languageStore.t(text);

const marqueeItems = computed(() => [
  { icon: Wallet, label: t("Encaissements") },
  { icon: FileText, label: t("Dossiers candidats") },
  { icon: BookOpen, label: t("Code") },
  { icon: Car, label: t("Conduite") },
  { icon: Archive, label: t("Archivage") },
  { icon: CalendarCheck, label: t("Planning") },
  { icon: Users, label: t("Moniteurs") },
  { icon: Receipt, label: t("Reçus") },
  { icon: BarChart3, label: t("Statistiques") },
  { icon: Bell, label: t("Rappels") },
  { icon: ShieldCheck, label: t("Sécurité") },
  { icon: WifiOff, label: t("Mode hors ligne") },
]);

/* ───────── Vague de points (canvas) ───────── */
const GRID = 14; // espacement des points (px)
const WAVELENGTH = 260; // distance entre deux crêtes (px)
const PERIOD = 9; // secondes pour qu'une crête avance d'une longueur d'onde
const SHARPNESS = 8; // plus grand = anneau plus fin
const BASE_ALPHA = 0.07; // opacité des points au repos
const PEAK_ALPHA = 0.28; // opacité maximale au passage de l'onde
const BASE_RADIUS = 0.6; // rayon des points au repos (px)
const PEAK_GROW = 0.5; // grossissement max au passage de l'onde (px)
const WAVE_COLOR = "rgb(37, 99, 235)"; // à remplacer par votre brand-500
const FRAME_MS = 33; // ~30 images/seconde

const hostRef = ref<HTMLElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);
let raf = 0;
let running = false;
let visible = true;
let last = 0;
let ro: ResizeObserver | null = null;
let io: IntersectionObserver | null = null;
let reduced = false;
let w = 0;
let h = 0;
let dpr = 1;
let baseColor = "rgb(0, 0, 0)";
let lastColorRead = -Infinity;

function resize() {
  const c = canvasRef.value;
  if (!c) return;
  const rect = c.getBoundingClientRect();
  dpr = Math.min(window.devicePixelRatio || 1, 2);
  w = rect.width;
  h = rect.height;
  c.width = Math.round(w * dpr);
  c.height = Math.round(h * dpr);
}

function draw(now: number) {
  const c = canvasRef.value;
  const ctx = c?.getContext("2d");
  if (!c || !ctx || w === 0) return;

  // Relit la couleur du thème une fois par seconde (mode clair/sombre)
  if (now - lastColorRead > 1000) {
    baseColor = getComputedStyle(c).color;
    lastColorRead = now;
  }

  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);
  ctx.clearRect(0, 0, w, h);

  const t = now / 1000;
  const cx = w / 2;
  const cy = h / 2;
  const rx = 0.6 * w; // même forme que l'ancien masque CSS
  const ry = 0.7 * h;
  const cols = Math.ceil(w / GRID) + 1;
  const rows = Math.ceil(h / GRID) + 1;

  let currentFill = "";

  for (let i = 0; i < cols; i++) {
    for (let j = 0; j < rows; j++) {
      const x = i * GRID + GRID / 2;
      const y = j * GRID + GRID / 2;

      // Masque intégré : transparent au centre (zone du texte), visible sur les côtés
      const r = Math.hypot((x - cx) / rx, (y - cy) / ry);
      const mask = Math.min(Math.max((r - 0.25) / 0.65, 0), 1);
      if (mask < 0.03) continue;

      const d = Math.hypot(x - cx, y - cy);
      const wave =
        0.5 + 0.5 * Math.cos(2 * Math.PI * (d / WAVELENGTH - t / PERIOD));
      const k = Math.pow(wave, SHARPNESS);

      const fill = k < 0.05 ? baseColor : WAVE_COLOR;
      if (fill !== currentFill) {
        ctx.fillStyle = fill;
        currentFill = fill;
      }
      ctx.globalAlpha = (BASE_ALPHA + (PEAK_ALPHA - BASE_ALPHA) * k) * mask;
      const s = (BASE_RADIUS + PEAK_GROW * k) * 2;
      ctx.fillRect(x - s / 2, y - s / 2, s, s);
    }
  }
  ctx.globalAlpha = 1;
}

function loop(now: number) {
  if (!running) return;
  raf = requestAnimationFrame(loop);
  if (now - last < FRAME_MS) return;
  last = now;
  draw(now);
}

function start() {
  if (running || reduced || !visible) return;
  running = true;
  raf = requestAnimationFrame(loop);
}

function stop() {
  running = false;
  cancelAnimationFrame(raf);
}

onMounted(() => {
  reduced = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  resize();

  ro = new ResizeObserver(() => {
    resize();
    if (!running) draw(performance.now());
  });
  if (canvasRef.value) ro.observe(canvasRef.value);

  // Met la trame en pause quand le hero n'est plus visible
  io = new IntersectionObserver(([entry]) => {
    visible = entry.isIntersecting;
    visible ? start() : stop();
  });
  if (hostRef.value) io.observe(hostRef.value);

  if (reduced)
    draw(0); // image fixe
  else start();
});

onBeforeUnmount(() => {
  stop();
  ro?.disconnect();
  io?.disconnect();
});
</script>

<template>
  <section class="relative bg-background">
    <!-- Conteneur : ses bordures forment les deux traits verticaux -->
    <div class="relative mx-auto max-w-6xl border-x border-border">
      <!-- Bloc texte : la trame est limitée à cette zone -->
      <div
        ref="hostRef"
        class="relative overflow-hidden px-5 pb-14 pt-16 text-center sm:px-6 sm:pb-16 sm:pt-24"
      >
        <!-- Trame de fond (le masque est calculé dans le canvas) -->
        <canvas
          ref="canvasRef"
          class="pointer-events-none absolute inset-0 h-full w-full text-foreground"
          aria-hidden="true"
        />

        <!-- Halos de marque statiques (un flou fixe n'est calculé qu'une fois) -->
        <div
          class="pointer-events-none absolute -left-16 -top-16 h-80 w-80 rounded-full bg-brand/15 blur-3xl"
          aria-hidden="true"
        />
        <div
          class="pointer-events-none absolute -bottom-20 -right-16 h-80 w-80 rounded-full bg-brand/10 blur-3xl"
          aria-hidden="true"
        />

        <!-- Contenu (au-dessus de la trame) -->
        <div class="relative">
          <span
            class="inline-flex items-center gap-2 rounded-[4px] border border-border bg-card px-2.5 py-1 text-[11px] font-medium text-muted-foreground"
          >
            <span class="h-1.5 w-1.5 rounded-full bg-brand-500" />
            {{ t("Plateforme de gestion pour auto-écoles") }}
          </span>

          <h1
            class="mx-auto mt-6 max-w-3xl text-[32px] font-semibold leading-[1.1] tracking-[-0.02em] text-foreground sm:text-[48px]"
          >
            {{ t("Pilotez votre auto-école,") }}
            <span class="text-brand-600">{{ t("du guichet à la route.") }}</span>
          </h1>

          <p
            class="mx-auto mt-5 max-w-xl text-[15px] leading-relaxed text-muted-foreground sm:text-base"
          >
            {{
              t(
                "Encaissements, dossiers candidats, code, conduite et archivage réunis dans une seule plateforme, en temps réel et même sans connexion."
              )
            }}
          </p>

          <div
            class="mt-8 flex flex-col items-center justify-center gap-2.5 sm:flex-row"
          >
            <button
              type="button"
              class="inline-flex w-full items-center justify-center gap-1.5 rounded-[4px] bg-brand-600 px-4 py-2 text-[13px] font-medium text-white transition hover:bg-brand-700 active:scale-[0.98] sm:w-auto"
              @click="emit('navigate-auth')"
            >
              {{ t("Accéder à mon espace") }}
              <ArrowRight class="h-3.5 w-3.5" />
            </button>
            <button
              type="button"
              class="inline-flex w-full items-center justify-center gap-1.5 rounded-[4px] border border-border bg-card px-4 py-2 text-[13px] font-medium text-foreground transition hover:bg-muted active:scale-[0.98] sm:w-auto"
              @click="emit('scroll-to', 'solutions')"
            >
              {{ t("Découvrir les solutions") }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Trait horizontal pleine largeur (après les CTA) -->
    <div class="relative border-t border-border">
      <div class="mx-auto max-w-6xl border-x border-border">
        <!-- Marquee -->
        <div
          class="marquee relative overflow-hidden py-5"
          aria-label="Fonctionnalités de la plateforme"
        >
          <div class="marquee-track flex w-max">
            <!-- Deux groupes identiques : le second est masqué aux lecteurs d'écran -->
            <div
              v-for="n in 2"
              :key="n"
              class="flex shrink-0 items-center gap-10 pr-10"
              :aria-hidden="n === 2"
            >
              <div
                v-for="item in marqueeItems"
                :key="`${n}-${item.label}`"
                class="flex shrink-0 items-center gap-2 text-muted-foreground"
              >
                <component
                  :is="item.icon"
                  class="h-4 w-4"
                  :stroke-width="1.75"
                />
                <span class="whitespace-nowrap text-[13px] font-medium">
                  {{ item.label }}
                </span>
              </div>
            </div>
          </div>

          <div
            class="pointer-events-none absolute inset-y-0 left-0 w-24 bg-gradient-to-r from-background to-transparent"
          />
          <div
            class="pointer-events-none absolute inset-y-0 right-0 w-24 bg-gradient-to-l from-background to-transparent"
          />
        </div>
      </div>
    </div>

    <!-- Second trait horizontal pleine largeur -->
    <div class="relative border-t border-border">
      <div
        class="mx-auto max-w-6xl border-x border-border px-5 pb-20 pt-12 sm:px-6 sm:pt-16"
      >
        <div class="mx-auto max-w-5xl">
          <HeroDashboardPreview />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
/* Couche isolée : le défilement ne dépend plus du reste de la page */
.marquee {
  contain: layout paint;
  transform: translateZ(0);
}
.marquee-track {
  animation: marquee 60s linear infinite;
  will-change: transform;
  backface-visibility: hidden;
}
.marquee:hover .marquee-track {
  animation-play-state: paused;
}
/* Chaque groupe = moitié exacte de la piste : boucle sans saut */
@keyframes marquee {
  from {
    transform: translate3d(0, 0, 0);
  }
  to {
    transform: translate3d(-50%, 0, 0);
  }
}
/* Animations réduites : on ralentit au lieu d'arrêter */
@media (prefers-reduced-motion: reduce) {
  .marquee-track {
    animation-duration: 120s;
  }
}
</style>
