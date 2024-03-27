<script lang="ts">
  import AuthNav from "./AuthNav.svelte";
  import MainNav from "./MainNav.svelte";
  import ModeToggle from "./ModeToggle.svelte";
</script>

<header class="bg-background sticky top-0 z-40 w-full border-b">
  <div
    class="container flex h-16 items-center space-x-4 sm:justify-between sm:space-x-0"
  >
    <MainNav />
    <div class="flex flex-1 items-center justify-end space-x-2">
      <ModeToggle />
      <AuthNav />
    </div>
  </div>
</header>
