<script lang="ts">
  import Router from "svelte-spa-router";
  import { routes } from "./routes";
  import Navbar from "$lib/components/nav/Navbar.svelte";
  import { ModeWatcher } from "mode-watcher";
  import Providers from "$lib/providers/Providers.svelte";
  import { Toaster } from "svelte-french-toast";
</script>

<Providers>
  <div class="relative flex min-h-screen flex-col">
    <Toaster />
    <ModeWatcher />
    <Navbar />
    <Router {routes} />
  </div>
</Providers>
