<script lang="ts">
  import { colorCodes } from '$lib/display/colors';
  import ColoredText from '$lib/components/common/ColoredText.svelte';

  type Section = {
    large?: boolean;
    color?: string;
    italic?: boolean;
    text: string[];
  };

  let {
    sections = [],
    color = colorCodes.WHITE
  }: {
    sections?: Section[];
    color?: string;
  } = $props();
</script>

<div class="grid grid-flow-row min-w-[19em] max-w-[25vw] w-fit h-fit border-2" style="border-color: {color}">
  {#each sections as section}
    <div
      class="grid grid-flow-row section p-1.5 text-center"
      class:text-xl={section.large}
      class:italic={section.italic}
      style="--border-color: {color}; color: {section.color || '#ffffff'}">
      {#each section.text as line}
        <div>
          <ColoredText text={line} />
        </div>
      {/each}
    </div>
  {/each}
</div>

<style lang="postcss">
  .section:not(:last-child) {
    border-bottom: 1px solid var(--border-color, #ffffff);
  }
</style>
