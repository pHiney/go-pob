<script lang="ts">
  import '../app.postcss';
  import Header from '../lib/components/Header.svelte';
  import Sidebar from '../lib/components/Sidebar.svelte';

  import { assets } from '$app/paths';
  import { browser } from '$app/environment';
  import { syncWrap } from '$lib/go/worker';
  import { proxy } from 'comlink';
  import type { Outputs } from '$lib/custom_types';
  import { outputs, currentBuild, sampleBuildCode } from '$lib/global';
  import OverlayController from '$lib/components/overlays/OverlayController.svelte';
  import { fontScaling } from '$lib/global.js';
  import { logError } from '$lib/utils';
  import type { Snippet } from 'svelte';

  let {
    children
  }: {
    children?: Snippet;
  } = $props();

  let wasmLoading = $state(true);

  let loadingMessage = $state('Initializing...');
  let loadingStage = $state('');

  if (browser) {
    if (!syncWrap || syncWrap === null) {
      loadingMessage = 'Failed to initialize worker';
    } else {
      syncWrap.booted
        .then((booted) => {
          if (booted) {
            wasmLoading = false;
            return;
          }

          fetch(assets + '/go-pob.wasm')
            .then((data) => data.arrayBuffer())
            .then((data) => {
              syncWrap
                ?.boot(
                  data,
                  proxy((out: Outputs) => {
                    outputs.set(out);
                  }),
                  proxy(currentBuild)
                )
                .then(async () => {
                  console.log('worker booted');

                  loadingMessage = 'Loading data...';
                  await syncWrap?.loadData(
                    // eslint-disable-next-line @typescript-eslint/require-await
                    proxy(async (stage: string) => {
                      loadingMessage = 'Loading data:';
                      loadingStage = stage;
                    })
                  );

                  wasmLoading = false;

                  // TODO Remove from Prod
                  syncWrap
                    ?.ImportCode(sampleBuildCode)
                    .then(() => {
                      syncWrap?.Tick('importBuildFromCode').catch(logError);
                    })
                    .catch(logError);
                })
                .catch(logError);
            })
            .catch(logError);
        })
        .catch(logError);
    }
  }
</script>

<div class="w-screen h-screen max-w-screen max-h-screen overflow-hidden flex flex-col" style="font-size: {$fontScaling}pt">
  {#if wasmLoading}
    <div class="flex flex-row justify-center h-full">
      <div class="flex flex-col justify-center text-5xl text-center">
        {loadingMessage}
        {#if loadingStage !== ''}
          <br />
          {loadingStage}
        {/if}
      </div>
    </div>
  {:else}
    <Header />

    <div class="flex flex-row h-full full-page">
      <Sidebar />

      <div class="h-full w-full overflow-hidden">
        {@render children?.()}
      </div>
    </div>

    <OverlayController />
  {/if}
</div>
