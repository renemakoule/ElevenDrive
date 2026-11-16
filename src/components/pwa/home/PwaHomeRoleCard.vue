<script setup lang="ts">
import type { UserRole } from '../../../types';
import { CheckCircle2, ArrowRight } from 'lucide-vue-next';

interface Props {
  role: UserRole;
  title: string;
  tag: string;
  description: string;
  icon: any;
  color: 'purple' | 'sky' | 'indigo';
  features: string[];
  actionLabel: string;
}

const props = defineProps<Props>();

const emit = defineEmits<{
  (e: 'select', role: UserRole): void;
}>();

const colorStyles: Record<string, {
  hoverBorder: string;
  iconBox: string;
  iconText: string;
  tagBadge: string;
  titleHover: string;
  checkIcon: string;
  actionText: string;
}> = {
  purple: {
    hoverBorder: 'hover:border-purple-500/50',
    iconBox: 'bg-purple-500/10 border-purple-500/20',
    iconText: 'text-purple-600 dark:text-purple-400',
    tagBadge: 'bg-purple-500/10 text-purple-700 dark:text-purple-300 border-purple-500/20',
    titleHover: 'group-hover:text-purple-600 dark:group-hover:text-purple-400',
    checkIcon: 'text-purple-500',
    actionText: 'text-purple-600 dark:text-purple-400',
  },
  sky: {
    hoverBorder: 'hover:border-sky-500/50',
    iconBox: 'bg-sky-500/10 border-sky-500/20',
    iconText: 'text-sky-600 dark:text-sky-400',
    tagBadge: 'bg-sky-500/10 text-sky-700 dark:text-sky-300 border-sky-500/20',
    titleHover: 'group-hover:text-sky-600 dark:group-hover:text-sky-400',
    checkIcon: 'text-sky-500',
    actionText: 'text-sky-600 dark:text-sky-400',
  },
  indigo: {
    hoverBorder: 'hover:border-indigo-500/50',
    iconBox: 'bg-indigo-500/10 border-indigo-500/20',
    iconText: 'text-indigo-600 dark:text-indigo-400',
    tagBadge: 'bg-indigo-500/10 text-indigo-700 dark:text-indigo-300 border-indigo-500/20',
    titleHover: 'group-hover:text-indigo-600 dark:group-hover:text-indigo-400',
    checkIcon: 'text-indigo-500',
    actionText: 'text-indigo-600 dark:text-indigo-400',
  },
};

const currentStyle = colorStyles[props.color] || colorStyles.purple;
</script>

<template>
  <div
    @click="emit('select', role)"
    class="group relative flex flex-col justify-between rounded-[4px] border border-border bg-card p-6 hover:shadow-md transition cursor-pointer"
    :class="currentStyle.hoverBorder"
  >
    <div>
      <div class="flex items-center justify-between mb-4">
        <div
          class="w-12 h-12 rounded-[4px] border flex items-center justify-center group-hover:scale-105 transition-transform"
          :class="[currentStyle.iconBox, currentStyle.iconText]"
        >
          <component :is="icon" class="w-6 h-6" />
        </div>
        <span
          class="px-2.5 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider border"
          :class="currentStyle.tagBadge"
        >
          {{ tag }}
        </span>
      </div>

      <h4
        class="text-base font-bold text-foreground transition"
        :class="currentStyle.titleHover"
      >
        {{ title }}
      </h4>
      <p class="mt-2 text-xs text-muted-foreground leading-relaxed">
        {{ description }}
      </p>

      <ul class="mt-4 space-y-1.5 text-xs text-muted-foreground">
        <li v-for="(feat, idx) in features" :key="idx" class="flex items-center gap-1.5">
          <CheckCircle2 class="w-3.5 h-3.5 shrink-0" :class="currentStyle.checkIcon" />
          <span>{{ feat }}</span>
        </li>
      </ul>
    </div>

    <div
      class="mt-6 pt-4 border-t border-border flex items-center justify-between text-xs font-semibold"
      :class="currentStyle.actionText"
    >
      <span>{{ actionLabel }}</span>
      <ArrowRight class="w-4 h-4 transform group-hover:translate-x-1 transition-transform" />
    </div>
  </div>
</template>
