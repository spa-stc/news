<script lang="ts">
  import Router from "svelte-spa-router";
  import { routes } from "./routes";
  import Navbar from "$lib/components/nav/Navbar.svelte";
  import { ModeWatcher } from "mode-watcher";
  import Providers from "$lib/providers/Providers.svelte";
</script>

<Providers>
  <div class="relative flex min-h-screen flex-col">
    <ModeWatcher />
    <Navbar />
    <Router {routes} />
  </div>
</Providers>
